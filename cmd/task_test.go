package cmd

import (
	"bytes"
	"testing"
)

func TestTaskListHelpDoesNotRequireMySQL(t *testing.T) {
	resetRootCmdFlags(t)
	t.Setenv("DB_HOST", "")
	t.Setenv("DB_PORT", "")
	t.Setenv("DB_USER", "")
	t.Setenv("DB_NAME", "")

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"task", "list", "--help"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(buf.Bytes(), []byte("Listar tareas")) {
		t.Fatalf("unexpected help output: %s", buf.String())
	}
}
