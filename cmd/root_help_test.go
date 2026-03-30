package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestRootHelp_ListaComandos(t *testing.T) {
	resetRootCmdFlags(t)
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"--help"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	for _, needle := range []string{"login", "logout", "switch", "task", "version", "Gestionar tareas"} {
		if !strings.Contains(out, needle) {
			t.Fatalf("falta %q en help: %s", needle, out)
		}
	}
}
