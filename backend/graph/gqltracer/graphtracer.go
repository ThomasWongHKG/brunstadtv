package gqltracer

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/davecgh/go-spew/spew"
	"go.opentelemetry.io/otel"
)

type GraphTracer struct {
}

// ExtensionName name of extension
func (t *GraphTracer) ExtensionName() string {
	return "GraphTracer"
}

// InterceptField intercepts
func (t *GraphTracer) InterceptField(ctx context.Context, next graphql.Resolver) (any, error) {
	field := graphql.GetFieldContext(ctx)

	fieldName := field.Field.Name
	if fieldName == "" {
		fieldName = "unknown"
	}
	ctx, span := otel.Tracer("graph-field").Start(ctx, fieldName)
	defer span.End()
	return next(ctx)
}

// InterceptOperation intercepts
func (t *GraphTracer) InterceptOperation(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
	op := graphql.GetOperationContext(ctx)
	opName := op.OperationName
	if opName == "" {
		opName = "unknown"
	}
	ctx, span := otel.Tracer("graph-operation").Start(ctx, opName)
	defer span.End()
	return next(ctx)
}

// Validate is unnecessary in our case
func (t *GraphTracer) Validate(_ graphql.ExecutableSchema) error {
	return nil
}
