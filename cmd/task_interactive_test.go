package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestResolveTaskID_FromArgs(t *testing.T) {
	id, err := resolveTaskID([]string{"abc-123"}, strings.NewReader(""), &bytes.Buffer{})
	if err != nil {
		t.Fatal(err)
	}
	if id != "abc-123" {
		t.Fatalf("id=%q", id)
	}
}

func TestResolveTaskID_FromPrompt(t *testing.T) {
	id, err := resolveTaskID(nil, strings.NewReader("id-999\n"), &bytes.Buffer{})
	if err != nil {
		t.Fatal(err)
	}
	if id != "id-999" {
		t.Fatalf("id=%q", id)
	}
}

func TestPromptAddFields(t *testing.T) {
	in := strings.NewReader("Tarea X\nDescripcion X\n8\n2026-04-30\n")
	name, desc, rel, due, err := promptAddFields(in, &bytes.Buffer{})
	if err != nil {
		t.Fatal(err)
	}
	if name != "Tarea X" || desc != "Descripcion X" || rel != 8 || due != "2026-04-30" {
		t.Fatalf("unexpected values: %q %q %d %q", name, desc, rel, due)
	}
}

