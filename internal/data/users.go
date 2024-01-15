package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var (
	ErrDuplicateEmail    = errors.New("duplicate email")
	ErrDuplicateToken    = errors.New("duplicate token")
	ErrDuplicateNickname = errors.New("duplicate nickname")
	ErrNoRecord          = errors.New("no record found")
	ErrEditConflict      = errors.New("edit conflict")
)

type User struct {
	ID        pgtype.UUID `form:"id"`
	Name      string      `form:"name"`
	Nickname  string      `form:"nickname"`
	Email     string      `form:"email"`
	Password  password    `form:"-"`
	Activated bool        `form:"activated"`
	Version   int         `form:"-"`
	CreatedAt time.Time   `form:"created_at"`
}

type password struct {
	plainText *string
	hashed    []byte
}

type UserModel struct {
	DB *pgxpool.Pool
}

func (m UserModel) Insert(user *User) error {
	query := `
	INSERT INTO users (name, nickname, email, password_hash, activated) VALUES ($1, $2, $3, $4, $5)
	RETURNING id, created_at, version;
	`
	args := []any{user.Name, user.Nickname, user.Email, user.Password.hashed, user.Activated}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRow(ctx, query, args...).Scan(&user.ID, &user.CreatedAt, &user.Version)
	if err != nil {
		var pgErr *pgconn.PgError
		fmt.Println(err)
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return ErrDuplicateEmail
			}
		}
		return err
	}
	return nil
}

func (m UserModel) GetByNickname(nickname string) (*User, error) {
	query := `
	SELECT * FROM users 
	         WHERE nickname = $1;
`
	var user User
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRow(ctx, query, nickname).Scan(
		&user.ID,
		&user.Name,
		&user.Nickname,
		&user.Email,
		&user.Password.hashed,
		&user.Activated,
		&user.Version,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}
	return &user, nil
}

func (m UserModel) Update(user *User) error {
	query := `
	UPDATE users
	SET name = $1, nickname = $2, email = $3, password_hash = $4, activated = $5, version = version + 1
	WHERE id = $6 AND version = $7
	RETURNING version
`
	//TODO: discuss the question "with which factor we will update data in db (id, nickname, email)"
	args := []any{
		user.Name,
		user.Nickname,
		user.Email,
		user.Password.hashed,
		user.Activated,
		user.ID,
		user.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRow(ctx, query, args...).Scan(&user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		case err.Error() == `pq: duplicate key value violates unique constraint "users_nickname_key"`:
			return ErrDuplicateNickname
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}
func (m UserModel) Delete(user *User) error {
	return nil
}

func (p *password) Set(plaintTextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintTextPassword), 12)
	if err != nil {
		return err
	}
	p.plainText = &plaintTextPassword
	p.hashed = hash
	return nil
}
