package taskllm

import (
	"errors"
	"strings"
	"testing"
)

func TestConfigFromEnv_sinAPIKey(t *testing.T) {
	t.Setenv(envAPIKey, "")
	t.Setenv(envProvider, "")
	t.Setenv(envModel, "")
	_, err := ConfigFromEnv()
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, ErrLLMNotConfigured) {
		t.Fatalf("want ErrLLMNotConfigured, got %v", err)
	}
	if !strings.Contains(err.Error(), envAPIKey) {
		t.Fatalf("mensaje debe mencionar la variable: %v", err)
	}
}

func TestConfigFromEnv_conAPIKey(t *testing.T) {
	t.Setenv(envAPIKey, "test-key")
	t.Setenv(envProvider, "openai")
	t.Setenv(envModel, "gpt-4o-mini")
	cfg, err := ConfigFromEnv()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.APIKey != "test-key" || cfg.Provider != "openai" || cfg.Model != "gpt-4o-mini" {
		t.Fatalf("%+v", cfg)
	}
	if err := cfg.Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestConfig_Validate_sinProveedor(t *testing.T) {
	c := Config{Provider: "", APIKey: "k", Model: "m"}
	err := c.Validate()
	if err == nil || !errors.Is(err, ErrLLMNotConfigured) {
		t.Fatalf("want ErrLLMNotConfigured, got %v", err)
	}
	if !strings.Contains(err.Error(), envProvider) {
		t.Fatalf("mensaje debe mencionar %s: %v", envProvider, err)
	}
}

func TestConfig_Validate_sinModelo(t *testing.T) {
	c := Config{Provider: "gemini", APIKey: "k", Model: ""}
	err := c.Validate()
	if err == nil || !errors.Is(err, ErrLLMNotConfigured) {
		t.Fatalf("want ErrLLMNotConfigured, got %v", err)
	}
	if !strings.Contains(err.Error(), envModel) {
		t.Fatalf("mensaje debe mencionar %s: %v", envModel, err)
	}
}

func TestConfig_Validate_sinProveedorNiModelo(t *testing.T) {
	c := Config{Provider: "", APIKey: "k", Model: ""}
	err := c.Validate()
	if err == nil || !errors.Is(err, ErrLLMNotConfigured) {
		t.Fatalf("want ErrLLMNotConfigured, got %v", err)
	}
	if !strings.Contains(err.Error(), envProvider) || !strings.Contains(err.Error(), envModel) {
		t.Fatalf("mensaje debe listar proveedor y modelo: %v", err)
	}
}

func TestConfig_Validate_OK(t *testing.T) {
	c := Config{Provider: "gemini", APIKey: "k", Model: "gemini-2.0-flash"}
	if err := c.Validate(); err != nil {
		t.Fatal(err)
	}
}
