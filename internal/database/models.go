// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package database

import (
	"time"
)

type Request struct {
	ID             int64
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Offset         int64
	ResultsPerPage int64
}
