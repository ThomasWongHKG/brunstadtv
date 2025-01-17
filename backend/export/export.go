package export

import (
	"context"
	"database/sql"
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/gin-gonic/gin"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/bcc-code/brunstadtv/backend/user"
	"github.com/bcc-code/brunstadtv/backend/utils"
	"github.com/google/uuid"

	"github.com/ansel1/merry/v2"
	"github.com/bcc-code/brunstadtv/backend/applications"
	"github.com/bcc-code/brunstadtv/backend/batchloaders"
	"github.com/bcc-code/brunstadtv/backend/common"
	"github.com/bcc-code/brunstadtv/backend/export/sqlexport"
	"github.com/bcc-code/brunstadtv/backend/items/collection"
	"github.com/bcc-code/brunstadtv/backend/sqlc"
	"github.com/bcc-code/mediabank-bridge/log"
	_ "github.com/glebarez/go-sqlite"
	"github.com/samber/lo"

	"github.com/pressly/goose/v3"
)

// SQLiteExportDBVersion is in semver format (https://semver.org/)
// What constitutes the various levels of change:
// * Anything with a pre-release tag is a wildcard. You can do whatever you want
// * Patch version is when the changes could not possibly break existing queries:
//   - Cosmetic changes, like comments
//   - Adding a constraint, index (since this is a change from the generator side)
//   - Adding a view
//   - Adding a table
//   - NOTE: Adding a field is explicitly excluded from patch level changes as that can break `SELECT * FROM` style queries
//   - NOTE: This is unikely to be used, and should only be if the change is really and clearly "tame".
//
// * Minor changes are changes like:
//   - Adding a field
//   - Changing a view/stored procedure while maintaining the In/Out signature
//   - NOTE: You should be quick to bump the minor version if there is any suspicion that the change could have impact on existing queries, but is stil (on "paper") backwards compatible
//   - NOTE: Changing column types is a major change, because of static langs
//
// * Major change: Everything else
const SQLiteExportDBVersion = "v0.0.1-beta"

type serviceProvider interface {
	GetQueries() *sqlc.Queries
	GetLoaders() *common.BatchLoaders
	GetFilteredLoaders(ctx context.Context) *common.FilteredLoaders
	GetS3Client() *s3.Client
}

// Embed the migrations into the binary
//
//go:embed sqlc/migrations/*.sql
var embedMigrations embed.FS

func initDB() (*sql.DB, string, error) {
	dir := os.TempDir()
	filename, _ := uuid.NewUUID()
	fn := filepath.Join(dir, fmt.Sprintf("%s.db", filename))

	log.L.Info().Str("dbPath", fn).Msg("SQLite DB created")
	db, err := sql.Open("sqlite", fn)
	if err != nil {
		return nil, "", merry.Wrap(err)
	}

	return db, fn, nil
}

func migrate(db *sql.DB) error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("sqlite"); err != nil {
		return merry.Wrap(err)
	}

	err := goose.Up(db, "sqlc/migrations")
	return merry.Wrap(err)
}

func exportShows(ctx *gin.Context, q serviceProvider, liteQueries *sqlexport.Queries) ([]int, error) {
	showIDs, err := q.GetQueries().ListAllPermittedShowIDs(ctx, user.GetRolesFromCtx(ctx))
	if err != nil {
		err = merry.Wrap(err)
		return nil, err
	}

	shows, err := batchloaders.GetMany(ctx, q.GetLoaders().ShowLoader, showIDs)
	if err != nil {
		err = merry.Wrap(err)
		return nil, err
	}

	for _, s := range shows {
		de := sql.NullString{}
		_ = de.Scan(s.DefaultEpisode)

		err := liteQueries.InsertShow(ctx, sqlexport.InsertShowParams{
			ID:             int64(s.ID),
			Type:           s.Type,
			LegacyID:       s.LegacyID.NullInt64,
			Title:          string(s.Title.AsJSON()),
			Description:    string(s.Description.AsJSON()),
			Image:          s.Image.NullString,
			DefaultEpisode: de,
		})

		if err != nil {
			err = merry.Wrap(err)
			return nil, err
		}
	}

	return showIDs, nil
}

