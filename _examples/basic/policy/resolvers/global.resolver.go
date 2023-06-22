package resolvers

import (
	"context"

	"github.com/endigma/toucan/_examples/basic/gen/toucan"
	"github.com/endigma/toucan/_examples/basic/models"
)

func (g globalResolver) HasAttributeProfilesArePublic(context context.Context) (bool, error) {
	return true, nil
}

func (g globalResolver) HasRoleAdmin(context context.Context, actor *models.User) (bool, error) {
	return actor.GlobalAdmin, nil
}

type globalResolver struct{ *Resolver }

// Global returns graph.GlobalResolver implementation.
func (r *Resolver) Global() toucan.GlobalResolver { return &globalResolver{r} }
