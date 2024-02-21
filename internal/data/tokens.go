package data

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

// TODO: userSignIn handler --> userSignInService (returns tokens) --> CreateSession (returns tokens) --> insert into db
// TODO:					 getByCredentials, createSession				newAccess newRefresh
// TODO: find user using refresh token
// TODO: find user ID
// TODO: create newRefreshToken ; NewAccessToken
// TODO: create session
// TODO: insert/set session
// TODO: refresh tokens
// TODO: parse accessToken into userID
//

type Session struct {
	RefreshToken string    `json:"refreshToken" bson:"refreshToken"`
	ExpiredAt    time.Time `json:"expiresAt" bson:"expiresAt"`
}

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

func (m TokenModel) SetSession(user *User, session Session) error {
	//TODO: implement sessions into users table
	query := `UPDATE users 
	SET session = $1
	WHERE id = $2
	RETURNING version`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRow(ctx, query, session, user.ID).Scan(&user.Version)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
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
