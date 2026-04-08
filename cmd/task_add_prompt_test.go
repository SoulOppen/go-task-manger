package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestTaskAddPromptHelpSinLLMEnv(t *testing.T) {
	resetRootCmdFlags(t)
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"task", "add-prompt", "--help"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "GTM_LLM") {
		t.Fatalf("help inesperado: %s", out)
	}
}
