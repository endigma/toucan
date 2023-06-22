package resolvers

import (
	"context"

	"github.com/endigma/toucan/_examples/basic/gen/toucan"
	"github.com/endigma/toucan/_examples/basic/models"
)

func (e userResolver) HasRoleAdmin(context context.Context, actor *models.User, resource *models.User) (bool, error) {
	return true, nil
}

func (e userResolver) HasRoleSelf(context context.Context, actor *models.User, resource *models.User) (bool, error) {
	return true, nil
}

func (e userResolver) HasRoleViewer(context context.Context, actor *models.User, resource *models.User) (bool, error) {
	return true, nil
}

type userResolver struct{ *Resolver }

// User returns graph.UserResolver implementation.
func (r *Resolver) User() toucan.UserResolver { return &userResolver{r} }
