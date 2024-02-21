package data

import (
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Models struct {
	Users interface {
		Insert(*User) error
		GetByNickname(string) (*User, error)
		Update(*User) error
		Delete(*User) error
	}
	Tokens interface {
		SetSession(*User, Session) error
	}
	Posts interface {
		Insert(*Post) error
		GetAll(string, []string, Filters) ([]*Post, error)
		GetById(pgtype.UUID) (*Post, error)
		Update(*Post) error
		Delete(pgtype.UUID) error
	}
}

func NewModels(db *pgxpool.Pool) Models {
	return Models{
		Users:  UserModel{DB: db},
		Tokens: TokenModel{DB: db},
		Posts:  PostModel{DB: db},
	}
}
