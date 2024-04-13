package data

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base32"
	"time"

	"github.com/skraio/banner-service/internal/validator"
)

type Token struct {
	Plaintext string `json:"token"`
	Hash      []byte `json:"-"`
	UserID    int64  `json:"-"`
}

func generateToken(userID int64) (*Token, error) {
	token := &Token{
		UserID: userID,
	}

	randomBytes := make([]byte, 16)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]

	return token, nil
}

func ValidateTokenPlaintext(v *validator.Validator, tokenPlaintext string) {
	v.Check(tokenPlaintext != "", "token", "must be provided")
	v.Check(len(tokenPlaintext) == 26, "token", "must be 26 bytes long")
}

type TokenModel struct {
	DB *sql.DB
}

func (m TokenModel) New(userID int64) (*Token, error) {
	token, err := generateToken(userID)
	if err != nil {
		return nil, err
	}

	err = m.Insert(token)
	return token, err
}

func (m TokenModel) Insert(token *Token) error {
	query := `
        INSERT INTO tokens (token_hash, user_id)
        VALUES ($1, $2)`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, token.Hash, token.UserID)
	return err
}

func (m TokenModel) DeleteAll(userID int64) error {
	query := `
        DELETE FROM tokens
        WHERE user_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, userID)
	return err
}
