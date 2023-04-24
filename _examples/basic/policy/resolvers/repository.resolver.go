package resolvers

import (
	"context"

	"github.com/endigma/toucan/_examples/basic/gen/toucan"
	"github.com/endigma/toucan/_examples/basic/models"
)

func (e *repositoryResolver) HasRole(ctx context.Context, user *models.User, role toucan.RepositoryRole, repository *models.Repository) bool {
	return false
}

func (e *repositoryResolver) HasAttribute(ctx context.Context, attr toucan.RepositoryAttribute, repository *models.Repository) bool {
	return false
}

type repositoryResolver struct{ *Resolver }

// User returns graph.UserResolver implementation.
func (r *Resolver) Repository() toucan.RepositoryResolver { return &repositoryResolver{r} }
