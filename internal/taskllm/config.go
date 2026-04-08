package taskllm

import (
	"fmt"
	"os"
	"strings"
)

const (
	envProvider = "GTM_LLM_PROVIDER"
	envAPIKey   = "GTM_LLM_API_KEY"
	envModel    = "GTM_LLM_MODEL"
	envBaseURL  = "GTM_LLM_BASE_URL"
)

// Config selecciona el proveedor LLM (API key + modelo + opcional base URL).
type Config struct {
	Provider string // gemini | openai
	APIKey   string
	Model    string
	BaseURL  string // vacio: URL por defecto del driver
}

// ConfigFromEnv lee GTM_LLM_* desde el entorno.
func ConfigFromEnv() (Config, error) {
	prov := strings.TrimSpace(os.Getenv(envProvider))
	if prov == "" {
		prov = "gemini"
	}
	key := strings.TrimSpace(os.Getenv(envAPIKey))
	if key == "" {
		return Config{}, fmt.Errorf("defina %s en el entorno (.env)", envAPIKey)
	}
	model := strings.TrimSpace(os.Getenv(envModel))
	if model == "" {
		if prov == "gemini" {
			model = "gemini-2.0-flash"
		} else {
			model = "gpt-4o-mini"
		}
	}
	base := strings.TrimSpace(os.Getenv(envBaseURL))
	return Config{
		Provider: strings.ToLower(prov),
		APIKey:   key,
		Model:    model,
		BaseURL:  base,
	}, nil
}
