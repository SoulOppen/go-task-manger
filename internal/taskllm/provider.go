package taskllm

import (
	"context"
	"fmt"
	"strings"

	"github.com/SoulOppen/task-manager-go/internal/task"
)

// Provider llama a un backend LLM y devuelve solo el texto JSON esperado.
type Provider interface {
	CompleteJSON(ctx context.Context, system, user string) (string, error)
}

// NewProvider instancia el driver segun cfg.Provider.
func NewProvider(cfg Config) (Provider, error) {
	switch cfg.Provider {
	case "gemini":
		return newGeminiProvider(cfg), nil
	case "openai":
		return newOpenAIProvider(cfg), nil
	default:
		return nil, fmt.Errorf("GTM_LLM_PROVIDER desconocido: %q (use gemini u openai)", cfg.Provider)
	}
}

// ExtractTasksFromPrompt orquesta llamada LLM + parse + build.
func ExtractTasksFromPrompt(ctx context.Context, cfg Config, userText string) ([]*task.Task, error) {
	if strings.TrimSpace(userText) == "" {
		return nil, fmt.Errorf("el texto del usuario esta vacio")
	}
	p, err := NewProvider(cfg)
	if err != nil {
		return nil, err
	}
	raw, err := p.CompleteJSON(ctx, SystemPrompt(), userText)
	if err != nil {
		return nil, err
	}
	items, err := ParseTasksJSON(raw)
	if err != nil {
		return nil, err
	}
	return BuildTasks(items)
}
