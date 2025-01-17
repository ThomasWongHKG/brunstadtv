// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: calendar.sql

package sqlc

import (
	"context"
	"time"

	"github.com/lib/pq"
	"github.com/tabbed/pqtype"
	null_v4 "gopkg.in/guregu/null.v4"
)

const getCalendarEntries = `-- name: getCalendarEntries :many
WITH t AS (SELECT ts.calendarentries_id,
                  json_object_agg(ts.languages_code, ts.title)       AS title,
                  json_object_agg(ts.languages_code, ts.description) AS description
           FROM calendarentries_translations ts
           GROUP BY ts.calendarentries_id)
SELECT e.id,
       e.event_id,
       e.link_type,
       e.start,
       e.end,
       e.episode_id,
       e.season_id,
       e.show_id,
       t.title,
       t.description
FROM calendarentries e
         LEFT JOIN t ON e.id = t.calendarentries_id
WHERE e.status = 'published'
AND e.id = ANY($1::int[])
`

type getCalendarEntriesRow struct {
	ID          int32                 `db:"id" json:"id"`
	EventID     null_v4.Int           `db:"event_id" json:"eventID"`
	LinkType    null_v4.String        `db:"link_type" json:"linkType"`
	Start       time.Time             `db:"start" json:"start"`
	End         time.Time             `db:"end" json:"end"`
	EpisodeID   null_v4.Int           `db:"episode_id" json:"episodeID"`
	SeasonID    null_v4.Int           `db:"season_id" json:"seasonID"`
	ShowID      null_v4.Int           `db:"show_id" json:"showID"`
	Title       pqtype.NullRawMessage `db:"title" json:"title"`
	Description pqtype.NullRawMessage `db:"description" json:"description"`
}

