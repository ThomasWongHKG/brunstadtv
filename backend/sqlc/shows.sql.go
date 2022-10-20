// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: shows.sql

package sqlc

import (
	"context"
	"encoding/json"
	"time"

	"github.com/lib/pq"
	"github.com/tabbed/pqtype"
	null_v4 "gopkg.in/guregu/null.v4"
)

const getPermissionsForShows = `-- name: getPermissionsForShows :many
SELECT sh.id,
       access.published::boolean          AS published,
       access.available_from::timestamp   AS available_from,
       access.available_to::timestamp     AS available_to,
       roles.roles::varchar[]             AS usergroups,
       roles.roles_download::varchar[]    AS usergroups_downloads,
       roles.roles_earlyaccess::varchar[] AS usergroups_earlyaccess
FROM shows sh
         LEFT JOIN show_availability access ON access.id = sh.id
         LEFT JOIN show_roles roles ON roles.id = sh.id
WHERE sh.id = ANY ($1::int[])
`

type getPermissionsForShowsRow struct {
	ID                    int32     `db:"id" json:"id"`
	Published             bool      `db:"published" json:"published"`
	AvailableFrom         time.Time `db:"available_from" json:"availableFrom"`
	AvailableTo           time.Time `db:"available_to" json:"availableTo"`
	Usergroups            []string  `db:"usergroups" json:"usergroups"`
	UsergroupsDownloads   []string  `db:"usergroups_downloads" json:"usergroupsDownloads"`
	UsergroupsEarlyaccess []string  `db:"usergroups_earlyaccess" json:"usergroupsEarlyaccess"`
}

func (q *Queries) getPermissionsForShows(ctx context.Context, dollar_1 []int32) ([]getPermissionsForShowsRow, error) {
	rows, err := q.db.QueryContext(ctx, getPermissionsForShows, pq.Array(dollar_1))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []getPermissionsForShowsRow
	for rows.Next() {
		var i getPermissionsForShowsRow
		if err := rows.Scan(
			&i.ID,
			&i.Published,
			&i.AvailableFrom,
			&i.AvailableTo,
			pq.Array(&i.Usergroups),
			pq.Array(&i.UsergroupsDownloads),
			pq.Array(&i.UsergroupsEarlyaccess),
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

const getShows = `-- name: getShows :many
WITH ts AS (SELECT shows_id,
                   json_object_agg(languages_code, title)       AS title,
                   json_object_agg(languages_code, description) AS description
            FROM shows_translations
            GROUP BY shows_id),
     tags AS (SELECT shows_id,
                     array_agg(tags_id) AS tags
              FROM shows_tags
              GROUP BY shows_id),
     images AS (WITH images AS (SELECT show_id, style, language, filename_disk
                                FROM images img
                                         JOIN directus_files df on img.file = df.id)
                SELECT show_id, json_agg(images) as json
                FROM images
                GROUP BY show_id)
SELECT sh.id,
       sh.legacy_id,
       sh.type,
       fs.filename_disk            as image_file_name,
       tags.tags::int[]            AS tag_ids,
       COALESCE(images.json, '[]') as images,
       ts.title,
       ts.description
FROM shows sh
         LEFT JOIN tags ON tags.shows_id = sh.id
         LEFT JOIN ts ON sh.id = ts.shows_id
         LEFT JOIN images ON sh.id = images.show_id
         LEFT JOIN directus_files fs ON fs.id = sh.image_file_id
WHERE sh.id = ANY ($1::int[])
`

type getShowsRow struct {
	ID            int32                 `db:"id" json:"id"`
	LegacyID      null_v4.Int           `db:"legacy_id" json:"legacyID"`
	Type          string                `db:"type" json:"type"`
	ImageFileName null_v4.String        `db:"image_file_name" json:"imageFileName"`
	TagIds        []int32               `db:"tag_ids" json:"tagIds"`
	Images        json.RawMessage       `db:"images" json:"images"`
	Title         pqtype.NullRawMessage `db:"title" json:"title"`
	Description   pqtype.NullRawMessage `db:"description" json:"description"`
}

func (q *Queries) getShows(ctx context.Context, dollar_1 []int32) ([]getShowsRow, error) {
	rows, err := q.db.QueryContext(ctx, getShows, pq.Array(dollar_1))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []getShowsRow
	for rows.Next() {
		var i getShowsRow
		if err := rows.Scan(
			&i.ID,
			&i.LegacyID,
			&i.Type,
			&i.ImageFileName,
			pq.Array(&i.TagIds),
			&i.Images,
			&i.Title,
			&i.Description,
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

const listAllPermittedShowIDs = `-- name: listAllPermittedShowIDs :many
SELECT sh.id
FROM shows sh
         LEFT JOIN show_availability access ON access.id = sh.id
         LEFT JOIN show_roles roles ON roles.id = sh.id
WHERE access.available_from < NOW()
  AND access.available_to > NOW()
  AND roles.roles && ($1::character varying[])
`

func (q *Queries) listAllPermittedShowIDs(ctx context.Context, dollar_1 []string) ([]int32, error) {
	rows, err := q.db.QueryContext(ctx, listAllPermittedShowIDs, pq.Array(dollar_1))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []int32
	for rows.Next() {
		var id int32
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		items = append(items, id)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listShows = `-- name: listShows :many
WITH ts AS (SELECT shows_id,
                   json_object_agg(languages_code, title)       AS title,
                   json_object_agg(languages_code, description) AS description
            FROM shows_translations
            GROUP BY shows_id),
     tags AS (SELECT shows_id,
                     array_agg(tags_id) AS tags
              FROM shows_tags
              GROUP BY shows_id),
     images AS (WITH images AS (SELECT show_id, style, language, filename_disk
                                FROM images img
                                         JOIN directus_files df on img.file = df.id)
                SELECT show_id, json_agg(images) as json
                FROM images
                GROUP BY show_id)
SELECT sh.id,
       sh.legacy_id,
       sh.type,
       fs.filename_disk            as image_file_name,
       tags.tags::int[]            AS tag_ids,
       COALESCE(images.json, '[]') as images,
       ts.title,
       ts.description
FROM shows sh
         LEFT JOIN tags ON tags.shows_id = sh.id
         LEFT JOIN ts ON sh.id = ts.shows_id
         LEFT JOIN images ON sh.id = images.show_id
         LEFT JOIN directus_files fs ON fs.id = sh.image_file_id
`

type listShowsRow struct {
	ID            int32                 `db:"id" json:"id"`
	LegacyID      null_v4.Int           `db:"legacy_id" json:"legacyID"`
	Type          string                `db:"type" json:"type"`
	ImageFileName null_v4.String        `db:"image_file_name" json:"imageFileName"`
	TagIds        []int32               `db:"tag_ids" json:"tagIds"`
	Images        json.RawMessage       `db:"images" json:"images"`
	Title         pqtype.NullRawMessage `db:"title" json:"title"`
	Description   pqtype.NullRawMessage `db:"description" json:"description"`
}

func (q *Queries) listShows(ctx context.Context) ([]listShowsRow, error) {
	rows, err := q.db.QueryContext(ctx, listShows)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []listShowsRow
	for rows.Next() {
		var i listShowsRow
		if err := rows.Scan(
			&i.ID,
			&i.LegacyID,
			&i.Type,
			&i.ImageFileName,
			pq.Array(&i.TagIds),
			&i.Images,
			&i.Title,
			&i.Description,
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
