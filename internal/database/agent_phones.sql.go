// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: agent_phones.sql

package database

import (
	"context"
	"database/sql"
)

const createAgentPhone = `-- name: CreateAgentPhone :exec
INSERT INTO agent_phones (agent_id, phone_id)
VALUES (
    ?,
    ?
)
ON CONFLICT(agent_id, phone_id) DO NOTHING
`

type CreateAgentPhoneParams struct {
	AgentID sql.NullString
	PhoneID sql.NullInt64
}

func (q *Queries) CreateAgentPhone(ctx context.Context, arg CreateAgentPhoneParams) error {
	_, err := q.db.ExecContext(ctx, createAgentPhone, arg.AgentID, arg.PhoneID)
	return err
}
