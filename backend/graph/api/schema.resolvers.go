package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"strconv"
	"time"

	merry "github.com/ansel1/merry/v2"
	"github.com/bcc-code/brunstadtv/backend/applications"
	"github.com/bcc-code/brunstadtv/backend/auth0"
	"github.com/bcc-code/brunstadtv/backend/batchloaders"
	"github.com/bcc-code/brunstadtv/backend/common"
	"github.com/bcc-code/brunstadtv/backend/export"
	"github.com/bcc-code/brunstadtv/backend/graph/api/generated"
	"github.com/bcc-code/brunstadtv/backend/graph/api/model"
	"github.com/bcc-code/brunstadtv/backend/user"
	"github.com/bcc-code/brunstadtv/backend/utils"
	"github.com/samber/lo"
	null "gopkg.in/guregu/null.v4"
)

// SetDevicePushToken is the resolver for the setDevicePushToken field.
func (r *mutationRootResolver) SetDevicePushToken(ctx context.Context, token string, languages []string) (*model.Device, error) {
	ginCtx, err := utils.GinCtx(ctx)
	if err != nil {
		return nil, err
	}
	profile := user.GetProfileFromCtx(ginCtx)
	if profile == nil {
		return nil, merry.New(
			"profile is null",
			merry.WithUserMessage("device must be connected to a profile, which is not supported by anonymous accounts"),
		)
	}

	for i := 0; i < len(languages); i++ {
		if len(languages[i]) != 2 {
			return nil, merry.New("invalid language", merry.WithUserMessage("Probably invalid language code"))
		}
		if i > 4 {
			return nil, merry.New("too many languages", merry.WithUserMessage("Language array too large. Max 5 entries"))
		}
	}

	d := common.Device{
		Token:     token,
		ProfileID: profile.ID,
		Name:      "default",
		UpdatedAt: time.Now(),
		Languages: languages,
	}
	err = r.Queries.SaveDevice(ginCtx, d)
	if err != nil {
		return nil, err
	}
	return &model.Device{
		Token:     d.Token,
		UpdatedAt: d.UpdatedAt.Format(time.RFC3339),
	}, nil
}

// SetEpisodeProgress is the resolver for the episodeProgress field.
func (r *mutationRootResolver) SetEpisodeProgress(ctx context.Context, id string, progress *int, duration *int) (*model.Episode, error) {
	ginCtx, err := utils.GinCtx(ctx)
	if err != nil {
		return nil, err
	}
	p := user.GetProfileFromCtx(ginCtx)
	if p == nil {
		return nil, ErrProfileNotSet
	}
	e, err := r.QueryRoot().Episode(ctx, id)
	if err != nil {
		return nil, err
	}
	episodeID := utils.AsInt(e.ID)
	q := r.Queries.ProfileQueries(p.ID)
	var episodeProgress *common.Progress
	pl := r.ProfileLoaders(ctx).ProgressLoader
	if progress == nil {
		err = q.ClearProgress(ctx, episodeID)
	} else {
		episodeProgress, err = pl.Get(ctx, utils.AsInt(e.ID))
		if err != nil {
			return nil, err
		}
		if episodeProgress == nil {
			dur := e.Duration
			if duration != nil {
				dur = *duration
			}
			var showID null.Int
			if e.Season != nil {
				s, err := r.QueryRoot().Season(ctx, e.Season.ID)
				if err != nil {
					return nil, err
				}
				showID.SetValid(int64(utils.AsInt(s.Show.ID)))
			}
			episodeProgress = &common.Progress{
				EpisodeID: episodeID,
				ShowID:    showID,
				Progress:  *progress,
				Duration:  dur,
				UpdatedAt: time.Now(),
			}
		}

		if float64(episodeProgress.Progress)/float64(episodeProgress.Duration) > 0.8 {
			if !episodeProgress.WatchedAt.Valid || episodeProgress.WatchedAt.Time.After(time.Now().Add(time.Hour*-12)) {
				episodeProgress.Watched++
				episodeProgress.WatchedAt = null.TimeFrom(time.Now())
			}
		}

		err = q.SaveProgress(ctx, *episodeProgress)
	}
	pl.Clear(ctx, episodeID)
	pl.Prime(ctx, episodeID, episodeProgress)
	return e, err
}

