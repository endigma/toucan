package resolvers

import (
	"github.com/endigma/toucan/_examples/basic/gen/toucan"
	"github.com/endigma/toucan/_examples/basic/models"
	log "github.com/sirupsen/logrus"
)

func (e *repositoryResolver) HasRoleOwner(user *models.User, repository *models.Repository) bool {
	return false
}

func (e *repositoryResolver) HasRoleEditor(user *models.User, repository *models.Repository) bool {
	return false
}

func (e *repositoryResolver) HasRoleViewer(user *models.User, repository *models.Repository) bool {
	log.Info("Checking if user is viewer of repository", user, repository)
	return user.Name == "Tom" && repository.Label == "Facebook"
}

func (e *repositoryResolver) HasAttributePublic(repository *models.Repository) bool {
	return repository.Public
}

type repositoryResolver struct{ *Resolver }

// User returns graph.UserResolver implementation.
func (r *Resolver) Repository() toucan.RepositoryResolver { return &repositoryResolver{r} }
