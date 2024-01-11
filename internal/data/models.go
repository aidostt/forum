package data

import "github.com/jackc/pgx/v5/pgxpool"

type Models struct {
	Users interface {
		Insert(*User) error
		GetByNickName(string) (*User, error)
		Update(*User) error
		Delete(*User) error
	}
}

func NewModels(db *pgxpool.Pool) Models {
	return Models{
		Users: UserModel{DB: db},
	}
}
