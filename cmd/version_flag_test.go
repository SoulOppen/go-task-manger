package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/SoulOppen/task-manager-go/internal/config"
)

func TestRootVersionShortFlag(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"-v"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}
	got := strings.TrimSpace(buf.String())
	if got != config.AppVersion() {
		t.Fatalf("version output: got %q, want %q", got, config.AppVersion())
	}
}

func TestRootVersionLongFlag(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"--version"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}
	got := strings.TrimSpace(buf.String())
	if got != config.AppVersion() {
		t.Fatalf("version output: got %q, want %q", got, config.AppVersion())
	}
}
