package resolvers

import (
	"context"

	"github.com/endigma/toucan/_examples/basic/gen/toucan"
	"github.com/endigma/toucan/_examples/basic/models"
	"github.com/endigma/toucan/decision"
)

func (g globalResolver) HasAttributeProfilesArePublic(context context.Context) decision.Decision {
	return decision.False("attribute profiles are public")
}

func (g globalResolver) HasRoleAdmin(context context.Context, actor *models.User) decision.Decision {
	if actor.GlobalAdmin {
		return decision.True("actor is admin")
	}

	return decision.False("no admin role")
}

type globalResolver struct{ *Resolver }

// Global returns graph.GlobalResolver implementation.
func (r *Resolver) Global() toucan.GlobalResolver { return &globalResolver{r} }
