package main

import (
	"context"
	"database/sql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	cache "github.com/Code-Hex/go-generics-cache"
	"github.com/bcc-code/brunstadtv/backend/applications"
	"github.com/bcc-code/brunstadtv/backend/asset"
	"github.com/bcc-code/brunstadtv/backend/auth0"
	"github.com/bcc-code/brunstadtv/backend/common"
	graphadmin "github.com/bcc-code/brunstadtv/backend/graph/admin"
	graphadmingenerated "github.com/bcc-code/brunstadtv/backend/graph/admin/generated"
	graphapi "github.com/bcc-code/brunstadtv/backend/graph/api"
	graphapigenerated "github.com/bcc-code/brunstadtv/backend/graph/api/generated"
	graphpub "github.com/bcc-code/brunstadtv/backend/graph/public"
	graphpubgenerated "github.com/bcc-code/brunstadtv/backend/graph/public/generated"
	"github.com/bcc-code/brunstadtv/backend/items/collection"
	"github.com/bcc-code/brunstadtv/backend/items/episode"
	"github.com/bcc-code/brunstadtv/backend/items/page"
	"github.com/bcc-code/brunstadtv/backend/items/season"
	"github.com/bcc-code/brunstadtv/backend/items/section"
	"github.com/bcc-code/brunstadtv/backend/items/show"
	"github.com/bcc-code/brunstadtv/backend/members"
	"github.com/bcc-code/brunstadtv/backend/search"
	"github.com/bcc-code/brunstadtv/backend/signing"
	"github.com/bcc-code/brunstadtv/backend/sqlc"
	"github.com/bcc-code/brunstadtv/backend/user"
	"github.com/bcc-code/brunstadtv/backend/utils"
	"github.com/bcc-code/mediabank-bridge/log"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/graph-gophers/dataloader/v7"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/samber/lo"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"sort"
	"strings"
	"time"
)

var generalCache = cache.New[string, any]()

var rolesLoaderCache = map[string]*common.FilteredLoaders{}

func getLoadersForRoles(db *sql.DB, queries *sqlc.Queries, collectionLoader *dataloader.Loader[int, *common.Collection], roles []string) *common.FilteredLoaders {
	sort.Strings(roles)

	key := strings.Join(roles, "-")

	if loaders, ok := rolesLoaderCache[key]; ok {
		return loaders
	}

	loaders := &common.FilteredLoaders{
		EpisodesLoader: common.NewRelationBatchLoader(func(ctx context.Context, ids []int) ([]common.Relation[int, int], error) {
			return queries.GetEpisodeIDsForSeasonsWithRoles(ctx, ids, roles)
		}),
		SeasonsLoader: common.NewRelationBatchLoader(func(ctx context.Context, ids []int) ([]common.Relation[int, int], error) {
			return queries.GetSeasonIDsForShowsWithRoles(ctx, ids, roles)
		}),
		CollectionItemsLoader: common.NewListBatchLoader(func(ctx context.Context, ids []int) ([]common.CollectionItem, error) {
			return queries.GetItemsForCollectionsWithRoles(ctx, ids, roles)
		}, func(i common.CollectionItem) int {
			return i.CollectionID
		}),
		CollectionItemIDsLoader: collection.NewCollectionItemIdsLoader(db, collectionLoader, roles),
	}

	rolesLoaderCache[key] = loaders

	return loaders
}

func filteredLoaderFactory(db *sql.DB, queries *sqlc.Queries, collectionLoader *dataloader.Loader[int, *common.Collection]) func(ctx context.Context) *common.FilteredLoaders {
	return func(ctx context.Context) *common.FilteredLoaders {
		ginCtx, err := utils.GinCtx(ctx)
		var roles []string
		if err != nil {
			log.L.Error().Err(err).Msg("failed to get gin ctx from context")
			roles = []string{"unknown"}
		} else {
			roles = user.GetRolesFromCtx(ginCtx)
		}
		return getLoadersForRoles(db, queries, collectionLoader, roles)
	}
}

