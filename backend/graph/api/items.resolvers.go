package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"strconv"

	"github.com/bcc-code/brunstadtv/backend/common"
	"github.com/bcc-code/brunstadtv/backend/graph/api/generated"
	"github.com/bcc-code/brunstadtv/backend/graph/api/model"
	"github.com/bcc-code/brunstadtv/backend/items/show"
	"github.com/bcc-code/brunstadtv/backend/utils"
)

// Image is the resolver for the image field.
func (r *episodeResolver) Image(ctx context.Context, obj *model.Episode, style *model.ImageStyle) (*string, error) {
	e, err := common.GetFromLoaderByID(ctx, r.Loaders.EpisodeLoader, utils.AsInt(obj.ID))
	if err != nil {
		return nil, err
	}
	return imageOrFallback(ctx, e.Images, e.Image, style), nil
}

// Streams is the resolver for the streams field.
func (r *episodeResolver) Streams(ctx context.Context, obj *model.Episode) ([]*model.Stream, error) {
	intID, _ := strconv.ParseInt(obj.ID, 10, 32)
	streams, err := common.GetFromLoaderForKey(ctx, r.Resolver.Loaders.StreamsLoader, int(intID))
	if err != nil {
		return nil, err
	}

	var out []*model.Stream
	for _, s := range streams {
		stream, err := model.StreamFrom(ctx, r.URLSigner, r.Resolver.APIConfig, s)
		if err != nil {
			return nil, err
		}

		out = append(out, stream)
	}

	return out, nil
}

// Files is the resolver for the files field.
func (r *episodeResolver) Files(ctx context.Context, obj *model.Episode) ([]*model.File, error) {
	intID, err := strconv.ParseInt(obj.ID, 10, 32)
	if err != nil {
		return nil, err
	}

	files, err := common.GetFromLoaderForKey(ctx, r.Resolver.Loaders.FilesLoader, int(intID))
	if err != nil {
		return nil, err
	}

	var out []*model.File
	for _, f := range files {
		out = append(out, model.FileFrom(ctx, r.URLSigner, r.Resolver.APIConfig.GetFilesCDNDomain(), f))
	}
	return out, nil
}

// Season is the resolver for the season field.
func (r *episodeResolver) Season(ctx context.Context, obj *model.Episode) (*model.Season, error) {
	if obj.Season != nil {
		return r.QueryRoot().Season(ctx, obj.Season.ID)
	}
	return nil, nil
}

// Progress is the resolver for the progress field.
func (r *episodeResolver) Progress(ctx context.Context, obj *model.Episode) (*int, error) {
	profileLoaders := r.ProfileLoaders(ctx)
	if profileLoaders == nil {
		return nil, nil
	}
	progress, err := common.GetFromLoaderByID(ctx, profileLoaders.ProgressLoader, utils.AsInt(obj.ID))
	if err != nil || progress == nil {
		return nil, err
	}
	return &progress.Progress, nil
}

// Image is the resolver for the image field.
func (r *seasonResolver) Image(ctx context.Context, obj *model.Season, style *model.ImageStyle) (*string, error) {
	e, err := common.GetFromLoaderByID(ctx, r.Loaders.SeasonLoader, utils.AsInt(obj.ID))
	if err != nil {
		return nil, err
	}
	return imageOrFallback(ctx, e.Images, e.Image, style), nil
}

// Show is the resolver for the show field.
func (r *seasonResolver) Show(ctx context.Context, obj *model.Season) (*model.Show, error) {
	return r.QueryRoot().Show(ctx, obj.Show.ID)
}

// Episodes is the resolver for the episodes field.
func (r *seasonResolver) Episodes(ctx context.Context, obj *model.Season, first *int, offset *int, dir *string) (*model.EpisodePagination, error) {
	intID, err := strconv.ParseInt(obj.ID, 10, 64)
	if err != nil {
		return nil, err
	}

	itemIDs, err := common.GetFromLoaderForKey(ctx, r.FilteredLoaders(ctx).EpisodesLoader, int(intID))
	if err != nil {
		return nil, err
	}

	page := utils.Paginate(itemIDs, first, offset, dir)

	episodes, err := common.GetManyFromLoader(ctx, r.Loaders.EpisodeLoader, utils.PointerIntArrayToIntArray(page.Items))
	if err != nil {
		return nil, err
	}

	return &model.EpisodePagination{
		Total:  page.Total,
		First:  page.First,
		Offset: page.Offset,
		Items:  utils.MapWithCtx(ctx, episodes, model.EpisodeFrom),
	}, nil
}

