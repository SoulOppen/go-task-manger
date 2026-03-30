package auth

import (
	"context"
	"database/sql"
	"errors"
	"time"

	mysqldriver "github.com/go-sql-driver/mysql"
)

// MySQLUserStore implementa UserStore contra tabla users.
type MySQLUserStore struct {
	db *sql.DB
}

func NewMySQLUserStore(db *sql.DB) *MySQLUserStore {
	return &MySQLUserStore{db: db}
}

func (s *MySQLUserStore) GetByUsername(ctx context.Context, username string) (User, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT username, password_hash,
		       quick_connect_value, quick_connect_created_at, quick_connect_reset_date
		FROM users WHERE username = ?`, username)

	var u User
	var qcVal sql.NullString
	var qcCreated sql.NullTime
	var qcReset sql.NullTime

	err := row.Scan(&u.Username, &u.PasswordHash, &qcVal, &qcCreated, &qcReset)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrUserNotFound
		}
		return User{}, err
	}
	if qcVal.Valid {
		u.QuickConnectValue = qcVal.String
	}
	if qcCreated.Valid {
		u.QuickConnectCreatedAt = qcCreated.Time.UTC().Format(time.RFC3339)
	}
	if qcReset.Valid {
		u.QuickConnectResetDate = qcReset.Time.UTC().Format("2006-01-02")
	}
	return u, nil
}

func (s *MySQLUserStore) Create(ctx context.Context, u User) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO users (username, password_hash, quick_connect_value, quick_connect_created_at, quick_connect_reset_date)
		 VALUES (?, ?, NULL, NULL, NULL)`,
		u.Username, u.PasswordHash,
	)
	if err != nil {
		var me *mysqldriver.MySQLError
		if errors.As(err, &me) && me.Number == 1062 {
			return ErrUserExists
		}
		return err
	}
	return nil
}

func (s *MySQLUserStore) Update(ctx context.Context, u User) error {
	var qcCreated interface{}
	if u.QuickConnectCreatedAt != "" {
		t, err := time.Parse(time.RFC3339, u.QuickConnectCreatedAt)
		if err != nil {
			return err
		}
		qcCreated = t.UTC()
	} else {
		qcCreated = nil
	}

	var qcVal interface{}
	if u.QuickConnectValue != "" {
		qcVal = u.QuickConnectValue
	} else {
		qcVal = nil
	}

	var qcReset interface{}
	if u.QuickConnectResetDate != "" {
		d, err := time.Parse("2006-01-02", u.QuickConnectResetDate)
		if err != nil {
			return err
		}
		qcReset = d.UTC().Format("2006-01-02")
	} else {
		qcReset = nil
	}

	res, err := s.db.ExecContext(ctx,
		`UPDATE users SET password_hash = ?,
			quick_connect_value = ?, quick_connect_created_at = ?, quick_connect_reset_date = ?
		 WHERE username = ?`,
		u.PasswordHash, qcVal, qcCreated, qcReset, u.Username,
	)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrUserNotFound
	}
	return nil
}