// Defining the Graphql handler
func graphqlHandler(db *sql.DB, queries *sqlc.Queries, loaders *common.BatchLoaders, searchService *search.Service, urlSigner *signing.Signer, config envConfig) gin.HandlerFunc {

	resolver := graphapi.Resolver{
		Queries:         queries,
		Loaders:         loaders,
		FilteredLoaders: filteredLoaderFactory(db, queries, loaders.CollectionLoader),
		SearchService:   searchService,
		APIConfig:       config.CDNConfig,
		URLSigner:       urlSigner,
	}

	tracer := &graphTracer{}

	// NewExecutableSchema and Config are in the generated.go file
	// Resolver is in the resolver.go file
	h := handler.NewDefaultServer(graphapigenerated.NewExecutableSchema(graphapigenerated.Config{Resolvers: &resolver}))
	h.Use(tracer)

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func publicGraphqlHandler(loaders *common.BatchLoaders) gin.HandlerFunc {
	resolver := graphpub.Resolver{
		Loaders: &graphpub.Loaders{
			EpisodeLoader: loaders.EpisodeLoader,
			SeasonLoader:  loaders.SeasonLoader,
			ShowLoader:    loaders.ShowLoader,
		},
	}

	tracer := &graphTracer{}

	// NewExecutableSchema and Config are in the generated.go file
	// Resolver is in the resolver.go file
	h := handler.NewDefaultServer(graphpubgenerated.NewExecutableSchema(graphpubgenerated.Config{Resolvers: &resolver}))
	h.Use(tracer)

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func adminGraphqlHandler(config envConfig, db *sql.DB, queries *sqlc.Queries, loaders *common.BatchLoaders) gin.HandlerFunc {

	resolver := graphadmin.Resolver{
		DB:      db,
		Queries: queries,
		Loaders: loaders,
	}

	// NewExecutableSchema and Config are in the generated.go file
	// Resolver is in the resolver.go file
	h := handler.NewDefaultServer(graphadmingenerated.NewExecutableSchema(graphadmingenerated.Config{Resolvers: &resolver}))

	directusSecret := config.Secrets.Directus
	if directusSecret == "" {
		log.L.Debug().Msg("No secret for Directus found in environment. Disabling endpoint")
		return func(c *gin.Context) {
			c.AbortWithStatus(404)
			return
		}
	}

	return func(c *gin.Context) {
		headerValue := c.GetHeader("x-api-key")
		if headerValue != directusSecret {
			c.AbortWithStatus(403)
			return
		}

		h.ServeHTTP(c.Writer, c.Request)
	}
}

// Defining the Playground handler
func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func getApplications(ctx context.Context, queries *sqlc.Queries) []common.Application {
	var key = "applications"
	cached, ok := generalCache.Get(key)
	if ok {
		return cached.([]common.Application)
	} else {
		apps, err := queries.ListApplications(ctx)
		if err != nil {
			panic(err)
		}
		// Cache with expiration in case the container lives too long
		generalCache.Set(key, apps, cache.WithExpiration(time.Minute*5))
		return apps
	}
}

func applicationFactory(queries *sqlc.Queries) func(ctx context.Context, code string) *common.Application {
	return func(ctx context.Context, code string) *common.Application {
		apps := getApplications(ctx, queries)

		app, found := lo.Find(apps, func(i common.Application) bool {
			return i.Code == strings.ToLower(strings.TrimSpace(code))
		})
		if found {
			return &app
		}
		app, found = lo.Find(apps, func(i common.Application) bool {
			return i.Default
		})
		if found {
			return &app
		}
		return nil
	}
}

func initBatchLoaders(db *sql.DB, queries *sqlc.Queries) *common.BatchLoaders {
	collectionLoader := common.NewBatchLoader(queries.GetCollections)

	return &common.BatchLoaders{
		// App
		ApplicationLoader:           common.NewBatchLoader(queries.GetApplications),
		ApplicationIDFromCodeLoader: common.NewConversionBatchLoader(queries.GetApplicationIDsForCodes),
		// Item
		PageLoader:           common.NewBatchLoader(queries.GetPages),
		PageIDFromCodeLoader: common.NewConversionBatchLoader(queries.GetPageIDsForCodes),
		SectionLoader:        common.NewBatchLoader(queries.GetSections),
		SectionLinksLoader: common.NewListBatchLoader(queries.GetLinksForSections, func(i common.SectionLink) int {
			return i.SectionID
		}),
		ShowLoader:           common.NewBatchLoader(queries.GetShows),
		SeasonLoader:         common.NewBatchLoader(queries.GetSeasons),
		EpisodeLoader:        common.NewBatchLoader(queries.GetEpisodes),
		EventLoader:          common.NewBatchLoader(queries.GetEvents),
		CalendarEntryLoader:  common.NewBatchLoader(queries.GetCalendarEntries),
		FilesLoader:          asset.NewBatchFilesLoader(*queries),
		StreamsLoader:        asset.NewBatchStreamsLoader(*queries),
		CollectionLoader:     collectionLoader,
		CollectionItemLoader: collection.NewItemListBatchLoader(*queries),
		// Relations
		SectionsLoader: common.NewRelationBatchLoader(queries.GetSectionIDsForPages),
		// Permissions
		ShowPermissionLoader:    show.NewPermissionLoader(*queries),
		SeasonPermissionLoader:  season.NewPermissionLoader(*queries),
		EpisodePermissionLoader: episode.NewPermissionLoader(*queries),
		PagePermissionLoader:    page.NewPermissionLoader(*queries),
		SectionPermissionLoader: section.NewPermissionLoader(*queries),
		FAQCategoryLoader:       common.NewBatchLoader(queries.GetFAQCategories),
		QuestionLoader:          common.NewBatchLoader(queries.GetQuestions),
		QuestionsLoader:         common.NewRelationBatchLoader(queries.GetQuestionIDsForCategories),
		// User Data
		ProfilesLoader: common.NewListBatchLoader(queries.GetProfilesForUserIDs, func(i common.Profile) string {
			return i.UserID
		}),
	}
}

func main() {
	ctx := context.Background()

	log.ConfigureGlobalLogger(zerolog.DebugLevel)
	log.L.Debug().Msg("Setting up tracing!")

	// Here you can get a tracedHttpClient if useful anywhere
	utils.MustSetupTracing()
	ctx, span := otel.Tracer("api/core").Start(ctx, "init")

	config := getEnvConfig()
	log.L.Debug().Str("DBConnString", config.DB.ConnectionString).Msg("Connection to DB")
	db, err := sql.Open("postgres", config.DB.ConnectionString)
	if err != nil {
		log.L.Panic().Err(err).Msg("Unable to connect to DB")
		return
	}

	urlSigner, err := signing.NewSigner(config.CDNConfig)
	if err != nil {
		log.L.Panic().Err(err).Msg("Unable to create URL signers")
		return
	}

	db.SetMaxIdleConns(2)
	// TODO: What makes sense here? We should gather some metrics over time
	db.SetMaxOpenConns(10)

	err = db.PingContext(ctx)
	if err != nil {
		log.L.Panic().Err(err).Msg("Ping failed")
		return
	}

	queries := sqlc.New(db)
	queries.SetImageCDNDomain(config.CDNConfig.ImageCDNDomain)

	loaders := initBatchLoaders(db, queries)

	authClient := auth0.New(config.Auth0)
	membersClient := members.New(config.Members, func(ctx context.Context) string {
		token, err := authClient.GetToken(ctx, config.Members.Domain)
		if err != nil {
			log.L.Panic().Err(err).Msg("Failed to retrieve token for members")
		}
		return token
	})

	log.L.Debug().Msg("Set up HTTP server")
	r := gin.Default()
	r.Use(utils.GinContextToContextMiddleware())
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"POST", "GET"},
		AllowHeaders:     []string{"content-type", "authorization", "accept-language"},
		AllowCredentials: true,
	}))
	r.Use(otelgin.Middleware("api")) // Open
	r.Use(authClient.ValidateToken())
	r.Use(user.NewUserMiddleware(queries, membersClient))
	r.Use(user.NewProfileMiddleware(queries, loaders))

	r.Use(applications.ApplicationMiddleware(applicationFactory(queries)))
	r.Use(applications.RoleMiddleware())

	searchService := search.New(db, config.Algolia)
	r.POST("/query", graphqlHandler(db, queries, loaders, searchService, urlSigner, config))

	r.GET("/", playgroundHandler())

	r.POST("/admin", adminGraphqlHandler(config, db, queries, loaders))

	r.POST("/public", publicGraphqlHandler(loaders))

	log.L.Debug().Msgf("connect to http://localhost:%s/ for GraphQL playground", config.Port)

	span.End()

	err = r.Run(":" + config.Port)
	if err != nil {
		return
	}
}
