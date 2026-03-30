package cmd

import "testing"

// resetRootCmdFlags evita que --help / --version queden activos en el mismo rootCmd entre tests.
func resetRootCmdFlags(t *testing.T) {
	t.Helper()
	for _, name := range []string{"help", "version"} {
		f := rootCmd.Flags().Lookup(name)
		if f == nil {
			continue
		}
		if err := f.Value.Set("false"); err != nil {
			t.Fatalf("flag %s: %v", name, err)
		}
		f.Changed = false
	}
}
