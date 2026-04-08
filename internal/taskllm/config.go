package taskllm

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// ErrLLMNotConfigured indica que falta configuracion obligatoria para el LLM (API key, proveedor o modelo).
var ErrLLMNotConfigured = errors.New("LLM no esta configurado")

const (
	envProvider = "GTM_LLM_PROVIDER"
	envAPIKey   = "GTM_LLM_API_KEY"
	envModel    = "GTM_LLM_MODEL"
	envBaseURL  = "GTM_LLM_BASE_URL"
)

// Config selecciona el proveedor LLM (API key + modelo + opcional base URL).
// No hay valores por defecto en codigo: proveedor y modelo deben venir del entorno o de flags CLI.
type Config struct {
	Provider string // gemini | openai
	APIKey   string
	Model    string
	BaseURL  string // opcional: vacio usa el endpoint HTTP del driver elegido
}

// ConfigFromEnv lee GTM_LLM_* desde el entorno. Solo exige API key; proveedor y modelo pueden ir vacios y completarse con flags antes de Validate.
func ConfigFromEnv() (Config, error) {
	key := strings.TrimSpace(os.Getenv(envAPIKey))
	if key == "" {
		return Config{}, fmt.Errorf("%w: %s esta vacia o no definida; definala en .env o en el entorno (vea .envExample)", ErrLLMNotConfigured, envAPIKey)
	}
	prov := strings.TrimSpace(os.Getenv(envProvider))
	model := strings.TrimSpace(os.Getenv(envModel))
	base := strings.TrimSpace(os.Getenv(envBaseURL))
	return Config{
		Provider: strings.ToLower(prov),
		APIKey:   key,
		Model:    model,
		BaseURL:  base,
	}, nil
}

// Validate comprueba que proveedor y modelo esten definidos (entorno y/o flags ya fusionados en Config).
func (c *Config) Validate() error {
	var missing []string
	if strings.TrimSpace(c.Provider) == "" {
		missing = append(missing, fmt.Sprintf("proveedor (%s o --llm-provider: gemini|openai)", envProvider))
	}
	if strings.TrimSpace(c.Model) == "" {
		missing = append(missing, fmt.Sprintf("modelo (%s o --llm-model)", envModel))
	}
	if len(missing) == 0 {
		return nil
	}
	joined := strings.Join(missing, "; ")
	return fmt.Errorf("%w: falta %s (vea .envExample)", ErrLLMNotConfigured, joined)
}
