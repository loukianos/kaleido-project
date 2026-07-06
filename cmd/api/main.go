package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"kaleido-project/db/sqlc"
	"kaleido-project/internal/api"
	"kaleido-project/internal/auth"
	"kaleido-project/internal/config"
	"kaleido-project/internal/contracts"
	"kaleido-project/internal/eth"
	"kaleido-project/internal/identity"
	"kaleido-project/internal/keys"
	"kaleido-project/internal/loans"
)

var version = "dev"

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := run(ctx); err != nil {
		slog.Error("startup failed", "error", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	logger := slog.Default()
	logger.Info("configuration loaded", "config", cfg, "version", version)

	startupCtx, cancelStartup := context.WithTimeout(ctx, 10*time.Second)
	defer cancelStartup()

	conn, err := db.Open(startupCtx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer conn.Close()

	ethClient, err := eth.Dial(startupCtx, cfg.EthRPCURL, cfg.ChainID, cfg.DeployerPrivateKey)
	if err != nil {
		return fmt.Errorf("initialize ethereum client: %w", err)
	}
	logger.Info("ethereum client initialized",
		"signer_address", ethClient.SignerAddress().Hex(),
		"chain_id", ethClient.ChainID().String(),
	)

	queries := db.New(conn)
	lockManager := db.NewLockManager(queries)
	lockHolder := fmt.Sprintf("%s:%d", hostname(), os.Getpid())
	encryptor, err := keys.NewAESGCM(cfg.KeyEncryptionMasterKey)
	if err != nil {
		return fmt.Errorf("initialize key encryptor: %w", err)
	}
	identityService := identity.NewService(queries, encryptor)

	poolSigners, err := identityService.EnsureServicerPool(startupCtx, cfg.ServicerKeyPoolSize)
	if err != nil {
		return fmt.Errorf("provision servicer key pool: %w", err)
	}
	poolAddresses := make([]common.Address, 0, len(poolSigners))
	for _, signer := range poolSigners {
		poolAddresses = append(poolAddresses, signer.Address())
	}
	logger.Info("servicer key pool ready", "size", len(poolSigners))

	contractRepo := contracts.NewRepository(queries, conn)
	contractService := contracts.NewService(
		contractRepo,
		ethClient,
		lockManager,
		cfg.LoanBaseURI,
		lockHolder,
		poolAddresses,
	)
	loanRepo := loans.NewRepository(queries, conn)
	loanService := loans.NewService(loanRepo, ethClient, lockManager, lockHolder, identityService, cfg.OIDCIssuerURL, poolSigners)

	// Heal missing pool role grants on contracts deployed before this pool existed (or under a smaller pool size).
	// Best-effort: a down chain shouldn't block startup, and grants retry on the next boot.
	if err := contractService.EnsureAllRoles(startupCtx); err != nil {
		logger.Warn("pool role reconcile failed; pool writes to affected contracts will revert until re-run", "error", err)
	}

	verifier, err := auth.NewOIDCVerifier(startupCtx, auth.OIDCConfig{
		IssuerURL: cfg.OIDCIssuerURL,
		JWKSURL:   cfg.OIDCJWKSURL,
		Audience:  cfg.OIDCAudience,
	})
	if err != nil {
		return fmt.Errorf("initialize oidc verifier: %w", err)
	}

	server := &http.Server{
		Addr: ":" + cfg.Port,
		Handler: api.New(version, logger, api.Options{
			ReadinessChecks: []api.ReadinessCheck{
				{
					Name:  "database",
					Check: db.Checker{DB: conn}.Check,
				},
				{
					Name:  "ethereum",
					Check: ethClient.Check,
				},
			},
			Contracts:     contractService,
			Loans:         loanService,
			Verifier:      verifier,
			Identities:    identityService,
			SignerAddress: ethClient.SignerAddress().Hex(),
		}),
		ReadHeaderTimeout: 5 * time.Second,
	}

	serverErrors := make(chan error, 1)
	go func() {
		logger.Info("starting API", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrors <- err
		}
	}()

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)
	case <-ctx.Done():
		logger.Info("shutdown signal received")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("graceful shutdown: %w", err)
		}
		logger.Info("shutdown complete")
	}
	return nil
}

func hostname() string {
	name, err := os.Hostname()
	if err != nil || name == "" {
		return "api"
	}
	return name
}
