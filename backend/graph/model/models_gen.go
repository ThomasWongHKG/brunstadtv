// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package gqlmodel

import (
	"fmt"
	"io"
	"strconv"
)

type CalendarEntry interface {
	IsCalendarEntry()
}

type Item interface {
	IsItem()
}

type Pagination interface {
	IsPagination()
}

type SearchResultItem interface {
	IsSearchResultItem()
}

type Section interface {
	IsSection()
}

type AppConfig struct {
	MinVersion string `json:"minVersion"`
}

type Calendar struct {
	Period *CalendarPeriod `json:"period"`
	Day    *CalendarDay    `json:"day"`
}

type CalendarDay struct {
	Events  []*Event        `json:"events"`
	Entries []CalendarEntry `json:"entries"`
}

type CalendarPeriod struct {
	ActiveDays []string `json:"activeDays"`
	Events     []*Event `json:"events"`
}

type Chapter struct {
	ID    string `json:"id"`
	Start int    `json:"start"`
	Title string `json:"title"`
}

type Collection struct {
	ID    string                    `json:"id"`
	Items *CollectionItemPagination `json:"items"`
}

type CollectionItemPagination struct {
	Total  int    `json:"total"`
	First  int    `json:"first"`
	Offset int    `json:"offset"`
	Items  []Item `json:"items"`
}

func (CollectionItemPagination) IsPagination() {}

type Config struct {
	Global *GlobalConfig `json:"global"`
	App    *AppConfig    `json:"app"`
}

type Episode struct {
	ID                string     `json:"id"`
	LegacyID          *string    `json:"legacyID"`
	LegacyProgramID   *string    `json:"legacyProgramID"`
	Title             string     `json:"title"`
	Description       string     `json:"description"`
	ExtraDescription  string     `json:"extraDescription"`
	ImageURL          *string    `json:"imageUrl"`
	Streams           []*Stream  `json:"streams"`
	Files             []*File    `json:"files"`
	Chapters          []*Chapter `json:"chapters"`
	Season            *Season    `json:"season"`
	Duration          int        `json:"duration"`
	AudioLanguages    []Language `json:"audioLanguages"`
	SubtitleLanguages []Language `json:"subtitleLanguages"`
	Number            *int       `json:"number"`
}

