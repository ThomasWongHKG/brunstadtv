package search

import (
	"context"
	"strconv"
	"strings"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/opt"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/bcc-code/brunstadtv/backend/common"
	"github.com/bcc-code/mediabank-bridge/log"
	"github.com/google/uuid"
	"github.com/graph-gophers/dataloader/v7"
	"github.com/samber/lo"
)

const tempIndexName = "temp"

// Reindex every supported collection
func (service *Service) Reindex(ctx context.Context) error {
	res, err := service.algoliaClient.CopyIndex(indexName, tempIndexName)
	_ = res.Wait()
	index := service.algoliaClient.InitIndex(tempIndexName)
	_, err = index.ClearObjects()
	if err != nil {
		return err
	}

	// Make sure we're not fetching from cache anywhere,
	// although that shouldn't be an issue, as we're priming on fetch anyway
	service.loaders.ShowLoader.ClearAll()
	service.loaders.ShowPermissionLoader.ClearAll()
	service.loaders.SeasonLoader.ClearAll()
	service.loaders.SeasonPermissionLoader.ClearAll()
	service.loaders.EpisodeLoader.ClearAll()
	service.loaders.EpisodePermissionLoader.ClearAll()

	// Makes it possible to filter in query, which fields you are searching on
	// Also configures hits per page
	primaryFields, err := service.getPrimaryTranslatedFields()
	if err != nil {
		return err
	}
	relationalFields, err := service.getRelationalTranslatedFields()
	if err != nil {
		return err
	}
	searchableAttributes := opt.SearchableAttributes(
		strings.Join(primaryFields, ", "),
		strings.Join(relationalFields, ", "),
		strings.Join(getFunctionalFields(), ", "),
	)
	languages, err := service.getLanguageKeys()
	if err != nil {
		return err
	}

	supportedLanguages := []string{"da", "de", "en", "es", "fi", "fr", "hu", "it", "nl", "no", "pl", "pt", "ro", "ru", "tr"}

	languages = lo.Filter(languages, func(l string, _ int) bool {
		return lo.Contains(supportedLanguages, l)
	})

	_, err = index.SetSettings(search.Settings{
		IndexLanguages:        opt.IndexLanguages(languages...),
		QueryLanguages:        opt.QueryLanguages(languages...),
		SearchableAttributes:  searchableAttributes,
		AttributesForFaceting: opt.AttributesForFaceting(service.getFilterFields()...),
		HitsPerPage:           opt.HitsPerPage(hitsPerPage),
	})
	if err != nil {
		return err
	}

	log.L.Debug().Str("collection", "shows").Msg("Indexing")
	err = service.indexShows(ctx, index)
	if err != nil {
		return err
	}
	log.L.Debug().Str("collection", "seasons").Msg("Indexing")
	err = service.indexSeasons(ctx, index)
	if err != nil {
		return err
	}
	log.L.Debug().Str("collection", "episodes").Msg("Indexing")
	err = service.indexEpisodes(ctx, index)
	if err != nil {
		return err
	}

	res, err = service.algoliaClient.MoveIndex(tempIndexName, indexName)
	if err != nil {
		return err
	}
	return res.Wait()
}

func (service *Service) indexShows(ctx context.Context, index *search.Index) error {
	return indexCollection[int, common.Show](
		ctx,
		service,
		index,
		service.loaders.ShowLoader,
		service.loaders.ShowPermissionLoader,
		service.queries.ListShows,
		service.showToSearchItem,
	)
}

func (service *Service) indexShow(ctx context.Context, id int) error {
	i, err := service.loaders.ShowLoader.Load(ctx, id)()
	if err != nil {
		return err
	}
	p, err := common.GetFromLoaderByID(ctx, service.loaders.ShowPermissionLoader, id)
	if err != nil {
		return err
	}
	return indexObject[int, common.Show](ctx, service, *i, p, service.showToSearchItem)
}

func (service *Service) indexSeasons(ctx context.Context, index *search.Index) error {
	return indexCollection[int, common.Season](
		ctx,
		service,
		index,
		service.loaders.SeasonLoader,
		service.loaders.SeasonPermissionLoader,
		service.queries.ListSeasons,
		service.seasonToSearchItem,
	)
}

func (service *Service) indexSeason(ctx context.Context, id int) error {
	i, err := service.loaders.SeasonLoader.Load(ctx, id)()
	if err != nil {
		return err
	}
	p, err := common.GetFromLoaderByID(ctx, service.loaders.SeasonPermissionLoader, id)
	if err != nil {
		return err
	}
	return indexObject[int, common.Season](ctx, service, *i, p, service.seasonToSearchItem)
}

