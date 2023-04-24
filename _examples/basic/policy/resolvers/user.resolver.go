package resolvers

import (
	"context"

	"github.com/endigma/toucan/_examples/basic/gen/toucan"
	"github.com/endigma/toucan/_examples/basic/models"
)

func (e *userResolver) HasRole(ctx context.Context, user *models.User, role toucan.UserRole, repository *models.User) bool {
	return false
}

func (e *userResolver) HasAttribute(ctx context.Context, attr toucan.UserAttribute, repository *models.User) bool {
	return false
}

type userResolver struct{ *Resolver }

// User returns graph.UserResolver implementation.
func (r *Resolver) User() toucan.UserResolver { return &userResolver{r} }
