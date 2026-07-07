package config

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func clearEnv(t *testing.T) {
	t.Helper()
	for _, k := range []string{
		"PORT",
		"ETH_RPC_URL", "CHAIN_ID", "DATABASE_URL",
		"LOAN_BASE_URI", "DEPLOYER_PRIVATE_KEY", "KEY_ENCRYPTION_MASTER_KEY",
	} {
		t.Setenv(k, "")
	}
	t.Setenv("KEY_ENCRYPTION_MASTER_KEY", strings.Repeat("a", 64))
}

func TestLoadValidationErrors(t *testing.T) {
	cases := map[string]map[string]string{
		"bad port":       {"PORT": "notaport"},
		"port out range": {"PORT": "70000"},
		"bad chain id":   {"CHAIN_ID": "abc"},
		"zero chain id":  {"CHAIN_ID": "0"},
		"bad key":        {"DEPLOYER_PRIVATE_KEY": "0xnothex"},
	}

	for name, overrides := range cases {
		t.Run(name, func(t *testing.T) {
			clearEnv(t)
			for k, v := range overrides {
				t.Setenv(k, v)
			}
			_, err := Load()
			require.Error(t, err)
		})
	}
}

func TestLogValueRedactsKey(t *testing.T) {
	cfg := Config{DeployerPrivateKey: "0x12345"}
	for _, attr := range cfg.LogValue().Group() {
		require.NotEqual(t, "0x12345", attr.Value.String())
	}
}

func TestLogValueRedactsDatabaseURL(t *testing.T) {
	cfg := Config{DatabaseURL: "postgres://user:password@localhost:5432/app"}
	for _, attr := range cfg.LogValue().Group() {
		require.NotContains(t, attr.Value.String(), "password")
	}
}
