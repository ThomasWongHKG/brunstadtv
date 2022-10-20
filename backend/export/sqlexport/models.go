// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0

package sqlexport

import (
	"database/sql"
)

type Application struct {
	ID            int64         `db:"id" json:"id"`
	Code          string        `db:"code" json:"code"`
	ClientVersion string        `db:"client_version" json:"clientVersion"`
	DefaultPageID sql.NullInt64 `db:"default_page_id" json:"defaultPageID"`
}

type Episode struct {
	ID               int64          `db:"id" json:"id"`
	LegacyID         sql.NullInt64  `db:"legacy_id" json:"legacyID"`
	LegacyProgramID  sql.NullInt64  `db:"legacy_program_id" json:"legacyProgramID"`
	AgeRating        string         `db:"age_rating" json:"ageRating"`
	Title            string         `db:"title" json:"title"`
	Description      string         `db:"description" json:"description"`
	ExtraDescription string         `db:"extra_description" json:"extraDescription"`
	Image            sql.NullString `db:"image" json:"image"`
	ImageUrl         sql.NullString `db:"image_url" json:"imageUrl"`
	ProductionDate   sql.NullString `db:"production_date" json:"productionDate"`
	SeasonID         sql.NullInt64  `db:"season_id" json:"seasonID"`
	Duration         int64          `db:"duration" json:"duration"`
	Number           int64          `db:"number" json:"number"`
}

type Page struct {
	ID          int64          `db:"id" json:"id"`
	Code        string         `db:"code" json:"code"`
	Title       string         `db:"title" json:"title"`
	Description string         `db:"description" json:"description"`
	Image       sql.NullString `db:"image" json:"image"`
	SectionIds  string         `db:"section_ids" json:"sectionIds"`
}

type Season struct {
	ID          int64          `db:"id" json:"id"`
	LegacyID    sql.NullInt64  `db:"legacy_id" json:"legacyID"`
	TagIds      string         `db:"tag_ids" json:"tagIds"`
	Number      int64          `db:"number" json:"number"`
	AgeRating   string         `db:"age_rating" json:"ageRating"`
	Title       string         `db:"title" json:"title"`
	Description string         `db:"description" json:"description"`
	ShowID      int64          `db:"show_id" json:"showID"`
	Image       sql.NullString `db:"image" json:"image"`
}

type Show struct {
	ID             int64          `db:"id" json:"id"`
	Type           string         `db:"type" json:"type"`
	LegacyID       sql.NullInt64  `db:"legacy_id" json:"legacyID"`
	Title          string         `db:"title" json:"title"`
	Description    string         `db:"description" json:"description"`
	Image          sql.NullString `db:"image" json:"image"`
	DefaultEpisode sql.NullString `db:"default_episode" json:"defaultEpisode"`
}
