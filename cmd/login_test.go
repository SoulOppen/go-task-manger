package cmd

import (
	"bytes"
	"testing"
)

func TestLoginHelpDoesNotRequireMySQL(t *testing.T) {
	t.Setenv("DB_HOST", "")
	t.Setenv("DB_PORT", "")
	t.Setenv("DB_USER", "")
	t.Setenv("DB_NAME", "")

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"login", "--help"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(buf.Bytes(), []byte("signup")) {
		t.Fatalf("unexpected help: %s", buf.String())
	}
}
