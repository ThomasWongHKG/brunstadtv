package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"time"

	"github.com/bcc-code/brunstadtv/backend/common"
	"github.com/bcc-code/brunstadtv/backend/graph/api/generated"
	"github.com/bcc-code/brunstadtv/backend/graph/api/model"
	"github.com/bcc-code/brunstadtv/backend/utils"
	"github.com/samber/lo"
)

// Categories is the resolver for the categories field.
func (r *fAQResolver) Categories(ctx context.Context, obj *model.Faq, first *int, offset *int) (*model.FAQCategoryPagination, error) {
	items, err := withCache(ctx, "categories", r.Queries.ListFAQCategories, time.Minute*5)
	if err != nil {
		return nil, err
	}

	cats := utils.MapWithCtx(ctx, lo.Map(items, func(i common.FAQCategory, _ int) *common.FAQCategory {
		return &i
	}), model.FAQCategoryFrom)

	page := utils.Paginate(cats, first, offset, nil)

	return &model.FAQCategoryPagination{
		Total:  page.Total,
		First:  page.First,
		Offset: page.Offset,
		Items:  page.Items,
	}, nil
}

// Category is the resolver for the category field.
func (r *fAQResolver) Category(ctx context.Context, obj *model.Faq, id string) (*model.FAQCategory, error) {
	return resolverForIntID(ctx, &itemLoaders[int, common.FAQCategory]{
		Item: r.Loaders.FAQCategoryLoader,
	}, id, model.FAQCategoryFrom)
}

// Question is the resolver for the question field.
func (r *fAQResolver) Question(ctx context.Context, obj *model.Faq, id string) (*model.Question, error) {
	return resolverForIntID(ctx, &itemLoaders[int, common.Question]{
		Item: r.Loaders.QuestionLoader,
	}, id, model.QuestionFrom)
}

// Questions is the resolver for the questions field.
func (r *fAQCategoryResolver) Questions(ctx context.Context, obj *model.FAQCategory, first *int, offset *int) (*model.QuestionPagination, error) {
	items, err := itemsResolverForIntID(ctx, &itemLoaders[int, common.Question]{
		Item: r.Loaders.QuestionLoader,
	}, r.Loaders.QuestionsLoader, obj.ID, model.QuestionFrom)
	if err != nil {
		return nil, err
	}

	page := utils.Paginate(items, first, offset, nil)

	return &model.QuestionPagination{
		Total:  page.Total,
		First:  page.First,
		Offset: page.Offset,
		Items:  page.Items,
	}, nil
}

// Category is the resolver for the category field.
func (r *questionResolver) Category(ctx context.Context, obj *model.Question) (*model.FAQCategory, error) {
	return r.FAQ().Category(ctx, nil, obj.Category.ID)
}

// FAQ returns generated.FAQResolver implementation.
func (r *Resolver) FAQ() generated.FAQResolver { return &fAQResolver{r} }

// FAQCategory returns generated.FAQCategoryResolver implementation.
func (r *Resolver) FAQCategory() generated.FAQCategoryResolver { return &fAQCategoryResolver{r} }

// Question returns generated.QuestionResolver implementation.
func (r *Resolver) Question() generated.QuestionResolver { return &questionResolver{r} }

type fAQResolver struct{ *Resolver }
type fAQCategoryResolver struct{ *Resolver }
type questionResolver struct{ *Resolver }
