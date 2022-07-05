// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0
// source: languages.sql

package sqlc

import (
	"context"
)

const getLanguageKeys = `-- name: GetLanguageKeys :many
SELECT code FROM languages
`

func (q *Queries) GetLanguageKeys(ctx context.Context) ([]string, error) {
	rows, err := q.db.QueryContext(ctx, getLanguageKeys)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var code string
		if err := rows.Scan(&code); err != nil {
			return nil, err
		}
		items = append(items, code)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getLanguages = `-- name: GetLanguages :many
SELECT code, legacy_2_letter_code, legacy_3_letter_code, name FROM languages
`

func (q *Queries) GetLanguages(ctx context.Context) ([]Language, error) {
	rows, err := q.db.QueryContext(ctx, getLanguages)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Language
	for rows.Next() {
		var i Language
		if err := rows.Scan(
			&i.Code,
			&i.Legacy2LetterCode,
			&i.Legacy3LetterCode,
			&i.Name,
		); err != nil {
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