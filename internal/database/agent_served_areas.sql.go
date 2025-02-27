// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: agent_served_areas.sql

package database

import (
	"context"
	"database/sql"
)

const createAgentServedArea = `-- name: CreateAgentServedArea :exec
INSERT INTO agent_served_areas (agent_id, area_id)
VALUES (
    ?,
    ?
)
ON CONFLICT(agent_id, area_id) DO NOTHING
`

type CreateAgentServedAreaParams struct {
	AgentID sql.NullString
	AreaID  sql.NullInt64
}

func (q *Queries) CreateAgentServedArea(ctx context.Context, arg CreateAgentServedAreaParams) error {
	_, err := q.db.ExecContext(ctx, createAgentServedArea, arg.AgentID, arg.AreaID)
	return err
}