// Application is the resolver for the application field.
func (r *queryRootResolver) Application(ctx context.Context) (*model.Application, error) {
	ginCtx, err := utils.GinCtx(ctx)
	if err != nil {
		return nil, err
	}
	app, err := applications.GetFromCtx(ginCtx)
	if err != nil {
		return nil, err
	}

	var page *model.Page
	if app.DefaultPageID.Valid {
		page = &model.Page{
			ID: strconv.Itoa(int(app.DefaultPageID.Int64)),
		}
	}
	var searchPage *model.Page
	if app.SearchPageID.Valid {
		searchPage = &model.Page{
			ID: strconv.Itoa(int(app.SearchPageID.Int64)),
		}
	}

	return &model.Application{
		ID:            strconv.Itoa(app.ID),
		Code:          app.Code,
		Page:          page,
		SearchPage:    searchPage,
		ClientVersion: app.ClientVersion,
	}, nil
}

// Export is the resolver for the export field.
func (r *queryRootResolver) Export(ctx context.Context, groups []string) (*model.Export, error) {
	ginCtx, err := utils.GinCtx(ctx)
	if err != nil {
		return nil, err
	}

	profile := user.GetProfileFromCtx(ginCtx)
	if profile == nil {
		return nil, merry.New(
			"Not authorized",
			merry.WithUserMessage("you are not authorized for this query"),
		)
	}

	url, err := export.DoExport(ctx, r, r.AWSConfig.GetTempStorageBucket())
	if err != nil {
		return nil, err
	}

	return &model.Export{
		URL:       url,
		DbVersion: export.SQLiteExportDBVersion,
	}, nil
}

// Page is the resolver for the page field.
func (r *queryRootResolver) Page(ctx context.Context, id *string, code *string) (*model.Page, error) {
	if id != nil {
		return resolverForIntID(ctx, &itemLoaders[int, common.Page]{
			Item:        r.Loaders.PageLoader,
			Permissions: r.Loaders.PagePermissionLoader,
		}, *id, model.PageFrom)
	}
	if code != nil {
		intID, err := batchloaders.GetByID(ctx, r.Loaders.PageIDFromCodeLoader, *code)
		if err != nil {
			return nil, err
		}
		if intID == nil {
			return nil, merry.Sentinel("No page found with that code")
		}
		return resolverFor(ctx, &itemLoaders[int, common.Page]{
			Item:        r.Loaders.PageLoader,
			Permissions: r.Loaders.PagePermissionLoader,
		}, *intID, model.PageFrom)
	}
	return nil, merry.Sentinel("Must specify either ID or code", merry.WithHTTPCode(400))
}

// Section is the resolver for the section field.
func (r *queryRootResolver) Section(ctx context.Context, id string, timestamp *string) (model.Section, error) {
	return resolverForIntID(ctx, &itemLoaders[int, common.Section]{
		Item:        r.Loaders.SectionLoader,
		Permissions: r.Loaders.SectionPermissionLoader,
	}, id, model.SectionFrom)
}

// Show is the resolver for the show field.
func (r *queryRootResolver) Show(ctx context.Context, id string) (*model.Show, error) {
	return resolverForIntID(ctx, &itemLoaders[int, common.Show]{
		Item:        r.Loaders.ShowLoader,
		Permissions: r.Loaders.ShowPermissionLoader,
	}, id, model.ShowFrom)
}

// Season is the resolver for the season field.
func (r *queryRootResolver) Season(ctx context.Context, id string) (*model.Season, error) {
	return resolverForIntID(ctx, &itemLoaders[int, common.Season]{
		Item:        r.Loaders.SeasonLoader,
		Permissions: r.Loaders.SeasonPermissionLoader,
	}, id, model.SeasonFrom)
}

// Episode is the resolver for the episode field.
func (r *queryRootResolver) Episode(ctx context.Context, id string) (*model.Episode, error) {
	return resolverForIntID(ctx, &itemLoaders[int, common.Episode]{
		Item:        r.Loaders.EpisodeLoader,
		Permissions: r.Loaders.EpisodePermissionLoader,
	}, id, model.EpisodeFrom)
}

