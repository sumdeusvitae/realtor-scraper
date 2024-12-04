// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: addresses.sql

package database

import (
	"context"
	"database/sql"
)

const createAddress = `-- name: CreateAddress :one
INSERT INTO addresses (line, line2, city, country, postal_code, state, state_code)
VALUES (?, ?, ?, ?, ?, ?, ?)
ON CONFLICT (line, line2, city, state_code, postal_code)
    DO NOTHING
RETURNING id
`

type CreateAddressParams struct {
	Line       sql.NullString
	Line2      sql.NullString
	City       sql.NullString
	Country    sql.NullString
	PostalCode sql.NullString
	State      sql.NullString
	StateCode  sql.NullString
}

func (q *Queries) CreateAddress(ctx context.Context, arg CreateAddressParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, createAddress,
		arg.Line,
		arg.Line2,
		arg.City,
		arg.Country,
		arg.PostalCode,
		arg.State,
		arg.StateCode,
	)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const getAddressID = `-- name: GetAddressID :one
SELECT id
FROM addresses
WHERE line = ?
  AND line2 = ?
  AND city = ?
  AND state_code = ?
  AND postal_code = ?
LIMIT 1
`

type GetAddressIDParams struct {
	Line       sql.NullString
	Line2      sql.NullString
	City       sql.NullString
	StateCode  sql.NullString
	PostalCode sql.NullString
}

func (q *Queries) GetAddressID(ctx context.Context, arg GetAddressIDParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, getAddressID,
		arg.Line,
		arg.Line2,
		arg.City,
		arg.StateCode,
		arg.PostalCode,
	)
	var id int64
	err := row.Scan(&id)
	return id, err
}
