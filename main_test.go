package main

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/SoulOppen/task-manager-go/internal/config"
)

// Contrato de salida: el binario responde -v como en internal/config.
func TestBinary_VersionFlag_consistente(t *testing.T) {
	if testing.Short() {
		t.Skip("omite go run en -short")
	}
	out, err := exec.Command("go", "run", ".", "-v").Output()
	if err != nil {
		if exit, ok := err.(*exec.ExitError); ok {
			t.Fatalf("go run: %v\n%s", err, string(exit.Stderr))
		}
		t.Fatal(err)
	}
	got := strings.TrimSpace(string(out))
	if got != config.AppVersion() {
		t.Fatalf("got %q, want %q", got, config.AppVersion())
	}
}