// Collection is the resolver for the collection field.
func (r *queryRootResolver) Collection(ctx context.Context, id string) (*model.Collection, error) {
	return resolverForIntID(ctx, &itemLoaders[int, common.Collection]{
		Item: r.Loaders.CollectionLoader,
	}, id, model.CollectionFrom)
}

// Search is the resolver for the search field.
func (r *queryRootResolver) Search(ctx context.Context, queryString string, first *int, offset *int, typeArg *string, minScore *int) (*model.SearchResult, error) {
	return searchResolver(r, ctx, queryString, first, offset, typeArg, minScore)
}

// Calendar is the resolver for the calendar field.
func (r *queryRootResolver) Calendar(ctx context.Context) (*model.Calendar, error) {
	return &model.Calendar{}, nil
}

// Event is the resolver for the event field.
func (r *queryRootResolver) Event(ctx context.Context, id string) (*model.Event, error) {
	return resolverForIntID(ctx, &itemLoaders[int, common.Event]{
		Item: r.Loaders.EventLoader,
	}, id, model.EventFrom)
}

// Faq is the resolver for the faq field.
func (r *queryRootResolver) Faq(ctx context.Context) (*model.Faq, error) {
	return &model.Faq{}, nil
}

// Me is the resolver for the me field.
func (r *queryRootResolver) Me(ctx context.Context) (*model.User, error) {
	gc, err := utils.GinCtx(ctx)
	if err != nil {
		return nil, err
	}

	usr := user.GetFromCtx(gc)

	u := &model.User{
		Anonymous: usr.IsAnonymous(),
		BccMember: usr.IsActiveBCC(),
		Roles:     usr.Roles,
	}

	if pid := gc.GetString(auth0.CtxUserID); pid != "" {
		u.ID = &pid
	}

	//if aud := gc.GetString(auth0.CtxAudience); aud != "" {
	//	u.Audience = &aud
	//}

	if usr.Email != "" {
		u.Email = &usr.Email
	}

	return u, nil
}

// Config is the resolver for the config field.
func (r *queryRootResolver) Config(ctx context.Context) (*model.Config, error) {
	return &model.Config{}, nil
}

// Profiles is the resolver for the profiles field.
func (r *queryRootResolver) Profiles(ctx context.Context) ([]*model.Profile, error) {
	ginCtx, err := utils.GinCtx(ctx)
	if err != nil {
		return nil, err
	}
	profiles := user.GetProfilesFromCtx(ginCtx)

	return lo.Map(profiles, func(i common.Profile, _ int) *model.Profile {
		return &model.Profile{
			ID:   i.ID.String(),
			Name: i.Name,
		}
	}), nil
}

// Profile is the resolver for the profile field.
func (r *queryRootResolver) Profile(ctx context.Context) (*model.Profile, error) {
	ginCtx, err := utils.GinCtx(ctx)
	if err != nil {
		return nil, err
	}
	profile := user.GetProfileFromCtx(ginCtx)

	return &model.Profile{
		ID:   profile.ID.String(),
		Name: profile.Name,
	}, nil
}

// LegacyIDLookup is the resolver for the legacyIDLookup field.
func (r *queryRootResolver) LegacyIDLookup(ctx context.Context, options *model.LegacyIDLookupOptions) (*model.LegacyIDLookup, error) {
	var id *int
	var err error
	if options.EpisodeID != nil {
		id, err = batchloaders.GetByID(ctx, r.Loaders.EpisodeIDFromLegacyIDLoader, *options.EpisodeID)
	}
	if options.ProgramID != nil {
		id, err = batchloaders.GetByID(ctx, r.Loaders.EpisodeIDFromLegacyProgramIDLoader, *options.ProgramID)
	}
	if err != nil {
		return nil, err
	}
	if id == nil {
		return nil, ErrItemNotFound
	}
	return &model.LegacyIDLookup{
		ID: strconv.Itoa(*id),
	}, nil
}

// MutationRoot returns generated.MutationRootResolver implementation.
func (r *Resolver) MutationRoot() generated.MutationRootResolver { return &mutationRootResolver{r} }

// QueryRoot returns generated.QueryRootResolver implementation.
func (r *Resolver) QueryRoot() generated.QueryRootResolver { return &queryRootResolver{r} }

type mutationRootResolver struct{ *Resolver }
type queryRootResolver struct{ *Resolver }
