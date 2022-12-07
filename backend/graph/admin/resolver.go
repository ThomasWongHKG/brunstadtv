package graph

import (
	"context"
	"database/sql"
	"github.com/bcc-code/brunstadtv/backend/batchloaders"
	"github.com/bcc-code/brunstadtv/backend/common"
	"github.com/bcc-code/brunstadtv/backend/graph/admin/model"
	"github.com/bcc-code/brunstadtv/backend/items/collection"
	"github.com/bcc-code/brunstadtv/backend/sqlc"
	"github.com/bcc-code/brunstadtv/backend/user"
	"github.com/bcc-code/brunstadtv/backend/utils"
	"github.com/samber/lo"
	"strconv"
)

// Resolver contains the common properties for all endpoints
type Resolver struct {
	DB      *sql.DB
	Queries *sqlc.Queries
	Loaders *common.BatchLoaders
}

func (r *previewResolver) getItemsForFilter(ctx context.Context, filter common.Filter) ([]*model.CollectionItem, error) {
	ids, err := collection.GetItemIDsForFilter(ctx, r.DB, nil, filter, false)
	if err != nil {
		return nil, err
	}

	ginCtx, err := utils.GinCtx(ctx)
	if err != nil {
		return nil, err
	}

	var items []*model.CollectionItem
	languages := user.GetLanguagesFromCtx(ginCtx)

	getIDs := func(col string, ids []common.Identifier) []int {
		return lo.Map(lo.Filter(ids, func(i common.Identifier, _ int) bool {
			return i.Collection == col
		}), func(i common.Identifier, _ int) int {
			return i.ID
		})
	}

	showIDs := getIDs("shows", ids)
	if len(showIDs) > 0 {
		r.Loaders.ShowLoader.LoadMany(ctx, showIDs)
	}
	seasonIDs := getIDs("seasons", ids)
	if len(seasonIDs) > 0 {
		r.Loaders.SeasonLoader.LoadMany(ctx, seasonIDs)
	}
	episodeIDs := getIDs("episodes", ids)
	if len(episodeIDs) > 0 {
		r.Loaders.EpisodeLoader.LoadMany(ctx, episodeIDs)
	}

	for _, e := range ids {
		switch e.Collection {
		case "shows":
			i, err := batchloaders.GetByID(ctx, r.Loaders.ShowLoader, e.ID)
			if err != nil {
				return nil, err
			}
			items = append(items, &model.CollectionItem{
				Collection: model.CollectionShows,
				ID:         strconv.Itoa(e.ID),
				Title:      i.Title.Get(languages),
			})
		case "seasons":
			i, err := batchloaders.GetByID(ctx, r.Loaders.SeasonLoader, e.ID)
			if err != nil {
				return nil, err
			}
			items = append(items, &model.CollectionItem{
				Collection: model.CollectionSeasons,
				ID:         strconv.Itoa(e.ID),
				Title:      i.Title.Get(languages),
			})
		case "episodes":
			i, err := batchloaders.GetByID(ctx, r.Loaders.EpisodeLoader, e.ID)
			if err != nil {
				return nil, err
			}
			items = append(items, &model.CollectionItem{
				Collection: model.CollectionEpisodes,
				ID:         strconv.Itoa(e.ID),
				Title:      i.Title.Get(languages),
			})
		}
	}

	return items, nil
}
