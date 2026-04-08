package taskllm

import (
	"errors"
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
}

func TestConfigFromEnv_conAPIKey(t *testing.T) {
	t.Setenv(envAPIKey, "test-key")
	t.Setenv(envProvider, "openai")
	t.Setenv(envModel, "")
	cfg, err := ConfigFromEnv()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.APIKey != "test-key" || cfg.Provider != "openai" {
		t.Fatalf("%+v", cfg)
	}
}
