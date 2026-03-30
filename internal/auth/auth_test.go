package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"testing"
	"time"
)

type memoryUserStore struct {
	users map[string]User
}

func newMemoryUserStore() *memoryUserStore {
	return &memoryUserStore{users: make(map[string]User)}
}

func (m *memoryUserStore) GetByUsername(ctx context.Context, username string) (User, error) {
	u, ok := m.users[username]
	if !ok {
		return User{}, ErrUserNotFound
	}
	return u, nil
}

func (m *memoryUserStore) Create(ctx context.Context, u User) error {
	if _, ok := m.users[u.Username]; ok {
		return ErrUserExists
	}
	m.users[u.Username] = u
	return nil
}

func (m *memoryUserStore) Update(ctx context.Context, u User) error {
	if _, ok := m.users[u.Username]; !ok {
		return ErrUserNotFound
	}
	m.users[u.Username] = u
	return nil
}

func setupConfigDir(t *testing.T) {
	t.Helper()
	tempDir := t.TempDir()
	t.Setenv("APPDATA", tempDir)
	t.Setenv("XDG_CONFIG_HOME", tempDir)
}

func sessionPathFromTemp(tempDir string) string {
	return filepath.Join(tempDir, "task-manager-go", "session.json")
}

func TestRunSignUpAndLoginSuccess(t *testing.T) {
	nowFunc = func() time.Time { return time.Date(2026, 3, 30, 10, 0, 0, 0, time.UTC) }
	t.Cleanup(func() { nowFunc = time.Now })

	tempDir := t.TempDir()
	t.Setenv("APPDATA", tempDir)
	t.Setenv("XDG_CONFIG_HOME", tempDir)
	sessionPath := filepath.Join(tempDir, "task-manager-go", "session.json")

	store := newMemoryUserStore()
	ctx := context.Background()
	signUpInput := bytes.NewBufferString("ariel\nsecret123\n")
	signUpOut := &bytes.Buffer{}
	if err := RunSignUp(ctx, store, signUpInput, signUpOut); err != nil {
		t.Fatalf("RunSignUp returned error: %v", err)
	}

	loginInput := bytes.NewBufferString("ariel\nsecret123\n")
	loginOut := &bytes.Buffer{}
	if err := RunLogin(ctx, store, loginInput, loginOut); err != nil {
		t.Fatalf("RunLogin returned error: %v", err)
	}
	if _, err := os.Stat(sessionPath); err != nil {
		t.Fatalf("expected session file to exist: %v", err)
	}

	user, err := store.GetByUsername(ctx, "ariel")
	if err != nil {
		t.Fatal(err)
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
	ctx := context.Background()
	store := newMemoryUserStore()
	setupConfigDir(t)

	if err := RunSignUp(ctx, store, bytes.NewBufferString("ariel\nsecret123\n"), &bytes.Buffer{}); err != nil {
		t.Fatalf("first signup failed: %v", err)
	}
	if err := RunSignUp(ctx, store, bytes.NewBufferString("ariel\notra-clave\n"), &bytes.Buffer{}); err == nil {
		t.Fatal("expected duplicate signup to fail")
	}
}

func TestRunLoginFailsWithWrongPassword(t *testing.T) {
	ctx := context.Background()
	store := newMemoryUserStore()
	setupConfigDir(t)

	if err := RunSignUp(ctx, store, bytes.NewBufferString("ariel\nsecret123\n"), &bytes.Buffer{}); err != nil {
		t.Fatalf("signup failed: %v", err)
	}
	if err := RunLogin(ctx, store, bytes.NewBufferString("ariel\nincorrecta\n"), &bytes.Buffer{}); err == nil {
		t.Fatal("expected login with wrong password to fail")
	}
}

func TestRunLoginFailsWhenUserNotFound(t *testing.T) {
	ctx := context.Background()
	store := newMemoryUserStore()
	setupConfigDir(t)

	if err := RunLogin(ctx, store, bytes.NewBufferString("inexistente\nsecret123\n"), &bytes.Buffer{}); err == nil {
		t.Fatal("expected login with missing user to fail")
	}
}

func TestRunLoginResetsQuickConnectOncePerDay(t *testing.T) {
	ctx := context.Background()
	store := newMemoryUserStore()
	setupConfigDir(t)

	baseTime := time.Date(2026, 3, 30, 8, 0, 0, 0, time.UTC)
	nowFunc = func() time.Time { return baseTime }
	t.Cleanup(func() { nowFunc = time.Now })

	if err := RunSignUp(ctx, store, bytes.NewBufferString("ariel\nsecret123\n"), &bytes.Buffer{}); err != nil {
		t.Fatalf("signup failed: %v", err)
	}
	if err := RunLogin(ctx, store, bytes.NewBufferString("ariel\nsecret123\n"), &bytes.Buffer{}); err != nil {
		t.Fatalf("first login failed: %v", err)
	}
	user, _ := store.GetByUsername(ctx, "ariel")
	firstValue := user.QuickConnectValue

	nowFunc = func() time.Time { return baseTime.Add(2 * time.Hour) }
	if err := RunLogin(ctx, store, bytes.NewBufferString("ariel\nsecret123\n"), &bytes.Buffer{}); err != nil {
		t.Fatalf("second same-day login failed: %v", err)
	}
	user, _ = store.GetByUsername(ctx, "ariel")
	if user.QuickConnectValue != firstValue {
		t.Fatal("quick connect should not reset twice in the same day")
	}

	nowFunc = func() time.Time { return baseTime.Add(24 * time.Hour) }
	if err := RunLogin(ctx, store, bytes.NewBufferString("ariel\nsecret123\n"), &bytes.Buffer{}); err != nil {
		t.Fatalf("next-day login failed: %v", err)
	}
	user, _ = store.GetByUsername(ctx, "ariel")
	if user.QuickConnectValue == firstValue {
		t.Fatal("quick connect must reset on first execution of a new day")
	}
}

func TestClearSession(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv("APPDATA", tempDir)
	t.Setenv("XDG_CONFIG_HOME", tempDir)

	if err := saveSession("ariel"); err != nil {
		t.Fatalf("saveSession returned error: %v", err)
	}
	if err := clearSession(); err != nil {
		t.Fatalf("clearSession returned error: %v", err)
	}
	path := sessionPathFromTemp(tempDir)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("expected session file to be deleted, got err=%v", err)
	}
}
