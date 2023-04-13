package resolvers

import (
	"github.com/endigma/toucan/_examples/basic/gen/toucan"
	"github.com/endigma/toucan/_examples/basic/models"
)

func (u *userResolver) HasRoleAdmin(actor *models.User, resource *models.User) bool {
	return false
}

func (u *userResolver) HasRoleSelf(actor *models.User, resource *models.User) bool {
	return false
}

func (u *userResolver) HasRoleViewer(actor *models.User, resource *models.User) bool {
	return u.HasRoleAdmin(actor, resource)
}

func (u *userResolver) HasAttributePublic(resource *models.User) bool {
	return false
}

type userResolver struct{ *Resolver }

// User returns graph.UserResolver implementation.
func (r *Resolver) User() toucan.UserResolver { return &userResolver{r} }
