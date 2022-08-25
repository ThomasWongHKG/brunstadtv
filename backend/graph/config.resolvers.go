package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"time"

	"github.com/bcc-code/brunstadtv/backend/graph/generated"
	gqlmodel "github.com/bcc-code/brunstadtv/backend/graph/model"
)

// App is the resolver for the App field.
func (r *configResolver) App(ctx context.Context, obj *gqlmodel.Config, timestamp *string) (*gqlmodel.AppConfig, error) {
	from, err := timestampFromString(timestamp)
	if err != nil {
		return nil, err
	}
	conf, err := withCacheAndTimestamp(ctx, "app_config", r.Queries.GetAppConfig, time.Second*30, from)
	if err != nil {
		return nil, err
	}
	return &gqlmodel.AppConfig{
		Live:       conf.Live,
		MinVersion: conf.MinVersion,
	}, nil
}

// Config returns generated.ConfigResolver implementation.
func (r *Resolver) Config() generated.ConfigResolver { return &configResolver{r} }

type configResolver struct{ *Resolver }