func exportSeasons(ctx context.Context, q serviceProvider, liteQueries *sqlexport.Queries, showIDs []int) ([]int, error) {
	filteredLoaders := q.GetFilteredLoaders(ctx)

	// TODO: Refactor? common.GetManyFromLoader is refusing to fit nicely.
	thunk := filteredLoaders.SeasonsLoader.LoadMany(ctx, showIDs)
	seasonIDsResult, errs := thunk()

	if len(errs) > 0 {
		log.L.Error().Errs("errs", errs)
		err := merry.New("err getting season IDs")
		return nil, err
	}

	seasonIDs := lo.Map(lo.Flatten(seasonIDsResult), func(i *int, _ int) int { return *i })
	seasons, err := batchloaders.GetMany(ctx, q.GetLoaders().SeasonLoader, seasonIDs)
	if err != nil {
		return nil, merry.Wrap(err)
	}

	for _, s := range seasons {
		tagIds, _ := json.Marshal(s.TagIDs)
		err := liteQueries.InsertSeason(ctx, sqlexport.InsertSeasonParams{
			ID:          int64(s.ID),
			LegacyID:    s.LegacyID.NullInt64,
			TagIds:      string(tagIds),
			Number:      int64(s.Number),
			AgeRating:   s.AgeRating,
			Title:       string(s.Title.AsJSON()),
			Description: string(s.Description.AsJSON()),
			ShowID:      int64(s.ShowID),
			Image:       s.Image.NullString,
		})

		if err != nil {
			log.L.Debug().Err(err).Msg("Err while inserting season")
		}
	}
	return seasonIDs, err
}

func exportEpisodes(ctx context.Context, q serviceProvider, liteQueries *sqlexport.Queries, seasonIDs []int) ([]int, error) {
	filteredLoaders := q.GetFilteredLoaders(ctx)

	// TODO: Refactor? common.GetManyFromLoader is refusing to fit nicely.
	thunk := filteredLoaders.EpisodesLoader.LoadMany(ctx, seasonIDs)
	episodesIDsResult, errs := thunk()

	if len(errs) > 0 {
		log.L.Error().Errs("errs", errs)
		err := merry.New("err getting episode IDs")
		return nil, err
	}

	episodeIDs := lo.Map(lo.Flatten(episodesIDsResult), func(i *int, _ int) int { return *i })
	episodes, err := batchloaders.GetMany(ctx, q.GetLoaders().EpisodeLoader, episodeIDs)
	if err != nil {
		return nil, merry.Wrap(err)
	}

	for _, e := range episodes {
		err := liteQueries.InsertEpisode(ctx, sqlexport.InsertEpisodeParams{
			ID:               int64(e.ID),
			LegacyID:         e.LegacyID.NullInt64,
			LegacyProgramID:  e.LegacyProgramID.NullInt64,
			AgeRating:        e.AgeRating,
			Title:            string(e.Title.AsJSON()),
			Description:      string(e.Description.AsJSON()),
			ExtraDescription: string(e.ExtraDescription.AsJSON()),
			Image:            e.Image.NullString,
			ProductionDate:   sql.NullString{String: e.PublishDate.Format(time.RFC1123), Valid: true},
			SeasonID:         e.SeasonID.NullInt64,
			Duration:         int64(e.Duration),
			Number:           e.Number.Int64,
		})

		if err != nil {
			log.L.Debug().Err(err).Msg("Err while inserting season")
		}
	}
	return seasonIDs, err
}

