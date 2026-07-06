package config

import (
	"cmp"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"kaleido-project/internal/eth"
	"kaleido-project/internal/keys"
)

// devKeyEncryptionMasterKey is a throwaway dev default matching .env.example, like the deployer key.
const devKeyEncryptionMasterKey = "6a006ea1d0bfd421d93890dfe78ec0fb16e74a9818e5a097a5d5cc0f62693051"

type Config struct {
	Port                     string
	EthRPCURL                string
	ChainID                  int64
	DatabaseURL              string
	LoanBaseURI              string
	DeployerPrivateKey       string
	KeyEncryptionMasterKey   string
	OIDCIssuerURL            string
	OIDCJWKSURL              string
	OIDCAudience             string
	ServicerKeyPoolSize      int
	ReconcileIntervalSeconds int
	// KeyEncryptor selects the signing-key encryption backend: keys.AESGCMScheme (local) or keys.KMSScheme (cloud).
	KeyEncryptor string
	KMSKeyID     string
}

func Load() (Config, error) {
	cfg := Config{
		Port:                   getenv("PORT", "8080"),
		EthRPCURL:              getenv("ETH_RPC_URL", "http://127.0.0.1:31545"),
		DatabaseURL:            getenv("DATABASE_URL", "postgres://loan_notes:loan_notes@127.0.0.1:5432/loan_notes?sslmode=disable"),
		LoanBaseURI:            getenv("LOAN_BASE_URI", "http://localhost:8080/loans/"),
		DeployerPrivateKey:     os.Getenv("DEPLOYER_PRIVATE_KEY"),
		KeyEncryptionMasterKey: getenv("KEY_ENCRYPTION_MASTER_KEY", devKeyEncryptionMasterKey),
		OIDCIssuerURL:          getenv("OIDC_ISSUER_URL", "http://localhost:8081/realms/loan-notes"),
		OIDCAudience:           getenv("OIDC_AUDIENCE", "loan-notes-api"),
		KeyEncryptor:           getenv("KEY_ENCRYPTOR", keys.AESGCMScheme),
		KMSKeyID:               os.Getenv("KMS_KEY_ID"),
	}
	// The JWKS fetch path can differ from the issuer name in tokens (compose service name vs host-published URL).
	cfg.OIDCJWKSURL = getenv("OIDC_JWKS_URL", cfg.OIDCIssuerURL+"/protocol/openid-connect/certs")

	chainID, err := parseChainID(getenv("CHAIN_ID", "1337"))
	if err != nil {
		return Config{}, err
	}
	cfg.ChainID = chainID

	poolSize, err := strconv.Atoi(getenv("SERVICER_KEY_POOL_SIZE", "2"))
	if err != nil || poolSize < 0 || poolSize > 32 {
		return Config{}, fmt.Errorf("SERVICER_KEY_POOL_SIZE must be an integer between 0 and 32, got %q", getenv("SERVICER_KEY_POOL_SIZE", "2"))
	}
	cfg.ServicerKeyPoolSize = poolSize

	reconcileInterval, err := strconv.Atoi(getenv("RECONCILE_INTERVAL_SECONDS", "5"))
	if err != nil || reconcileInterval < 1 || reconcileInterval > 3600 {
		return Config{}, fmt.Errorf("RECONCILE_INTERVAL_SECONDS must be an integer between 1 and 3600, got %q", getenv("RECONCILE_INTERVAL_SECONDS", "5"))
	}
	cfg.ReconcileIntervalSeconds = reconcileInterval

	if err := cfg.validate(); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func (c Config) validate() error {
	if p, err := strconv.Atoi(c.Port); err != nil || p < 1 || p > 65535 {
		return fmt.Errorf("PORT must be a valid TCP port, got %q", c.Port)
	}
	if c.EthRPCURL == "" {
		return fmt.Errorf("ETH_RPC_URL must not be empty")
	}
	if c.ChainID <= 0 {
		return fmt.Errorf("CHAIN_ID must be positive, got %d", c.ChainID)
	}
	if c.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL must not be empty")
	}
	if c.LoanBaseURI == "" {
		return fmt.Errorf("LOAN_BASE_URI must not be empty")
	}
	if c.DeployerPrivateKey != "" {
		if _, err := eth.ParsePrivateKey(c.DeployerPrivateKey); err != nil {
			return fmt.Errorf("DEPLOYER_PRIVATE_KEY is invalid: %w", err)
		}
	}
	switch c.KeyEncryptor {
	case keys.AESGCMScheme:
		if _, err := keys.NewAESGCM(c.KeyEncryptionMasterKey); err != nil {
			return fmt.Errorf("KEY_ENCRYPTION_MASTER_KEY is invalid: %w", err)
		}
	case keys.KMSScheme:
		if c.KMSKeyID == "" {
			return fmt.Errorf("KMS_KEY_ID is required when KEY_ENCRYPTOR is %s", keys.KMSScheme)
		}
	default:
		return fmt.Errorf("KEY_ENCRYPTOR must be %s or %s, got %q", keys.AESGCMScheme, keys.KMSScheme, c.KeyEncryptor)
	}
	if c.OIDCIssuerURL == "" || c.OIDCJWKSURL == "" || c.OIDCAudience == "" {
		return fmt.Errorf("OIDC_ISSUER_URL, OIDC_JWKS_URL, and OIDC_AUDIENCE must not be empty")
	}
	return nil
}

func (c Config) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("port", c.Port),
		slog.String("eth_rpc_url", c.EthRPCURL),
		slog.Int64("chain_id", c.ChainID),
		// Redact URL but log whether it's set
		slog.Bool("database_url_set", c.DatabaseURL != ""),
		slog.String("loan_base_uri", c.LoanBaseURI),
		// Redact private kkey but log whether it's set
		slog.Bool("deployer_key_set", c.DeployerPrivateKey != ""),
		// Redact master key but log whether it's set
		slog.Bool("key_encryption_master_key_set", c.KeyEncryptionMasterKey != ""),
		slog.String("oidc_issuer_url", c.OIDCIssuerURL),
		slog.String("oidc_jwks_url", c.OIDCJWKSURL),
		slog.String("oidc_audience", c.OIDCAudience),
	)
}

func parseChainID(s string) (int64, error) {
	v, err := strconv.ParseInt(strings.TrimSpace(s), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("CHAIN_ID must be an integer, got %q", s)
	}
	return v, nil
}

func getenv(key, fallback string) string {
	return cmp.Or(os.Getenv(key), fallback)
}
