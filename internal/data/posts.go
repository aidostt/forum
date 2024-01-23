package data

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type Post struct {
	ID          pgtype.UUID `form:"id"`
	AuthorID    pgtype.UUID `form:"author_id"`
	Heading     string      `form:"heading"`
	Description string      `form:"description"`
	Tags        []string    `form:"tags"`
	Version     int         `form:"-"`
	CreatedAt   time.Time   `form:"created_at"`
}

//TODO:make migrations files for posts
//TODO:create CRUD operations for posts

type PostModel struct {
	DB *pgxpool.Pool
}

func (m PostModel) Insert(post *Post) error {
	query := `INSERT INTO posts (author_id, heading, description, tags) 
VALUES ($1, $2, $3, $4) RETURNING id, version, created_at`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	args := []any{post.AuthorID, post.Heading, post.Description, post.Tags}
	err := m.DB.QueryRow(ctx, query, args...).Scan(&post.ID, &post.Version, &post.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (m PostModel) GeAll(heading string, tags []string, filters Filters) ([]*Post, error) {
	return nil, nil
}

func (m PostModel) GetById(id pgtype.UUID) (*Post, error) {
	query := `SELECT * FROM posts WHERE id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var post Post
	err := m.DB.QueryRow(ctx, query, id).Scan(
		&post.ID,
		&post.AuthorID,
		&post.Heading,
		&post.Description,
		&post.Tags,
		&post.Version,
		&post.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}

	return &post, nil
}

func (m PostModel) Update(post *Post) error {
	query := `UPDATE posts SET heading = $1, description = $2, tags = $3, version = version + 1, created_at = now()
	WHERE id = $4 AND version = $5
	RETURNING version`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	args := []any{
		post.Heading,
		post.Description,
		post.Tags,
		post.ID,
		post.Version,
	}

	err := m.DB.QueryRow(ctx, query, args...).Scan(&post.Version)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrEditConflict
		}
		return err
	}
	return nil
}

func (m PostModel) Delete(post *Post) error {
	query := `DELETE FROM posts WHERE id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.Exec(ctx, query, post.ID)
	return err
}