func exportCurrentApplication(ctx *gin.Context, liteQueries *sqlexport.Queries) error {
	app, err := applications.GetFromCtx(ctx)
	if err != nil {
		return err
	}

	return liteQueries.InsertApplication(ctx, sqlexport.InsertApplicationParams{
		ID:            int64(app.ID),
		Code:          app.Code,
		ClientVersion: app.ClientVersion,
		DefaultPageID: app.DefaultPageID.NullInt64,
	})
}

func exportSections(ctx context.Context, q serviceProvider, liteQueries *sqlexport.Queries) ([]int, []int, error) {
	filteredLoaders := q.GetFilteredLoaders(ctx)
	pages, err := q.GetQueries().ListPages(ctx)
	if err != err {
		return nil, nil, err
	}

	allPageIDs := lo.Map(pages, func(p common.Page, _ int) int { return p.ID })
	thunk := filteredLoaders.SectionsLoader.LoadMany(ctx, allPageIDs)
	sectionIDsResult, errs := thunk()
	if len(errs) > 0 {
		log.L.Error().Errs("errs", errs)
		err := merry.New("err getting section IDs")
		return nil, nil, err
	}

	sectionIDs := lo.Map(lo.Flatten(sectionIDsResult), func(i *int, _ int) int { return *i })
	sections, err := batchloaders.GetMany(ctx, q.GetLoaders().SectionLoader, sectionIDs)
	if err != nil {
		return nil, nil, merry.Wrap(err)
	}

	allowedPageIDs := map[int]interface{}{}
	neededCollectionIDs := map[int]interface{}{}
	for _, s := range sections {

		allowedPageIDs[s.PageID] = nil
		neededCollectionIDs[int(s.CollectionID.ValueOrZero())] = nil

		err := liteQueries.InsertSection(ctx, sqlexport.InsertSectionParams{
			ID:           int64(s.ID),
			Sort:         int64(s.Sort),
			PageID:       int64(s.PageID),
			Type:         s.Type,
			ShowTitle:    s.ShowTitle,
			Title:        string(s.Title.AsJSON()),
			Description:  string(s.Description.AsJSON()),
			Style:        s.Style,
			Size:         s.Size,
			CollectionID: s.CollectionID.NullInt64,
		})

		if err != nil {
			err = merry.Wrap(err)
			return nil, nil, err
		}
	}

	return lo.Keys(allowedPageIDs), lo.Keys(neededCollectionIDs), nil
}

func exportPages(ctx context.Context, q serviceProvider, liteQueries *sqlexport.Queries, pageIDs []int) error {
	pages, err := batchloaders.GetMany(ctx, q.GetLoaders().PageLoader, pageIDs)
	if err != nil {
		return merry.Wrap(err)
	}

	for _, p := range pages {
		img := sql.NullString{}
		_ = img.Scan(p.Images.GetDefault([]string{"no"}, "default"))

		err := liteQueries.InsertPage(ctx, sqlexport.InsertPageParams{
			ID:          int64(p.ID),
			Code:        p.Code,
			Title:       string(p.Title.AsJSON()),
			Description: string(p.Description.AsJSON()),
			Image:       img,
		})

		if err != nil {
			err = merry.Wrap(err)
			return err
		}
	}

	return nil
}

func exportCollections(ctx context.Context, q serviceProvider, liteQueries *sqlexport.Queries, collectionIDs []int) error {
	filteredLoaders := q.GetFilteredLoaders(ctx)
	collections, err := batchloaders.GetMany(ctx, q.GetLoaders().CollectionLoader, collectionIDs)
	if err != nil {
		return merry.Wrap(err)
	}

	for _, c := range collections {
		entries, err := collection.GetCollectionEntries(ctx, q.GetLoaders(), filteredLoaders, c.ID)
		if err != nil {
			return merry.Wrap(err)
		}

		entryIDs := lo.Map(entries, func(e collection.Entry, _ int) int { return e.ID })
		entryIDsJSON, _ := json.Marshal(entryIDs)

		err = liteQueries.InsertCollection(ctx, sqlexport.InsertCollectionParams{
			ID: int64(c.ID),
			// TODO: maybe remove name?
			Name:            c.Title.Get([]string{"no"}),
			Type:            c.Type,
			CollectionItems: string(entryIDsJSON),
		})

		if err != nil {
			err = merry.Wrap(err)
			return err
		}
	}
	return nil
}

