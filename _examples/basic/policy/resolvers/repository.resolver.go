package resolvers

import (
	"context"

	"github.com/endigma/toucan/_examples/basic/gen/toucan"
	"github.com/endigma/toucan/_examples/basic/models"
)

func (e repositoryResolver) HasRoleOwner(context context.Context, actor *models.User, resource *models.Repository) (bool, error) {
	for _, role := range actor.Roles {
		if role.Repo == resource.ID && role.Role == "owner" {
			return true, nil
		}
	}

	return false, nil
}

func (e repositoryResolver) HasRoleEditor(context context.Context, actor *models.User, resource *models.Repository) (bool, error) {
	for _, role := range actor.Roles {
		if role.Repo == resource.ID && role.Role == "editor" {
			return true, nil
		}
	}

	return false, nil
}

func (e repositoryResolver) HasRoleViewer(context context.Context, actor *models.User, resource *models.Repository) (bool, error) {
	for _, role := range actor.Roles {
		if role.Repo == resource.ID && role.Role == "viewer" {
			return true, nil
		}
	}

	return false, nil
}

func (e repositoryResolver) HasAttributePublic(context context.Context, resource *models.Repository) (bool, error) {
	return resource.Public, nil
}

type repositoryResolver struct{ *Resolver }

// User returns graph.UserResolver implementation.
func (r *Resolver) Repository() toucan.RepositoryResolver { return &repositoryResolver{r} }
