package data

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"time"

	"github.com/skraio/banner-service/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateUserName = errors.New("duplicate username")
)

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

var AnonymousUser = &User{}

type User struct {
	UserID    int64     `json:"user_id"`
	UserName  string    `json:"username"`
	Role      Role      `json:"role"`
	Password  password  `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

type password struct {
	plaintext *string
	hash      []byte
}

func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}

	p.plaintext = &plaintextPassword
	p.hash = hash

	return nil
}

func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 32, "password", "must not be more than 32 bytes long")
}

func ValidateUsername(v *validator.Validator, username string) {
	v.Check(username != "", "username", "must be provided")
	v.Check(len(username) >= 5, "username", "must be at least 5 bytes long")
	v.Check(len(username) <= 20, "username", "must not be more than 20 bytes long")
}

func ValidateUserCredentials(v *validator.Validator, user *User) {
	ValidateUsername(v, user.UserName)

	v.Check(user.Role == RoleUser || user.Role == RoleAdmin, "role", "must be either user or admin")

	if user.Password.plaintext != nil {
		ValidatePasswordPlaintext(v, *user.Password.plaintext)
	}

	if user.Password.hash == nil {
		panic("missing password hash for user")
	}
}

type UserModel struct {
	DB *sql.DB
}

func (m UserModel) Insert(user *User) error {
	query := `
        INSERT INTO users (username, role, password_hash)
        VALUES ($1, $2, $3)
        RETURNING user_id, created_at`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	args := []any{user.UserName, user.Role, user.Password.hash}

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.UserID, &user.CreatedAt)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"`:
			return ErrDuplicateUserName
		default:
			return err
		}
	}

	return nil
}

func (m UserModel) GetByUserName(userName string) (*User, error) {
	query := `
        SELECT user_id, username, role, password_hash, created_at
        FROM users
        WHERE username = $1`

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, userName).Scan(
		&user.UserID,
		&user.UserName,
		&user.Role,
		&user.Password.hash,
		&user.CreatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (m UserModel) Update(user *User) error {
	query := `
        UPDATE users
        SET username = $1, role = $2, password_hash = $3
        WHERE user_id = $3`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	args := []any{user.UserName, user.Role, user.Password.hash, user.UserID}

	err := m.DB.QueryRowContext(ctx, query, args...).Scan()
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"`:
			return ErrDuplicateUserName
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

func (m UserModel) GetByToken(tokenPlainText string) (*User, error) {
	tokenHash := sha256.Sum256([]byte(tokenPlainText))

	query := `
        SELECT users.user_id, users.username, users.role, users.password_hash, users.created_at
        FROM users
        INNER JOIN tokens
        ON users.user_id = tokens.user_id
        WHERE tokens.token_hash = $1`

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, tokenHash[:]).Scan(
		&user.UserID,
		&user.UserName,
		&user.Role,
		&user.Password.hash,
		&user.CreatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}