// Image is the resolver for the image field.
func (r *showResolver) Image(ctx context.Context, obj *model.Show, style *model.ImageStyle) (*string, error) {
	e, err := common.GetFromLoaderByID(ctx, r.Loaders.ShowLoader, utils.AsInt(obj.ID))
	if err != nil {
		return nil, err
	}
	return imageOrFallback(ctx, e.Images, e.Image, style), nil
}

// EpisodeCount is the resolver for the episodeCount field.
func (r *showResolver) EpisodeCount(ctx context.Context, obj *model.Show) (int, error) {
	seasonIDs, err := common.GetFromLoaderForKey(ctx, r.FilteredLoaders(ctx).SeasonsLoader, utils.AsInt(obj.ID))
	if err != nil {
		return 0, err
	}
	el := r.FilteredLoaders(ctx).EpisodesLoader
	for _, id := range seasonIDs {
		el.Load(ctx, *id)
	}

	count := 0
	for _, id := range seasonIDs {
		episodeIDs, err := common.GetFromLoaderForKey(ctx, el, *id)
		if err != nil {
			return 0, err
		}
		count += len(episodeIDs)
	}
	return count, nil
}

// SeasonCount is the resolver for the seasonCount field.
func (r *showResolver) SeasonCount(ctx context.Context, obj *model.Show) (int, error) {
	seasonIDs, err := common.GetFromLoaderForKey(ctx, r.FilteredLoaders(ctx).SeasonsLoader, utils.AsInt(obj.ID))
	if err != nil {
		return 0, err
	}
	return len(seasonIDs), nil
}

// Seasons is the resolver for the seasons field.
func (r *showResolver) Seasons(ctx context.Context, obj *model.Show, first *int, offset *int, dir *string) (*model.SeasonPagination, error) {
	intID, err := strconv.ParseInt(obj.ID, 10, 64)
	if err != nil {
		return nil, err
	}

	itemIDs, err := common.GetFromLoaderForKey(ctx, r.FilteredLoaders(ctx).SeasonsLoader, int(intID))
	if err != nil {
		return nil, err
	}

	page := utils.Paginate(itemIDs, first, offset, dir)

	seasons, err := common.GetManyFromLoader(ctx, r.Loaders.SeasonLoader, utils.PointerIntArrayToIntArray(page.Items))
	if err != nil {
		return nil, err
	}
	return &model.SeasonPagination{
		Total:  page.Total,
		First:  page.First,
		Offset: page.Offset,
		Items:  utils.MapWithCtx(ctx, seasons, model.SeasonFrom),
	}, nil
}

// DefaultEpisode is the resolver for the defaultEpisode field.
func (r *showResolver) DefaultEpisode(ctx context.Context, obj *model.Show) (*model.Episode, error) {
	s, err := common.GetFromLoaderByID(ctx, r.Loaders.ShowLoader, utils.AsInt(obj.ID))
	if err != nil {
		return nil, err
	}
	ls := r.FilteredLoaders(ctx)
	eID, err := show.DefaultEpisodeID(ctx, ls, s)
	if err != nil || eID == nil {
		return nil, err
	}
	return r.QueryRoot().Episode(ctx, strconv.Itoa(*eID))
}

// Episode returns generated.EpisodeResolver implementation.
func (r *Resolver) Episode() generated.EpisodeResolver { return &episodeResolver{r} }

// Season returns generated.SeasonResolver implementation.
func (r *Resolver) Season() generated.SeasonResolver { return &seasonResolver{r} }

// Show returns generated.ShowResolver implementation.
func (r *Resolver) Show() generated.ShowResolver { return &showResolver{r} }

type episodeResolver struct{ *Resolver }
type seasonResolver struct{ *Resolver }
type showResolver struct{ *Resolver }