type EpisodeCalendarEntry struct {
	ID          string   `json:"id"`
	Event       *Event   `json:"event"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Start       string   `json:"start"`
	End         string   `json:"end"`
	Episode     *Episode `json:"episode"`
}

func (EpisodeCalendarEntry) IsCalendarEntry() {}

type EpisodeItem struct {
	ID       string   `json:"id"`
	Sort     int      `json:"sort"`
	Title    string   `json:"title"`
	ImageURL *string  `json:"imageUrl"`
	Episode  *Episode `json:"episode"`
}

func (EpisodeItem) IsItem() {}

type EpisodePagination struct {
	Total  int        `json:"total"`
	First  int        `json:"first"`
	Offset int        `json:"offset"`
	Items  []*Episode `json:"items"`
}

func (EpisodePagination) IsPagination() {}

type EpisodeSearchItem struct {
	ID          string  `json:"id"`
	LegacyID    *string `json:"legacyID"`
	Collection  string  `json:"collection"`
	Title       string  `json:"title"`
	Header      *string `json:"header"`
	Description *string `json:"description"`
	Highlight   *string `json:"highlight"`
	Image       *string `json:"image"`
	URL         string  `json:"url"`
	ShowID      *string `json:"showId"`
	ShowTitle   *string `json:"showTitle"`
	Show        *Show   `json:"show"`
	SeasonID    *string `json:"seasonId"`
	SeasonTitle *string `json:"seasonTitle"`
	Season      *Season `json:"season"`
}

func (EpisodeSearchItem) IsSearchResultItem() {}

type Event struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Start string `json:"start"`
	End   string `json:"end"`
	Image string `json:"image"`
}

type Faq struct {
	Categories *FAQCategoryPagination `json:"categories"`
	Category   *FAQCategory           `json:"category"`
	Question   *Question              `json:"question"`
}

type FAQCategory struct {
	ID        string              `json:"id"`
	Title     string              `json:"title"`
	Questions *QuestionPagination `json:"questions"`
}

type FAQCategoryPagination struct {
	Total  int            `json:"total"`
	First  int            `json:"first"`
	Offset int            `json:"offset"`
	Items  []*FAQCategory `json:"items"`
}

func (FAQCategoryPagination) IsPagination() {}

type File struct {
	ID               string    `json:"id"`
	URL              string    `json:"url"`
	AudioLanguage    Language  `json:"audioLanguage"`
	SubtitleLanguage *Language `json:"subtitleLanguage"`
	Size             *int      `json:"size"`
	FileName         string    `json:"fileName"`
	MimeType         string    `json:"mimeType"`
}

type GlobalConfig struct {
	LiveOnline  bool `json:"liveOnline"`
	NpawEnabled bool `json:"npawEnabled"`
}

type ItemSection struct {
	ID    string                    `json:"id"`
	Page  *Page                     `json:"page"`
	Title string                    `json:"title"`
	Type  ItemSectionType           `json:"type"`
	Style string                    `json:"style"`
	Items *CollectionItemPagination `json:"items"`
}

func (ItemSection) IsSection() {}

type MaintenanceMessage struct {
	Message string  `json:"message"`
	Details *string `json:"details"`
}

type Messages struct {
	Maintenance []*MaintenanceMessage `json:"maintenance"`
}

type Page struct {
	ID          string             `json:"id"`
	Code        string             `json:"code"`
	Title       string             `json:"title"`
	Description *string            `json:"description"`
	Sections    *SectionPagination `json:"sections"`
}

type PageItem struct {
	ID       string  `json:"id"`
	Sort     int     `json:"sort"`
	Title    string  `json:"title"`
	ImageURL *string `json:"imageUrl"`
	Page     *Page   `json:"page"`
}

func (PageItem) IsItem() {}

type PagePagination struct {
	Total  int     `json:"total"`
	First  int     `json:"first"`
	Offset int     `json:"offset"`
	Items  []*Page `json:"items"`
}

func (PagePagination) IsPagination() {}

type Question struct {
	ID       string       `json:"id"`
	Category *FAQCategory `json:"category"`
	Question string       `json:"question"`
	Answer   string       `json:"answer"`
}

type QuestionPagination struct {
	Total  int         `json:"total"`
	First  int         `json:"first"`
	Offset int         `json:"offset"`
	Items  []*Question `json:"items"`
}

func (QuestionPagination) IsPagination() {}

type SearchResult struct {
	Hits   int                `json:"hits"`
	Page   int                `json:"page"`
	Result []SearchResultItem `json:"result"`
}

type Season struct {
	ID          string             `json:"id"`
	LegacyID    *string            `json:"legacyID"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	ImageURL    *string            `json:"imageUrl"`
	Number      int                `json:"number"`
	Show        *Show              `json:"show"`
	Episodes    *EpisodePagination `json:"episodes"`
}

type SeasonCalendarEntry struct {
	ID          string  `json:"id"`
	Event       *Event  `json:"event"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Start       string  `json:"start"`
	End         string  `json:"end"`
	Season      *Season `json:"season"`
}

func (SeasonCalendarEntry) IsCalendarEntry() {}

type SeasonItem struct {
	ID       string  `json:"id"`
	Sort     int     `json:"sort"`
	Title    string  `json:"title"`
	ImageURL *string `json:"imageUrl"`
	Season   *Season `json:"season"`
}

func (SeasonItem) IsItem() {}

type SeasonPagination struct {
	Total  int       `json:"total"`
	First  int       `json:"first"`
	Offset int       `json:"offset"`
	Items  []*Season `json:"items"`
}

func (SeasonPagination) IsPagination() {}

type SeasonSearchItem struct {
	ID          string  `json:"id"`
	LegacyID    *string `json:"legacyID"`
	Collection  string  `json:"collection"`
	Title       string  `json:"title"`
	Header      *string `json:"header"`
	Description *string `json:"description"`
	Highlight   *string `json:"highlight"`
	Image       *string `json:"image"`
	URL         string  `json:"url"`
	ShowID      string  `json:"showId"`
	ShowTitle   string  `json:"showTitle"`
	Show        *Show   `json:"show"`
}

func (SeasonSearchItem) IsSearchResultItem() {}

type SectionPagination struct {
	Total  int       `json:"total"`
	First  int       `json:"first"`
	Offset int       `json:"offset"`
	Items  []Section `json:"items"`
}

func (SectionPagination) IsPagination() {}

type Settings struct {
	AudioLanguages    []Language `json:"audioLanguages"`
	SubtitleLanguages []Language `json:"subtitleLanguages"`
}

type Show struct {
	ID           string            `json:"id"`
	LegacyID     *string           `json:"legacyID"`
	Title        string            `json:"title"`
	Description  string            `json:"description"`
	ImageURL     *string           `json:"imageUrl"`
	EpisodeCount int               `json:"episodeCount"`
	SeasonCount  int               `json:"seasonCount"`
	Seasons      *SeasonPagination `json:"seasons"`
}

type ShowCalendarEntry struct {
	ID          string `json:"id"`
	Event       *Event `json:"event"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Start       string `json:"start"`
	End         string `json:"end"`
	Show        *Show  `json:"show"`
}

