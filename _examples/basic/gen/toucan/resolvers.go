// Code generated by toucan. DO NOT EDIT.
package toucan

import (
	"context"
	models "github.com/endigma/toucan/_examples/basic/models"
	decision "github.com/endigma/toucan/decision"
)

// Resolver for resource `global`
type GlobalResolver interface {
	HasRoleAdmin(ctx context.Context, actor *models.User) decision.Decision
	HasAttributeProfilesArePublic(ctx context.Context) decision.Decision
}

// Resolver for resource `repository`
type RepositoryResolver interface {
	CacheKey(resource *models.Repository) string

	HasRoleOwner(ctx context.Context, actor *models.User, resource *models.Repository) decision.Decision
	HasRoleEditor(ctx context.Context, actor *models.User, resource *models.Repository) decision.Decision
	HasRoleViewer(ctx context.Context, actor *models.User, resource *models.Repository) decision.Decision
	HasAttributePublic(ctx context.Context, resource *models.Repository) decision.Decision
}

// Resolver for resource `user`
type UserResolver interface {
	CacheKey(resource *models.User) string

	HasRoleAdmin(ctx context.Context, actor *models.User, resource *models.User) decision.Decision
	HasRoleSelf(ctx context.Context, actor *models.User, resource *models.User) decision.Decision
	HasRoleViewer(ctx context.Context, actor *models.User, resource *models.User) decision.Decision
}

// Root Resolver
type Resolver interface {
	CacheKey(actor *models.User) string

	Global() GlobalResolver
	Repository() RepositoryResolver
	User() UserResolver
}
