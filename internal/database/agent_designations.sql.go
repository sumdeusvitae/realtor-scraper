// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: agent_designations.sql

package database

import (
	"context"
	"database/sql"
)

const createAgentDesignation = `-- name: CreateAgentDesignation :exec
INSERT INTO agent_designations (agent_id, designation_id)
VALUES (
    ?,
    ?
)
ON CONFLICT(agent_id, designation_id) DO NOTHING
`

type CreateAgentDesignationParams struct {
	AgentID       sql.NullString
	DesignationID sql.NullInt64
}

func (q *Queries) CreateAgentDesignation(ctx context.Context, arg CreateAgentDesignationParams) error {
	_, err := q.db.ExecContext(ctx, createAgentDesignation, arg.AgentID, arg.DesignationID)
	return err
}
