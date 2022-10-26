// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: user-profiles.sql

package sqlc

import (
	"context"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

const getProfiles = `-- name: getProfiles :many
SELECT id, user_id, name
FROM users.profiles
WHERE user_id = ANY ($1::varchar[])
`

func (q *Queries) getProfiles(ctx context.Context, dollar_1 []string) ([]UsersProfile, error) {
	rows, err := q.db.QueryContext(ctx, getProfiles, pq.Array(dollar_1))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []UsersProfile
	for rows.Next() {
		var i UsersProfile
		if err := rows.Scan(&i.ID, &i.UserID, &i.Name); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const saveProfile = `-- name: saveProfile :exec
INSERT INTO users.profiles (id, user_id, name)
VALUES ($1::uuid, $2::varchar, $3::varchar)
ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name
`

type saveProfileParams struct {
	Column1 uuid.UUID `db:"column_1" json:"column1"`
	Column2 string    `db:"column_2" json:"column2"`
	Column3 string    `db:"column_3" json:"column3"`
}

func (q *Queries) saveProfile(ctx context.Context, arg saveProfileParams) error {
	_, err := q.db.ExecContext(ctx, saveProfile, arg.Column1, arg.Column2, arg.Column3)
	return err
}