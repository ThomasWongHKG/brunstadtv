package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/bcc-code/brunstadtv/backend/graph/public/generated"
	"github.com/bcc-code/brunstadtv/backend/graph/public/model"
)

func (r *seasonResolver) Show(ctx context.Context, obj *model.Season) (*model.Show, error) {
	return r.QueryRoot().Show(ctx, obj.Show.ID)
}

// Season returns generated.SeasonResolver implementation.
func (r *Resolver) Season() generated.SeasonResolver { return &seasonResolver{r} }

type seasonResolver struct{ *Resolver }
