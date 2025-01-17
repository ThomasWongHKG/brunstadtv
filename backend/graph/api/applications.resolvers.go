package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/bcc-code/brunstadtv/backend/graph/api/generated"
	gqlmodel "github.com/bcc-code/brunstadtv/backend/graph/api/model"
)

// Page is the resolver for the page field.
func (r *applicationResolver) Page(ctx context.Context, obj *gqlmodel.Application) (*gqlmodel.Page, error) {
	if obj.Page != nil {
		return r.QueryRoot().Page(ctx, &obj.Page.ID, nil)
	}
	return nil, nil
}

// SearchPage is the resolver for the searchPage field.
func (r *applicationResolver) SearchPage(ctx context.Context, obj *gqlmodel.Application) (*gqlmodel.Page, error) {
	if obj.SearchPage != nil {
		return r.QueryRoot().Page(ctx, &obj.SearchPage.ID, nil)
	}
	return nil, nil
}

// Application returns generated.ApplicationResolver implementation.
func (r *Resolver) Application() generated.ApplicationResolver { return &applicationResolver{r} }

type applicationResolver struct{ *Resolver }
