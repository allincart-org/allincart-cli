package shop

import (
	"testing"
)

func Test_newShopCredentials_envIntegration(t *testing.T) {
	t.Setenv("ALLINCART_CLI_API_CLIENT_ID", "id")
	t.Setenv("ALLINCART_CLI_API_CLIENT_SECRET", "secret")

	cfg := &Config{}
	creds, err := newShopCredentials(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if creds == nil {
		t.Fatal("expected credentials, got nil")
	}
}

func Test_newShopCredentials_envPassword(t *testing.T) {
	t.Setenv("ALLINCART_CLI_API_USERNAME", "user")
	t.Setenv("ALLINCART_CLI_API_PASSWORD", "pass")

	cfg := &Config{}
	creds, err := newShopCredentials(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if creds == nil {
		t.Fatal("expected credentials, got nil")
	}
}

func Test_newShopCredentials_configPassword(t *testing.T) {
	// No env vars needed for this test, but ensure clean state
	t.Setenv("ALLINCART_CLI_API_CLIENT_ID", "")
	t.Setenv("ALLINCART_CLI_API_CLIENT_SECRET", "")
	t.Setenv("ALLINCART_CLI_API_USERNAME", "")
	t.Setenv("ALLINCART_CLI_API_PASSWORD", "")

	cfg := &Config{
		AdminApi: &ConfigAdminApi{
			Username: "user",
			Password: "pass",
		},
	}
	creds, err := newShopCredentials(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if creds == nil {
		t.Fatal("expected credentials, got nil")
	}
}

func Test_newShopCredentials_configIntegration(t *testing.T) {
	// No env vars needed for this test, but ensure clean state
	t.Setenv("ALLINCART_CLI_API_CLIENT_ID", "")
	t.Setenv("ALLINCART_CLI_API_CLIENT_SECRET", "")
	t.Setenv("ALLINCART_CLI_API_USERNAME", "")
	t.Setenv("ALLINCART_CLI_API_PASSWORD", "")

	cfg := &Config{
		AdminApi: &ConfigAdminApi{
			ClientId:     "id",
			ClientSecret: "secret",
		},
	}
	creds, err := newShopCredentials(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if creds == nil {
		t.Fatal("expected credentials, got nil")
	}
}

func Test_newShopCredentials_noConfig(t *testing.T) {
	// Ensure clean env state
	t.Setenv("ALLINCART_CLI_API_CLIENT_ID", "")
	t.Setenv("ALLINCART_CLI_API_CLIENT_SECRET", "")
	t.Setenv("ALLINCART_CLI_API_USERNAME", "")
	t.Setenv("ALLINCART_CLI_API_PASSWORD", "")

	cfg := &Config{}
	_, err := newShopCredentials(cfg)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func Test_NewShopClient(t *testing.T) {
	t.Setenv("ALLINCART_CLI_API_URL", "http://localhost")
	t.Setenv("ALLINCART_CLI_API_CLIENT_ID", "id")
	t.Setenv("ALLINCART_CLI_API_CLIENT_SECRET", "secret")

	cfg := &Config{}
	client, err := NewShopClient(t.Context(), cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if client == nil {
		t.Fatal("expected client, got nil")
	}
}

func Test_NewShopClient_configUrl(t *testing.T) {
	// Ensure env URL is not set
	t.Setenv("ALLINCART_CLI_API_URL", "")
	t.Setenv("ALLINCART_CLI_API_CLIENT_ID", "id")
	t.Setenv("ALLINCART_CLI_API_CLIENT_SECRET", "secret")

	cfg := &Config{URL: "http://config-url"}
	client, err := NewShopClient(t.Context(), cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if client == nil {
		t.Fatal("expected client, got nil")
	}
	// Ideally, we'd check the client's configured URL here, but the SDK doesn't expose it easily.
}

func Test_NewShopClient_skipSSLCheck_env(t *testing.T) {
	t.Setenv("ALLINCART_CLI_API_URL", "http://localhost")
	t.Setenv("ALLINCART_CLI_API_CLIENT_ID", "id")
	t.Setenv("ALLINCART_CLI_API_CLIENT_SECRET", "secret")
	t.Setenv("ALLINCART_CLI_API_DISABLE_SSL_CHECK", "true")

	cfg := &Config{}
	client, err := NewShopClient(t.Context(), cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if client == nil {
		t.Fatal("expected client, got nil")
	}
	// We cannot easily assert the TLS config here without reflection or modifying the original code.
}

func Test_NewShopClient_skipSSLCheck_config(t *testing.T) {
	t.Setenv("ALLINCART_CLI_API_URL", "http://localhost")
	t.Setenv("ALLINCART_CLI_API_CLIENT_ID", "id")
	t.Setenv("ALLINCART_CLI_API_CLIENT_SECRET", "secret")
	// Ensure env var is not set or false
	t.Setenv("ALLINCART_CLI_API_DISABLE_SSL_CHECK", "false")

	cfg := &Config{
		AdminApi: &ConfigAdminApi{
			DisableSSLCheck: true,
		},
	}
	client, err := NewShopClient(t.Context(), cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if client == nil {
		t.Fatal("expected client, got nil")
	}
	// We cannot easily assert the TLS config here without reflection or modifying the original code.
}

func Test_NewShopClient_NoURL(t *testing.T) {
	// Ensure env URL is not set
	t.Setenv("ALLINCART_CLI_API_URL", "")
	t.Setenv("ALLINCART_CLI_API_CLIENT_ID", "id")
	t.Setenv("ALLINCART_CLI_API_CLIENT_SECRET", "secret")

	// Config with empty URL
	cfg := &Config{URL: ""}
	_, err := NewShopClient(t.Context(), cfg)
	// The current implementation doesn't check for empty URL
	// The error would come from the SDK later when it tries to make a request
	if err != nil {
		t.Fatalf("expected no error for empty URL at this stage, got %v", err)
	}
}

func Test_NewShopClient_CredentialsError(t *testing.T) {
	// Ensure clean env state
	t.Setenv("ALLINCART_CLI_API_URL", "http://localhost")
	t.Setenv("ALLINCART_CLI_API_CLIENT_ID", "")
	t.Setenv("ALLINCART_CLI_API_CLIENT_SECRET", "")
	t.Setenv("ALLINCART_CLI_API_USERNAME", "")
	t.Setenv("ALLINCART_CLI_API_PASSWORD", "")

	// Config without credentials
	cfg := &Config{}
	_, err := NewShopClient(t.Context(), cfg)
	if err == nil {
		t.Fatal("expected error for missing credentials, got nil")
	}
}
