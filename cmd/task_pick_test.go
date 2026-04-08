package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestTaskPickHelpSinMySQL(t *testing.T) {
	resetRootCmdFlags(t)
	t.Setenv("DB_HOST", "")
	t.Setenv("DB_PORT", "")
	t.Setenv("DB_USER", "")
	t.Setenv("DB_NAME", "")

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"task", "pick", "--help"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "al azar") {
		t.Fatalf("help inesperado: %s", buf.String())
	}
}
