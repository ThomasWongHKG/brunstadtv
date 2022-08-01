// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0
// source: sections.sql

package sqlc

import (
	"context"
	"time"

	"github.com/lib/pq"
	"github.com/tabbed/pqtype"
	null_v4 "gopkg.in/guregu/null.v4"
)

const getSections = `-- name: GetSections :many
WITH t AS (SELECT ts.sections_id,
                  json_object_agg(ts.languages_code, ts.title)       AS title,
                  json_object_agg(ts.languages_code, ts.description) AS description
           FROM sections_translations ts
           GROUP BY ts.sections_id),
     u AS (SELECT ug.sections_id,
                  array_agg(ug.usergroups_code) AS roles
           FROM sections_usergroups ug
           GROUP BY ug.sections_id)
SELECT s.id,
       s.page_id,
       s.style,
       s.sort,
       s.status::text = 'published'::text AS published,
       s.date_created,
       s.date_updated,
       s.collection_id,
       t.title,
       t.description,
       u.roles::character varying[] AS roles
FROM sections s
         LEFT JOIN t ON s.id = t.sections_id
         LEFT JOIN u ON s.id = u.sections_id
WHERE s.id = ANY($1::int[])
`

type GetSectionsRow struct {
	ID           int32                 `db:"id" json:"id"`
	PageID       null_v4.Int           `db:"page_id" json:"pageID"`
	Style        null_v4.String        `db:"style" json:"style"`
	Sort         null_v4.Int           `db:"sort" json:"sort"`
	Published    bool                  `db:"published" json:"published"`
	DateCreated  time.Time             `db:"date_created" json:"dateCreated"`
	DateUpdated  time.Time             `db:"date_updated" json:"dateUpdated"`
	CollectionID null_v4.Int           `db:"collection_id" json:"collectionID"`
	Title        pqtype.NullRawMessage `db:"title" json:"title"`
	Description  pqtype.NullRawMessage `db:"description" json:"description"`
	Roles        []string              `db:"roles" json:"roles"`
}

func (q *Queries) GetSections(ctx context.Context, dollar_1 []int32) ([]GetSectionsRow, error) {
	rows, err := q.db.QueryContext(ctx, getSections, pq.Array(dollar_1))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetSectionsRow
	for rows.Next() {
		var i GetSectionsRow
		if err := rows.Scan(
			&i.ID,
			&i.PageID,
			&i.Style,
			&i.Sort,
			&i.Published,
			&i.DateCreated,
			&i.DateUpdated,
			&i.CollectionID,
			&i.Title,
			&i.Description,
			pq.Array(&i.Roles),
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

const getSectionsForPageIDs = `-- name: GetSectionsForPageIDs :many
WITH t AS (SELECT ts.sections_id,
                  json_object_agg(ts.languages_code, ts.title)       AS title,
                  json_object_agg(ts.languages_code, ts.description) AS description
           FROM sections_translations ts
           GROUP BY ts.sections_id),
     u AS (SELECT ug.sections_id,
                  array_agg(ug.usergroups_code) AS roles
           FROM sections_usergroups ug
           GROUP BY ug.sections_id)
SELECT s.id,
       s.page_id,
       s.style,
       s.sort,
       s.status::text = 'published'::text AS published,
       s.date_created,
       s.date_updated,
       s.collection_id,
       t.title,
       t.description,
       u.roles::character varying[] AS roles
FROM sections s
         LEFT JOIN t ON s.id = t.sections_id
         LEFT JOIN u ON s.id = u.sections_id
WHERE s.page_id = ANY($1::int[])
`

type GetSectionsForPageIDsRow struct {
	ID           int32                 `db:"id" json:"id"`
	PageID       null_v4.Int           `db:"page_id" json:"pageID"`
	Style        null_v4.String        `db:"style" json:"style"`
	Sort         null_v4.Int           `db:"sort" json:"sort"`
	Published    bool                  `db:"published" json:"published"`
	DateCreated  time.Time             `db:"date_created" json:"dateCreated"`
	DateUpdated  time.Time             `db:"date_updated" json:"dateUpdated"`
	CollectionID null_v4.Int           `db:"collection_id" json:"collectionID"`
	Title        pqtype.NullRawMessage `db:"title" json:"title"`
	Description  pqtype.NullRawMessage `db:"description" json:"description"`
	Roles        []string              `db:"roles" json:"roles"`
}

func (q *Queries) GetSectionsForPageIDs(ctx context.Context, dollar_1 []int32) ([]GetSectionsForPageIDsRow, error) {
	rows, err := q.db.QueryContext(ctx, getSectionsForPageIDs, pq.Array(dollar_1))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetSectionsForPageIDsRow
	for rows.Next() {
		var i GetSectionsForPageIDsRow
		if err := rows.Scan(
			&i.ID,
			&i.PageID,
			&i.Style,
			&i.Sort,
			&i.Published,
			&i.DateCreated,
			&i.DateUpdated,
			&i.CollectionID,
			&i.Title,
			&i.Description,
			pq.Array(&i.Roles),
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

const listSections = `-- name: ListSections :many
WITH t AS (SELECT ts.sections_id,
                  json_object_agg(ts.languages_code, ts.title)       AS title,
                  json_object_agg(ts.languages_code, ts.description) AS description
           FROM sections_translations ts
           GROUP BY ts.sections_id),
     u AS (SELECT ug.sections_id,
                  array_agg(ug.usergroups_code) AS roles
           FROM sections_usergroups ug
           GROUP BY ug.sections_id)
SELECT s.id,
       s.page_id,
       s.style,
       s.sort,
       s.status::text = 'published'::text AS published,
       s.date_created,
       s.date_updated,
       s.collection_id,
       t.title,
       t.description,
       u.roles::character varying[] AS roles
FROM sections s
         LEFT JOIN t ON s.id = t.sections_id
         LEFT JOIN u ON s.id = u.sections_id
`

type ListSectionsRow struct {
	ID           int32                 `db:"id" json:"id"`
	PageID       null_v4.Int           `db:"page_id" json:"pageID"`
	Style        null_v4.String        `db:"style" json:"style"`
	Sort         null_v4.Int           `db:"sort" json:"sort"`
	Published    bool                  `db:"published" json:"published"`
	DateCreated  time.Time             `db:"date_created" json:"dateCreated"`
	DateUpdated  time.Time             `db:"date_updated" json:"dateUpdated"`
	CollectionID null_v4.Int           `db:"collection_id" json:"collectionID"`
	Title        pqtype.NullRawMessage `db:"title" json:"title"`
	Description  pqtype.NullRawMessage `db:"description" json:"description"`
	Roles        []string              `db:"roles" json:"roles"`
}

func (q *Queries) ListSections(ctx context.Context) ([]ListSectionsRow, error) {
	rows, err := q.db.QueryContext(ctx, listSections)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListSectionsRow
	for rows.Next() {
		var i ListSectionsRow
		if err := rows.Scan(
			&i.ID,
			&i.PageID,
			&i.Style,
			&i.Sort,
			&i.Published,
			&i.DateCreated,
			&i.DateUpdated,
			&i.CollectionID,
			&i.Title,
			&i.Description,
			pq.Array(&i.Roles),
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