func (service *Service) indexEpisodes(ctx context.Context, index *search.Index) error {
	return indexCollection[int, common.Episode](
		ctx,
		service,
		index,
		service.loaders.EpisodeLoader,
		service.loaders.EpisodePermissionLoader,
		service.queries.ListEpisodes,
		service.episodeToSearchItem,
	)
}

func (service *Service) indexEpisode(ctx context.Context, id int) error {
	i, err := service.loaders.EpisodeLoader.Load(ctx, id)()
	if err != nil {
		return err
	}
	p, err := common.GetFromLoaderByID(ctx, service.loaders.EpisodePermissionLoader, id)
	if err != nil {
		return err
	}
	return indexObject[int, common.Episode](ctx, service, *i, p, service.episodeToSearchItem)
}

type indexable[k comparable] interface {
	GetKey() k
	GetImage() uuid.NullUUID
}

func indexCollection[k comparable, t indexable[k]](
	ctx context.Context,
	service *Service,
	index *search.Index,
	loader *dataloader.Loader[k, *t],
	permissionLoader *dataloader.Loader[k, *common.Permissions[k]],
	factory func(context.Context) ([]t, error),
	converter func(context.Context, t) (searchItem, error),
) error {
	items, err := factory(ctx)
	if err != nil {
		return err
	}

	ids := lo.Map(items, func(i t, _ int) k {
		return i.GetKey()
	})

	imageIds := lo.Map(lo.Filter(items, func(i t, _ int) bool {
		return i.GetImage().Valid
	}), func(i t, _ int) uuid.UUID {
		return i.GetImage().UUID
	})

	permissionLoader.LoadMany(ctx, ids)()
	service.loaders.ImageLoader.LoadMany(ctx, imageIds)()

	var searchItems []searchObject
	pushItems := func(force bool) error {
		if len(searchItems) > 200 || (force && len(searchItems) > 0) {
			_, err := index.SaveObjects(searchItems)
			if err != nil {
				return err
			}
			searchItems = []searchObject{}
		}
		return nil
	}

	for _, i := range items {
		p := i
		loader.Prime(ctx, p.GetKey(), &p)

		item, err := converter(ctx, p)
		if err != nil {
			return err
		}

		perm, err := common.GetFromLoaderByID(ctx, permissionLoader, i.GetKey())
		if err != nil {
			return err
		}

		item.assignVisibility(perm.Availability)
		item.assignRoles(perm.Roles)
		err = item.assignImage(ctx, service.loaders, p)
		if err != nil {
			return err
		}

		searchItems = append(searchItems, item.toSearchObject())
		err = pushItems(false)
	}
	return pushItems(true)
}

func indexObject[k comparable, t indexable[k]](
	ctx context.Context,
	service *Service,
	obj t,
	perms *common.Permissions[k],
	converter func(context.Context, t) (searchItem, error),
) error {
	item, err := converter(ctx, obj)
	if err != nil {
		return err
	}

	item.assignVisibility(perms.Availability)
	item.assignRoles(perms.Roles)
	err = item.assignImage(ctx, service.loaders, obj)
	if err != nil {
		return err
	}
	_, err = service.index.SaveObject(item)
	return err
}

var supportedCollections = []string{
	"shows",
	"seasons",
	"episodes",
}

// DeleteModel from index by collection and id
func (service *Service) DeleteModel(_ context.Context, collection string, id int) error {

	if !lo.Contains(supportedCollections, collection) {
		// no reason to send a request if the collection isn't supported
		return nil
	}
	_, err := service.index.DeleteObject(collection + "-" + strconv.Itoa(id))
	return err
}

// IndexModel by collection and id
func (service *Service) IndexModel(ctx context.Context, collection string, id int) (err error) {
	// Clearing the loaders of cached instances
	// and indexing to the search engine
	log.L.Debug().Str("collection", collection).Int("id", id).Msg("Indexing item")
	switch collection {
	case "shows":
		service.loaders.ShowLoader.Clear(ctx, id)
		return service.indexShow(ctx, id)
	case "seasons":
		service.loaders.SeasonLoader.Clear(ctx, id)
		return service.indexSeason(ctx, id)
	case "episodes":
		service.loaders.EpisodeLoader.Clear(ctx, id)
		return service.indexEpisode(ctx, id)
	}
	// no reason to return errors, as we know quite well what is supported
	return nil
}
