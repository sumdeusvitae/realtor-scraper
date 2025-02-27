// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: languages.sql

package database

import (
	"context"
	"database/sql"
)

const createLanguage = `-- name: CreateLanguage :one
INSERT INTO languages (name)
VALUES (?)
ON CONFLICT(name) DO NOTHING
RETURNING id
`

func (q *Queries) CreateLanguage(ctx context.Context, name sql.NullString) (int64, error) {
	row := q.db.QueryRowContext(ctx, createLanguage, name)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const getLanguageID = `-- name: GetLanguageID :one
SELECT id FROM languages 
WHERE name = ? 
LIMIT 1
`

func (q *Queries) GetLanguageID(ctx context.Context, name sql.NullString) (int64, error) {
	row := q.db.QueryRowContext(ctx, getLanguageID, name)
	var id int64
	err := row.Scan(&id)
	return id, err
}
