package auth

import (
	"context"
	"errors"
)

var (
	ErrUserNotFound = errors.New("usuario no encontrado")
	ErrUserExists   = errors.New("el usuario ya existe")
)

// UserStore persiste credenciales y campos quick-connect en MySQL o en memoria (tests).
type UserStore interface {
	GetByUsername(ctx context.Context, username string) (User, error)
	Create(ctx context.Context, u User) error
	Update(ctx context.Context, u User) error
}
