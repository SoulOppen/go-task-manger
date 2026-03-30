package cmd

import (
	"bytes"
	"testing"
)

func TestSwitchHelpSinMySQL(t *testing.T) {
	resetRootCmdFlags(t)
	t.Setenv("DB_HOST", "")
	t.Setenv("DB_PORT", "")
	t.Setenv("DB_USER", "")
	t.Setenv("DB_NAME", "")

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"switch", "--help"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(buf.Bytes(), []byte("Login como otro usuario")) {
		t.Fatalf("help inesperado: %s", buf.String())
	}
}
