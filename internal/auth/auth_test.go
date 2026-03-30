package auth

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"testing"
	"time"
)

func setupConfigDir(t *testing.T) string {
	t.Helper()

	tempDir := t.TempDir()
	t.Setenv("APPDATA", tempDir)
	t.Setenv("XDG_CONFIG_HOME", tempDir)

	return filepath.Join(tempDir, "task-manager-go", "users.json")
}

func TestRunSignUpAndLoginSuccess(t *testing.T) {
	nowFunc = func() time.Time { return time.Date(2026, 3, 30, 10, 0, 0, 0, time.UTC) }
	t.Cleanup(func() { nowFunc = time.Now })

	usersPath := setupConfigDir(t)

	signUpInput := bytes.NewBufferString("ariel\nsecret123\n")
	signUpOut := &bytes.Buffer{}
	if err := runSignUp(signUpInput, signUpOut); err != nil {
		t.Fatalf("runSignUp returned error: %v", err)
	}

	if _, err := os.Stat(usersPath); err != nil {
		t.Fatalf("expected users file to exist: %v", err)
	}

	loginInput := bytes.NewBufferString("ariel\nsecret123\n")
	loginOut := &bytes.Buffer{}
	if err := runLogin(loginInput, loginOut); err != nil {
		t.Fatalf("runLogin returned error: %v", err)
	}

	users, err := loadUsers()
	if err != nil {
		t.Fatalf("loadUsers returned error: %v", err)
	}
	user, _, exists := findUser(users, "ariel")
	if !exists {
		t.Fatal("expected user ariel to exist")
	}
	if matched, _ := regexp.MatchString("^[0-9a-f]{24}$", user.QuickConnectValue); !matched {
		t.Fatalf("quick connect value must be 24 hex chars, got %q", user.QuickConnectValue)
	}

	qcPath, err := quickConnectFile("ariel")
	if err != nil {
		t.Fatalf("quickConnectFile returned error: %v", err)
	}
	data, err := os.ReadFile(qcPath)
	if err != nil {
		t.Fatalf("expected quick connect file to exist: %v", err)
	}
	var qc QuickConnectFile
	if err := json.Unmarshal(data, &qc); err != nil {
		t.Fatalf("invalid quick connect json: %v", err)
	}
	if qc.Valor != user.QuickConnectValue {
		t.Fatalf("quick connect file value mismatch: %q != %q", qc.Valor, user.QuickConnectValue)
	}
}

func TestRunSignUpDuplicateUser(t *testing.T) {
	setupConfigDir(t)

	firstInput := bytes.NewBufferString("ariel\nsecret123\n")
	if err := runSignUp(firstInput, &bytes.Buffer{}); err != nil {
		t.Fatalf("first signup failed: %v", err)
	}

	secondInput := bytes.NewBufferString("ariel\notra-clave\n")
	if err := runSignUp(secondInput, &bytes.Buffer{}); err == nil {
		t.Fatal("expected duplicate signup to fail")
	}
}

func TestRunLoginFailsWithWrongPassword(t *testing.T) {
	setupConfigDir(t)

	signUpInput := bytes.NewBufferString("ariel\nsecret123\n")
	if err := runSignUp(signUpInput, &bytes.Buffer{}); err != nil {
		t.Fatalf("signup failed: %v", err)
	}

	loginInput := bytes.NewBufferString("ariel\nincorrecta\n")
	if err := runLogin(loginInput, &bytes.Buffer{}); err == nil {
		t.Fatal("expected login with wrong password to fail")
	}
}

func TestRunLoginFailsWhenUserNotFound(t *testing.T) {
	setupConfigDir(t)

	loginInput := bytes.NewBufferString("inexistente\nsecret123\n")
	if err := runLogin(loginInput, &bytes.Buffer{}); err == nil {
		t.Fatal("expected login with missing user to fail")
	}
}

func TestRunLoginResetsQuickConnectOncePerDay(t *testing.T) {
	setupConfigDir(t)

	baseTime := time.Date(2026, 3, 30, 8, 0, 0, 0, time.UTC)
	nowFunc = func() time.Time { return baseTime }
	t.Cleanup(func() { nowFunc = time.Now })

	if err := runSignUp(bytes.NewBufferString("ariel\nsecret123\n"), &bytes.Buffer{}); err != nil {
		t.Fatalf("signup failed: %v", err)
	}

	if err := runLogin(bytes.NewBufferString("ariel\nsecret123\n"), &bytes.Buffer{}); err != nil {
		t.Fatalf("first login failed: %v", err)
	}
	users, _ := loadUsers()
	user, _, _ := findUser(users, "ariel")
	firstValue := user.QuickConnectValue

	nowFunc = func() time.Time { return baseTime.Add(2 * time.Hour) }
	if err := runLogin(bytes.NewBufferString("ariel\nsecret123\n"), &bytes.Buffer{}); err != nil {
		t.Fatalf("second same-day login failed: %v", err)
	}
	users, _ = loadUsers()
	user, _, _ = findUser(users, "ariel")
	if user.QuickConnectValue != firstValue {
		t.Fatal("quick connect should not reset twice in the same day")
	}

	nowFunc = func() time.Time { return baseTime.Add(24 * time.Hour) }
	if err := runLogin(bytes.NewBufferString("ariel\nsecret123\n"), &bytes.Buffer{}); err != nil {
		t.Fatalf("next-day login failed: %v", err)
	}
	users, _ = loadUsers()
	user, _, _ = findUser(users, "ariel")
	if user.QuickConnectValue == firstValue {
		t.Fatal("quick connect must reset on first execution of a new day")
	}
}