// DoExport exports some key data into a SQLite DB and uploads that to the provided S3 bucket
// It then returns a pre-signed link to the file that remains valid for 1 hour
//
// The rest if the functions in this file are not exported because they are currently dependent on each other and
// are basically split only on order to understand the flow better.
func DoExport(ctx context.Context, q serviceProvider, bucketName string) (string, error) {
	//TODO: Caching?
	gctx, err := utils.GinCtx(ctx)
	if err != nil {
		return "", merry.Wrap(err)
	}

	db, dbPath, err := initDB()
	if err != nil {
		return "", merry.Wrap(err, merry.WithUserMessage("Unable to generate export file"))
	}

	err = migrate(db)
	if err != nil {
		return "", merry.Wrap(err, merry.WithUserMessage("Unable to generate export file"))
	}
	liteQueries := sqlexport.New(db)

	showIDs, err := exportShows(gctx, q, liteQueries)
	if err != nil {
		return "", merry.Wrap(err, merry.WithUserMessage("Unable to generate export file"))
	}

	seasonIDs, err := exportSeasons(ctx, q, liteQueries, showIDs)
	if err != nil {
		return "", merry.Wrap(err, merry.WithUserMessage("Unable to generate export file"))
	}

	_, err = exportEpisodes(ctx, q, liteQueries, seasonIDs)
	if err != nil {
		return "", merry.Wrap(err, merry.WithUserMessage("Unable to generate export file"))
	}

	// Just the current app for now. We can look into expanding later
	err = exportCurrentApplication(gctx, liteQueries)
	if err != nil {
		return "", merry.Wrap(err, merry.WithUserMessage("Unable to generate export file"))
	}

	pagesToExport, collectionsToExport, err := exportSections(ctx, q, liteQueries)
	if err != nil {
		return "", merry.Wrap(err, merry.WithUserMessage("Unable to generate export file"))
	}

	err = exportPages(ctx, q, liteQueries, pagesToExport)
	if err != nil {
		return "", merry.Wrap(err, merry.WithUserMessage("Unable to generate export file"))
	}

	err = exportCollections(ctx, q, liteQueries, collectionsToExport)
	if err != nil {
		return "", merry.Wrap(err, merry.WithUserMessage("Unable to generate export file"))
	}

	err = db.Close()
	if err != nil {
		return "", merry.Wrap(err, merry.WithUserMessage("Unable to generate export file"))
	}

	f, err := os.Open(dbPath)
	if err != nil {
		return "", merry.Wrap(err, merry.WithUserMessage("Unable to generate export file"))
	}

	s3DestinationPath := aws.String(path.Join("/sqliteexport", filepath.Base(dbPath)))

	_, err = q.GetS3Client().PutObject(ctx, &s3.PutObjectInput{
		Body:         f,
		Bucket:       aws.String(bucketName),
		CacheControl: aws.String("Cache-Control: private, max-age=604800, immutable"),
		Key:          s3DestinationPath,
	})
	if err != nil {
		return "", merry.Wrap(err, merry.WithUserMessage("Unable to generate export file"))
	}

	presignClient := s3.NewPresignClient(q.GetS3Client(), s3.WithPresignExpires(1*time.Hour))
	if err != nil {
		return "", merry.Wrap(err, merry.WithUserMessage("Unable to generate export file"))
	}

	res, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName), // TODO - from external
		Key:    s3DestinationPath,
	})
	if err != nil {
		return "", merry.Wrap(err, merry.WithUserMessage("Unable to generate export file"))
	}
	return res.URL, nil
}
