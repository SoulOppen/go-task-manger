package auth

import (
	"bufio"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username              string `json:"username"`
	PasswordHash          string `json:"password_hash"`
	QuickConnectValue     string `json:"quick_connect_value,omitempty"`
	QuickConnectCreatedAt string `json:"quick_connect_created_at,omitempty"`
	QuickConnectResetDate string `json:"quick_connect_reset_date,omitempty"`
}

type QuickConnectFile struct {
	Valor      string `json:"valor"`
	Creacion   string `json:"creacion"`
	Expiracion string `json:"expiracion"`
	OS         string `json:"os"`
	PCUID      string `json:"pc_uid"`
	Username   string `json:"username"`
}

type Session struct {
	Username string `json:"username"`
	LoginAt  string `json:"login_at"`
	OS       string `json:"os"`
	PCUID    string `json:"pc_uid"`
}

var nowFunc = time.Now

func sessionFile() (string, error) {
	confDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(confDir, "task-manager-go", "session.json"), nil
}

func scanLine(scanner *bufio.Scanner) (string, error) {
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return "", err
		}
		return "", io.EOF
	}
	return strings.TrimSpace(scanner.Text()), nil
}

// RunSignUp registra un usuario en el store (MySQL en produccion).
func RunSignUp(ctx context.Context, store UserStore, in io.Reader, out io.Writer) error {
	scanner := bufio.NewScanner(in)

	fmt.Fprintln(out, "Bienvenido al sistema de conexión")
	fmt.Fprintln(out, "¿Cuál es tu nombre de usuario?")
	username, err := scanLine(scanner)
	if err != nil {
		return fmt.Errorf("no se pudo leer username: %w", err)
	}
	if username == "" {
		return errors.New("el nombre de usuario es obligatorio")
	}

	fmt.Fprintln(out, "¿Cuál es tu clave?")
	password, err := scanLine(scanner)
	if err != nil {
		return fmt.Errorf("no se pudo leer password: %w", err)
	}
	if password == "" {
		return errors.New("la clave es obligatoria")
	}

	if _, err := store.GetByUsername(ctx, username); err == nil {
		return errors.New("el usuario ya existe")
	} else if !errors.Is(err, ErrUserNotFound) {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u := User{
		Username:     username,
		PasswordHash: string(hashedPassword),
	}
	if err := store.Create(ctx, u); err != nil {
		if errors.Is(err, ErrUserExists) {
			return errors.New("el usuario ya existe")
		}
		return err
	}

	fmt.Fprintln(out, "Fuiste registrado con éxito")
	return nil
}

// RunLogin autentica y actualiza quick-connect en el store.
func RunLogin(ctx context.Context, store UserStore, in io.Reader, out io.Writer) error {
	scanner := bufio.NewScanner(in)

	fmt.Fprintln(out, "¿Cuál es tu nombre de usuario?")
	username, err := scanLine(scanner)
	if err != nil {
		return fmt.Errorf("no se pudo leer username: %w", err)
	}
	if username == "" {
		return errors.New("el nombre de usuario es obligatorio")
	}

	fmt.Fprintln(out, "¿Cuál es tu clave?")
	password, err := scanLine(scanner)
	if err != nil {
		return fmt.Errorf("no se pudo leer password: %w", err)
	}
	if password == "" {
		return errors.New("la clave es obligatoria")
	}

	user, err := store.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return errors.New("no existe usuario")
		}
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return errors.New("clave incorrecta")
	}

	updatedUser, err := syncQuickConnect(user)
	if err != nil {
		return err
	}

	if err := store.Update(ctx, updatedUser); err != nil {
		return err
	}

	if err := saveSession(updatedUser.Username); err != nil {
		return err
	}

	fmt.Fprintln(out, "Login exitoso")
	return nil
}

func syncQuickConnect(user User) (User, error) {
	now := nowFunc().UTC()
	today := now.Format("2006-01-02")

	createdAt, err := parseQuickConnectCreatedAt(user.QuickConnectCreatedAt)
	if err != nil {
		createdAt = time.Time{}
	}

	needsReset := user.QuickConnectValue == "" ||
		user.QuickConnectResetDate != today ||
		createdAt.IsZero() ||
		now.After(createdAt.AddDate(0, 1, 0))

	if needsReset {
		value, err := generateHexValue24()
		if err != nil {
			return user, err
		}
		user.QuickConnectValue = value
		user.QuickConnectCreatedAt = now.Format(time.RFC3339)
		user.QuickConnectResetDate = today
		createdAt = now
	}

	filePath, err := quickConnectFile(user.Username)
	if err != nil {
		return user, err
	}

	valid, err := validateQuickConnectFile(filePath, user, now)
	if err != nil {
		return user, err
	}
	if !valid {
		payload := QuickConnectFile{
			Valor:      user.QuickConnectValue,
			Creacion:   user.QuickConnectCreatedAt,
			Expiracion: createdAt.AddDate(0, 1, 0).Format(time.RFC3339),
			OS:         runtime.GOOS,
			PCUID:      machineUID(),
			Username:   user.Username,
		}
		if err := writeQuickConnectFile(filePath, payload); err != nil {
			return user, err
		}
	}

	return user, nil
}

func parseQuickConnectCreatedAt(value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, errors.New("quick connect sin fecha de creacion")
	}
	return time.Parse(time.RFC3339, value)
}

func generateHexValue24() (string, error) {
	b := make([]byte, 12)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func quickConnectFile(username string) (string, error) {
	confDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	safeUsername := strings.ReplaceAll(strings.TrimSpace(username), " ", "_")
	return filepath.Join(confDir, "task-manager-go", "quick_connect_"+safeUsername+".json"), nil
}

func writeQuickConnectFile(path string, payload QuickConnectFile) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func validateQuickConnectFile(path string, user User, now time.Time) (bool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, err
	}

	var payload QuickConnectFile
	if err := json.Unmarshal(data, &payload); err != nil {
		return false, nil
	}

	createdAt, err := parseQuickConnectCreatedAt(payload.Creacion)
	if err != nil {
		return false, nil
	}

	if payload.Valor == "" || payload.Valor != user.QuickConnectValue {
		return false, nil
	}

	if payload.Username != user.Username {
		return false, nil
	}

	if now.After(createdAt.AddDate(0, 1, 0)) {
		return false, nil
	}

	return true, nil
}

func saveSession(username string) error {
	path, err := sessionFile()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	session := Session{
		Username: username,
		LoginAt:  nowFunc().UTC().Format(time.RFC3339),
		OS:       runtime.GOOS,
		PCUID:    machineUID(),
	}
	data, err := json.MarshalIndent(session, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func clearSession() error {
	path, err := sessionFile()
	if err != nil {
		return err
	}
	if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}

func machineUID() string {
	hostname, _ := os.Hostname()
	home, _ := os.UserHomeDir()
	raw := strings.Join([]string{
		hostname,
		home,
		runtime.GOOS,
		runtime.GOARCH,
	}, "|")
	hash := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(hash[:16])
}

func Logout() {
	if err := clearSession(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Fprintln(os.Stdout, "Logout exitoso")
}
