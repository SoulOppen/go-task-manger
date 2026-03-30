package config

import (
	"strings"
	"testing"
)

func TestAppNameMatchesDefaultNormalized(t *testing.T) {
	want := strings.ReplaceAll(strings.TrimSpace(DefaultName), " ", "-")
	if got := AppName(); got != want {
		t.Fatalf("AppName() = %q, want %q", got, want)
	}
}

func TestAppVersionMatchesDefaultTrimmed(t *testing.T) {
	want := strings.TrimSpace(DefaultVersion)
	if got := AppVersion(); got != want {
		t.Fatalf("AppVersion() = %q, want %q", got, want)
	}
}
