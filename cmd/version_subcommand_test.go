package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/SoulOppen/task-manager-go/internal/config"
)

func TestVersionSubcommand_Output(t *testing.T) {
	resetRootCmdFlags(t)
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"version"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}
	got := strings.TrimSpace(buf.String())
	if got != config.AppVersion() {
		t.Fatalf("got %q, want %q", got, config.AppVersion())
	}
}
