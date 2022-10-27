package main

import (
	"context"
	"database/sql"
	"sort"
	"strings"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	cache "github.com/Code-Hex/go-generics-cache"
	awsSDKConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/bcc-code/brunstadtv/backend/applications"
	"github.com/bcc-code/brunstadtv/backend/asset"
	"github.com/bcc-code/brunstadtv/backend/auth0"
	"github.com/bcc-code/brunstadtv/backend/batchloaders"
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
	"github.com/bcc-code/brunstadtv/backend/version"
	"github.com/bcc-code/mediabank-bridge/log"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"github.com/google/uuid"
	"github.com/graph-gophers/dataloader/v7"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/samber/lo"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
)

var generalCache = cache.New[string, any]()

var rolesLoaderCache = cache.New[string, *common.FilteredLoaders]()

func getLoadersForRoles(db *sql.DB, queries *sqlc.Queries, collectionLoader *dataloader.Loader[int, *common.Collection], roles []string) *common.FilteredLoaders {
	sort.Strings(roles)

	key := strings.Join(roles, "-")

	if loaders, ok := rolesLoaderCache.Get(key); ok {
		return loaders
	}

	rq := queries.RoleQueries(roles)

	loaders := &common.FilteredLoaders{
		ShowFilterLoader:    batchloaders.NewFilterLoader(rq.GetShowIDsWithRoles),
		SeasonFilterLoader:  batchloaders.NewFilterLoader(rq.GetSeasonIDsWithRoles),
		EpisodeFilterLoader: batchloaders.NewFilterLoader(rq.GetEpisodeIDsWithRoles),
		SeasonsLoader:       batchloaders.NewRelationLoader(rq.GetSeasonIDsForShowsWithRoles),
		SectionsLoader:      batchloaders.NewRelationLoader(rq.GetSectionIDsForPagesWithRoles),
		EpisodesLoader:      batchloaders.NewRelationLoader(rq.GetEpisodeIDsForSeasonsWithRoles),
		CollectionItemsLoader: batchloaders.NewListLoader(rq.GetItemsForCollectionsWithRoles, func(i common.CollectionItem) int {
			return i.CollectionID
		}),
		CollectionItemIDsLoader: collection.NewCollectionItemIdsLoader(db, collectionLoader, roles),
	}

	rolesLoaderCache.Set(key, loaders)

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

var profilesLoaderCache = cache.New[uuid.UUID, *common.ProfileLoaders]()

func getLoadersForProfile(queries *sqlc.Queries, profileID uuid.UUID) *common.ProfileLoaders {
	if loaders, ok := profilesLoaderCache.Get(profileID); ok {
		return loaders
	}

	profileQueries := queries.ProfileQueries(profileID)
	loaders := &common.ProfileLoaders{
		ProgressLoader: batchloaders.New(profileQueries.GetProgressForEpisodes, batchloaders.WithMemoryCache(time.Second*5)),
	}

	profilesLoaderCache.Set(profileID, loaders, cache.WithExpiration(time.Minute*5))

	return loaders
}

func profileLoaderFactory(queries *sqlc.Queries) func(ctx context.Context) *common.ProfileLoaders {
	return func(ctx context.Context) *common.ProfileLoaders {
		ginCtx, err := utils.GinCtx(ctx)
		if err != nil {
			return nil
		}
		p := user.GetProfileFromCtx(ginCtx)
		if p == nil {
			return nil
		}
		return getLoadersForProfile(queries, p.ID)
	}
}

// Defining the Graphql handler
func graphqlHandler(
	db *sql.DB,
	queries *sqlc.Queries,
	loaders *common.BatchLoaders,
	searchService *search.Service,
	urlSigner *signing.Signer,
	config envConfig,
	s3client *s3.Client,
) gin.HandlerFunc {

	resolver := graphapi.Resolver{
		Queries:         queries,
		Loaders:         loaders,
		FilteredLoaders: filteredLoaderFactory(db, queries, loaders.CollectionLoader),
		ProfileLoaders:  profileLoaderFactory(queries),
		SearchService:   searchService,
		APIConfig:       config.CDNConfig,
		AWSConfig:       config.AWS,
		URLSigner:       urlSigner,
		S3Client:        s3client,
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
		// Cache with expiration in case the container lives too long.
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

func initBatchLoaders(queries *sqlc.Queries) *common.BatchLoaders {
	collectionLoader := batchloaders.New(queries.GetCollections).Loader

	return &common.BatchLoaders{
		// App
		ApplicationLoader:           batchloaders.NewLoader(queries.GetApplications),
		ApplicationIDFromCodeLoader: batchloaders.NewConversionLoader[string, int](queries.GetApplicationIDsForCodes),
		// Item
		PageLoader:           batchloaders.New(queries.GetPages).Loader,
		PageIDFromCodeLoader: batchloaders.NewConversionLoader[string, int](queries.GetPageIDsForCodes),
		SectionLoader:        batchloaders.New(queries.GetSections).Loader,
		ShowLoader:           batchloaders.New(queries.GetShows).Loader,
		SeasonLoader:         batchloaders.New(queries.GetSeasons).Loader,
		EpisodeLoader:        batchloaders.New(queries.GetEpisodes).Loader,
		LinkLoader:           batchloaders.New(queries.GetLinks).Loader,
		EventLoader:          batchloaders.New(queries.GetEvents).Loader,
		CalendarEntryLoader:  batchloaders.New(queries.GetCalendarEntries).Loader,
		FilesLoader:          asset.NewBatchFilesLoader(*queries),
		StreamsLoader:        asset.NewBatchStreamsLoader(*queries),
		CollectionLoader:     collectionLoader,
		CollectionItemLoader: collection.NewItemListBatchLoader(*queries),
		// Relations
		SectionsLoader: batchloaders.NewRelationLoader(queries.GetSectionIDsForPages),
		// Permissions
		ShowPermissionLoader:    show.NewPermissionLoader(*queries),
		SeasonPermissionLoader:  season.NewPermissionLoader(*queries),
		EpisodePermissionLoader: episode.NewPermissionLoader(*queries),
		PagePermissionLoader:    page.NewPermissionLoader(*queries),
		SectionPermissionLoader: section.NewPermissionLoader(*queries),
		FAQCategoryLoader:       batchloaders.NewLoader(queries.GetFAQCategories),
		QuestionLoader:          batchloaders.NewLoader(queries.GetQuestions),
		QuestionsLoader:         batchloaders.NewRelationLoader(queries.GetQuestionIDsForCategories),
		MessageGroupLoader:      batchloaders.NewLoader(queries.GetMessageGroups),
		// User Data
		ProfilesLoader: batchloaders.NewListLoader(queries.GetProfilesForUserIDs, func(i common.Profile) string {
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

	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Redis.Address,
		Password: config.Redis.Password,
		Username: config.Redis.Username,
		DB:       config.Redis.Database,
	})

	status := rdb.Ping(ctx)
	if status.Err() != nil {
		log.L.Panic().Err(status.Err()).Msg("Failed to ping redis database")
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

	loaders := initBatchLoaders(queries)

	authClient := auth0.New(config.Auth0)
	membersClient := members.New(config.Members, func(ctx context.Context) string {
		token, err := authClient.GetToken(ctx, config.Members.Domain)
		if err != nil {
			log.L.Panic().Err(err).Msg("Failed to retrieve token for members")
		}
		return token
	})

	awsConfig, err := awsSDKConfig.LoadDefaultConfig(ctx)
	if err != nil {
		log.L.Panic().Err(err).Msg("Failed to configure AWS SDK")
	}
	awsConfig.Region = config.AWS.Region

	s3Client := s3.NewFromConfig(awsConfig)

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
	r.Use(user.NewProfileMiddleware(queries, rdb))

	r.Use(applications.ApplicationMiddleware(applicationFactory(queries)))
	r.Use(applications.RoleMiddleware())

	searchService := search.New(db, config.Algolia)

	gqlHandler := graphqlHandler(db, queries, loaders, searchService, urlSigner, config, s3Client)
	r.POST("/query", gqlHandler)

	r.GET("/", playgroundHandler())

	r.POST("/admin", adminGraphqlHandler(config, db, queries, loaders))

	r.POST("/public", publicGraphqlHandler(loaders))

	r.GET("/versionz", version.GinHandler)

	log.L.Debug().Msgf("connect to http://localhost:%s/ for GraphQL playground", config.Port)

	span.End()

	err = r.Run(":" + config.Port)
	if err != nil {
		return
	}
}
