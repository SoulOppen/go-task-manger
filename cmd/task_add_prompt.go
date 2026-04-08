package cmd

import (
	"fmt"
	"io"
	"strings"

	"github.com/SoulOppen/task-manager-go/internal/task"
	"github.com/SoulOppen/task-manager-go/internal/taskllm"
	"github.com/spf13/cobra"
)

var (
	addPromptText     string
	addPromptProvider string
	addPromptModel    string
)

var taskAddPromptCmd = &cobra.Command{
	Use:   "add-prompt",
	Short: "Crear tarea(s) desde texto con un LLM (API)",
	Long:  `Envia el texto a un proveedor configurado con GTM_LLM_* y crea una o varias tareas segun el JSON devuelto.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		userText, err := readPromptText(cmd, addPromptText)
		if err != nil {
			return err
		}
		cfg, err := taskllm.ConfigFromEnv()
		if err != nil {
			return err
		}
		if p := strings.TrimSpace(addPromptProvider); p != "" {
			cfg.Provider = strings.ToLower(p)
		}
		if m := strings.TrimSpace(addPromptModel); m != "" {
			cfg.Model = m
		}

		tasks, err := taskllm.ExtractTasksFromPrompt(cmd.Context(), cfg, userText)
		if err != nil {
			return err
		}
		return withTaskRepo(cmd.Context(), func(repo *task.Repository) error {
			for _, t := range tasks {
				if err := repo.Create(cmd.Context(), t); err != nil {
					return fmt.Errorf("al crear tarea %q: %w", t.Name, err)
				}
				fmt.Fprintln(cmd.OutOrStdout(), t.ID)
			}
			return nil
		})
	},
}

func readPromptText(cmd *cobra.Command, flag string) (string, error) {
	if s := strings.TrimSpace(flag); s != "" {
		return s, nil
	}
	b, err := io.ReadAll(cmd.InOrStdin())
	if err != nil {
		return "", err
	}
	if len(strings.TrimSpace(string(b))) == 0 {
		return "", fmt.Errorf("indica --prompt o redirige texto por stdin")
	}
	return string(b), nil
}

func init() {
	taskCmd.AddCommand(taskAddPromptCmd)
	taskAddPromptCmd.Flags().StringVar(&addPromptText, "prompt", "", "texto a enviar al modelo (si vacio, lee stdin)")
	taskAddPromptCmd.Flags().StringVar(&addPromptProvider, "llm-provider", "", "sobrescribe GTM_LLM_PROVIDER (gemini|openai)")
	taskAddPromptCmd.Flags().StringVar(&addPromptModel, "llm-model", "", "sobrescribe GTM_LLM_MODEL")
	taskAddPromptCmd.SilenceUsage = true
}
