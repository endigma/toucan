package resolvers

import (
	"context"

	"github.com/endigma/toucan/_examples/basic/gen/toucan"
	"github.com/endigma/toucan/_examples/basic/models"
	"github.com/endigma/toucan/decision"
)

func (e userResolver) HasRoleAdmin(context context.Context, actor *models.User, resource *models.User) decision.Decision {
	return decision.True("")
}

func (e userResolver) HasRoleSelf(context context.Context, actor *models.User, resource *models.User) decision.Decision {
	return decision.True("")
}

func (e userResolver) HasRoleViewer(context context.Context, actor *models.User, resource *models.User) decision.Decision {
	return decision.True("")
}

type userResolver struct{ *Resolver }

// User returns graph.UserResolver implementation.
func (r *Resolver) User() toucan.UserResolver { return &userResolver{r} }
