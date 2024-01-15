package data

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type Token struct {
	PlainText string      `form:"token"`
	Hash      []byte      `form:"-"`
	UserID    pgtype.UUID `form:"-"`
	Expiry    time.Time   `form:"expiry"`
	Scope     string      `form:"-"`
}

type TokenModel struct {
	DB *pgxpool.Pool
}

func GenerateNewToken(userID pgtype.UUID, ttl time.Duration, scope string) (*Token, error) {
	token := &Token{
		UserID: userID,
		Scope:  scope,
		Expiry: time.Now().Add(ttl),
	}
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}
	token.PlainText = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	hash := sha256.Sum256([]byte(token.PlainText))
	token.Hash = hash[:]
	return token, nil
}

func (m *TokenModel) New(userID pgtype.UUID, ttl time.Duration, scope string) (*Token, error) {
	token, err := GenerateNewToken(userID, ttl, scope)
	if err != nil {
		return nil, err
	}
	err = m.Insert(token)
	if err != nil {
		return nil, err
	}
	return token, nil
}
func (m TokenModel) Insert(token *Token) error {
	query := `INSERT INTO tokens (hash, user_id, expiry, scope) VALUES ($1, $2, $3, $4)`

	args := []any{token.Hash, token.UserID, token.Expiry, token.Scope}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.Exec(ctx, query, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		fmt.Println(err)
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return ErrDuplicateToken
			}
		}
		return err
	}
	return nil
}

func (m TokenModel) DeleteAllForUser(scope string, userID pgtype.UUID) error {
	query := `DELETE FROM tokens WHERE scope = $1 AND user_id = $2`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.Exec(ctx, query, scope, userID)
	return err
}