func (ShowCalendarEntry) IsCalendarEntry() {}

type ShowItem struct {
	ID       string  `json:"id"`
	Sort     int     `json:"sort"`
	Title    string  `json:"title"`
	ImageURL *string `json:"imageUrl"`
	Show     *Show   `json:"show"`
}

func (ShowItem) IsItem() {}

type ShowPagination struct {
	Total  int     `json:"total"`
	First  int     `json:"first"`
	Offset int     `json:"offset"`
	Items  []*Show `json:"items"`
}

func (ShowPagination) IsPagination() {}

type ShowSearchItem struct {
	ID          string  `json:"id"`
	LegacyID    *string `json:"legacyID"`
	Collection  string  `json:"collection"`
	Title       string  `json:"title"`
	Header      *string `json:"header"`
	Description *string `json:"description"`
	Highlight   *string `json:"highlight"`
	Image       *string `json:"image"`
	URL         string  `json:"url"`
}

func (ShowSearchItem) IsSearchResultItem() {}

type SimpleCalendarEntry struct {
	ID          string `json:"id"`
	Event       *Event `json:"event"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Start       string `json:"start"`
	End         string `json:"end"`
}

func (SimpleCalendarEntry) IsCalendarEntry() {}

type Stream struct {
	ID                string     `json:"id"`
	URL               string     `json:"url"`
	AudioLanguages    []Language `json:"audioLanguages"`
	SubtitleLanguages []Language `json:"subtitleLanguages"`
	Type              StreamType `json:"type"`
}

type URLItem struct {
	ID       string  `json:"id"`
	Sort     int     `json:"sort"`
	Title    string  `json:"title"`
	ImageURL *string `json:"imageUrl"`
	URL      string  `json:"url"`
}

func (URLItem) IsItem() {}

type User struct {
	ID        *string   `json:"id"`
	Anonymous bool      `json:"anonymous"`
	BccMember bool      `json:"bccMember"`
	Audience  *string   `json:"audience"`
	Email     *string   `json:"email"`
	Settings  *Settings `json:"settings"`
	Roles     []string  `json:"roles"`
}

type ItemSectionType string

const (
	ItemSectionTypeCards  ItemSectionType = "cards"
	ItemSectionTypeSlider ItemSectionType = "slider"
)

var AllItemSectionType = []ItemSectionType{
	ItemSectionTypeCards,
	ItemSectionTypeSlider,
}

func (e ItemSectionType) IsValid() bool {
	switch e {
	case ItemSectionTypeCards, ItemSectionTypeSlider:
		return true
	}
	return false
}

func (e ItemSectionType) String() string {
	return string(e)
}

func (e *ItemSectionType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ItemSectionType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ItemSectionType", str)
	}
	return nil
}

func (e ItemSectionType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type Language string

const (
	LanguageEn Language = "en"
	LanguageNo Language = "no"
	LanguageDe Language = "de"
)

var AllLanguage = []Language{
	LanguageEn,
	LanguageNo,
	LanguageDe,
}

func (e Language) IsValid() bool {
	switch e {
	case LanguageEn, LanguageNo, LanguageDe:
		return true
	}
	return false
}

func (e Language) String() string {
	return string(e)
}

func (e *Language) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Language(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Language", str)
	}
	return nil
}

func (e Language) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type StreamType string

const (
	StreamTypeHlsTs   StreamType = "hls_ts"
	StreamTypeHlsCmaf StreamType = "hls_cmaf"
	StreamTypeDash    StreamType = "dash"
)

var AllStreamType = []StreamType{
	StreamTypeHlsTs,
	StreamTypeHlsCmaf,
	StreamTypeDash,
}

func (e StreamType) IsValid() bool {
	switch e {
	case StreamTypeHlsTs, StreamTypeHlsCmaf, StreamTypeDash:
		return true
	}
	return false
}

func (e StreamType) String() string {
	return string(e)
}

func (e *StreamType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = StreamType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid StreamType", str)
	}
	return nil
}

func (e StreamType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
