package resolvers

import (
	"context"

	"github.com/endigma/toucan/_examples/basic/gen/toucan"
	"github.com/endigma/toucan/_examples/basic/models"
	"github.com/endigma/toucan/decision"
)

func (e repositoryResolver) HasRoleOwner(context context.Context, actor *models.User, resource *models.Repository) decision.Decision {
	for _, role := range actor.Roles {
		if role.Repo == resource.ID && role.Role == "owner" {
			return decision.Allow("actor is viewer")
		}
	}

	return decision.Skip("no viewer role")
}

func (e repositoryResolver) HasRoleEditor(context context.Context, actor *models.User, resource *models.Repository) decision.Decision {
	for _, role := range actor.Roles {
		if role.Repo == resource.ID && role.Role == "editor" {
			return decision.Allow("actor is viewer")
		}
	}

	return decision.Skip("no viewer role")
}

func (e repositoryResolver) HasRoleViewer(context context.Context, actor *models.User, resource *models.Repository) decision.Decision {
	for _, role := range actor.Roles {
		if role.Repo == resource.ID && role.Role == "viewer" {
			return decision.Allow("actor is viewer")
		}
	}

	return decision.Skip("no viewer role")
}

func (e repositoryResolver) HasAttributePublic(context context.Context, resource *models.Repository) decision.Decision {
	if resource.Public {
		return decision.Allow("repository is public")
	} else {
		return decision.Skip("repository is private")
	}
}

type repositoryResolver struct{ *Resolver }

// User returns graph.UserResolver implementation.
func (r *Resolver) Repository() toucan.RepositoryResolver { return &repositoryResolver{r} }