func (q *Queries) getCalendarEntries(ctx context.Context, dollar_1 []int32) ([]getCalendarEntriesRow, error) {
	rows, err := q.db.QueryContext(ctx, getCalendarEntries, pq.Array(dollar_1))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []getCalendarEntriesRow
	for rows.Next() {
		var i getCalendarEntriesRow
		if err := rows.Scan(
			&i.ID,
			&i.EventID,
			&i.LinkType,
			&i.Start,
			&i.End,
			&i.EpisodeID,
			&i.SeasonID,
			&i.ShowID,
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

const getCalendarEntriesForEvents = `-- name: getCalendarEntriesForEvents :many
WITH t AS (SELECT ts.calendarentries_id,
                  json_object_agg(ts.languages_code, ts.title)       AS title,
                  json_object_agg(ts.languages_code, ts.description) AS description
           FROM calendarentries_translations ts
           GROUP BY ts.calendarentries_id)
SELECT e.id,
       e.event_id,
       e.link_type,
       e.start,
       e.end,
       e.episode_id,
       e.season_id,
       e.show_id,
       t.title,
       t.description
FROM calendarentries e
         LEFT JOIN t ON e.id = t.calendarentries_id
WHERE e.status = 'published'
  AND e.event_id = ANY($1::int[])
`

type getCalendarEntriesForEventsRow struct {
	ID          int32                 `db:"id" json:"id"`
	EventID     null_v4.Int           `db:"event_id" json:"eventID"`
	LinkType    null_v4.String        `db:"link_type" json:"linkType"`
	Start       time.Time             `db:"start" json:"start"`
	End         time.Time             `db:"end" json:"end"`
	EpisodeID   null_v4.Int           `db:"episode_id" json:"episodeID"`
	SeasonID    null_v4.Int           `db:"season_id" json:"seasonID"`
	ShowID      null_v4.Int           `db:"show_id" json:"showID"`
	Title       pqtype.NullRawMessage `db:"title" json:"title"`
	Description pqtype.NullRawMessage `db:"description" json:"description"`
}

func (q *Queries) getCalendarEntriesForEvents(ctx context.Context, dollar_1 []int32) ([]getCalendarEntriesForEventsRow, error) {
	rows, err := q.db.QueryContext(ctx, getCalendarEntriesForEvents, pq.Array(dollar_1))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []getCalendarEntriesForEventsRow
	for rows.Next() {
		var i getCalendarEntriesForEventsRow
		if err := rows.Scan(
			&i.ID,
			&i.EventID,
			&i.LinkType,
			&i.Start,
			&i.End,
			&i.EpisodeID,
			&i.SeasonID,
			&i.ShowID,
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

const getCalendarEntryIDsForPeriod = `-- name: getCalendarEntryIDsForPeriod :many
SELECT e.id
FROM calendarentries e
WHERE e.status = 'published'
  AND ((e.start >= $1::timestamptz AND e.start <= $2::timestamptz) OR
       (e.end >= $1::timestamptz AND e.end <= $2::timestamptz))
ORDER BY e.start
`

type getCalendarEntryIDsForPeriodParams struct {
	Column1 time.Time `db:"column_1" json:"column1"`
	Column2 time.Time `db:"column_2" json:"column2"`
}

func (q *Queries) getCalendarEntryIDsForPeriod(ctx context.Context, arg getCalendarEntryIDsForPeriodParams) ([]int32, error) {
	rows, err := q.db.QueryContext(ctx, getCalendarEntryIDsForPeriod, arg.Column1, arg.Column2)
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

const getEventIDsForPeriod = `-- name: getEventIDsForPeriod :many
SELECT e.id
FROM events e
WHERE e.status = 'published'
  AND ((e.start >= $1::timestamptz AND e.start <= $2::timestamptz) OR
       (e.end >= $1::timestamptz AND e.end <= $2::timestamptz))
ORDER BY e.start
`

type getEventIDsForPeriodParams struct {
	Column1 time.Time `db:"column_1" json:"column1"`
	Column2 time.Time `db:"column_2" json:"column2"`
}

func (q *Queries) getEventIDsForPeriod(ctx context.Context, arg getEventIDsForPeriodParams) ([]int32, error) {
	rows, err := q.db.QueryContext(ctx, getEventIDsForPeriod, arg.Column1, arg.Column2)
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

const getEvents = `-- name: getEvents :many
WITH t AS (SELECT ts.events_id,
                  json_object_agg(ts.languages_code, ts.title)       AS title
           FROM events_translations ts
           GROUP BY ts.events_id)
SELECT e.id,
       e.start,
       e.end,
       t.title
FROM events e
         LEFT JOIN t ON e.id = t.events_id
WHERE e.status = 'published'
  AND e.id = ANY($1::int[])
`

type getEventsRow struct {
	ID    int32                 `db:"id" json:"id"`
	Start time.Time             `db:"start" json:"start"`
	End   time.Time             `db:"end" json:"end"`
	Title pqtype.NullRawMessage `db:"title" json:"title"`
}

func (q *Queries) getEvents(ctx context.Context, dollar_1 []int32) ([]getEventsRow, error) {
	rows, err := q.db.QueryContext(ctx, getEvents, pq.Array(dollar_1))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []getEventsRow
	for rows.Next() {
		var i getEventsRow
		if err := rows.Scan(
			&i.ID,
			&i.Start,
			&i.End,
			&i.Title,
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

const listCalendarEntries = `-- name: listCalendarEntries :many
WITH t AS (SELECT ts.calendarentries_id,
                  json_object_agg(ts.languages_code, ts.title)       AS title,
                  json_object_agg(ts.languages_code, ts.description) AS description
           FROM calendarentries_translations ts
           GROUP BY ts.calendarentries_id)
SELECT e.id,
       e.event_id,
       e.link_type,
       e.start,
       e.end,
       e.episode_id,
       e.season_id,
       e.show_id,
       t.title,
       t.description
FROM calendarentries e
         LEFT JOIN t ON e.id = t.calendarentries_id
WHERE e.status = 'published'
`

type listCalendarEntriesRow struct {
	ID          int32                 `db:"id" json:"id"`
	EventID     null_v4.Int           `db:"event_id" json:"eventID"`
	LinkType    null_v4.String        `db:"link_type" json:"linkType"`
	Start       time.Time             `db:"start" json:"start"`
	End         time.Time             `db:"end" json:"end"`
	EpisodeID   null_v4.Int           `db:"episode_id" json:"episodeID"`
	SeasonID    null_v4.Int           `db:"season_id" json:"seasonID"`
	ShowID      null_v4.Int           `db:"show_id" json:"showID"`
	Title       pqtype.NullRawMessage `db:"title" json:"title"`
	Description pqtype.NullRawMessage `db:"description" json:"description"`
}

func (q *Queries) listCalendarEntries(ctx context.Context) ([]listCalendarEntriesRow, error) {
	rows, err := q.db.QueryContext(ctx, listCalendarEntries)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []listCalendarEntriesRow
	for rows.Next() {
		var i listCalendarEntriesRow
		if err := rows.Scan(
			&i.ID,
			&i.EventID,
			&i.LinkType,
			&i.Start,
			&i.End,
			&i.EpisodeID,
			&i.SeasonID,
			&i.ShowID,
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

const listEvents = `-- name: listEvents :many
WITH t AS (SELECT ts.events_id,
                  json_object_agg(ts.languages_code, ts.title)       AS title
           FROM events_translations ts
           GROUP BY ts.events_id)
SELECT e.id,
       e.start,
       e.end,
       t.title
FROM events e
         LEFT JOIN t ON e.id = t.events_id
WHERE e.status = 'published'
`

type listEventsRow struct {
	ID    int32                 `db:"id" json:"id"`
	Start time.Time             `db:"start" json:"start"`
	End   time.Time             `db:"end" json:"end"`
	Title pqtype.NullRawMessage `db:"title" json:"title"`
}

func (q *Queries) listEvents(ctx context.Context) ([]listEventsRow, error) {
	rows, err := q.db.QueryContext(ctx, listEvents)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []listEventsRow
	for rows.Next() {
		var i listEventsRow
		if err := rows.Scan(
			&i.ID,
			&i.Start,
			&i.End,
			&i.Title,
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
