// Code generated by toucan. DO NOT EDIT.
package toucan

import (
	"context"
	models "github.com/endigma/toucan/_examples/basic/models"
	cache "github.com/endigma/toucan/cache"
	decision "github.com/endigma/toucan/decision"
)

func (a Authorizer) AuthorizeRepository(ctx context.Context, actor *models.User, action RepositoryPermission, resource *models.Repository) decision.Decision {
	resolver := a.Repository()

	if !action.Valid() {
		return decision.Error(ErrInvalidRepositoryPermission)
	}

	if resource != nil {
		switch action {
		case RepositoryPermissionRead:
			// Source: attribute - Public
			if result := cache.QueryOr(ctx, cache.CacheKey{
				ActorKey:    actor.ToucanKey(),
				Resource:    "repository",
				ResourceKey: resource.ToucanKey(),
				SourceType:  "attribute",
				SourceName:  "Public",
			}, func() decision.Decision {
				return resolver.HasAttributePublic(ctx, resource)
			}); result.Allow {
				return result
			}

		}
	}

	if resource != nil && actor != nil {
		switch action {
		case RepositoryPermissionRead:
			// Source: role - Owner
			if result := cache.QueryOr(ctx, cache.CacheKey{
				ActorKey:    actor.ToucanKey(),
				Resource:    "repository",
				ResourceKey: resource.ToucanKey(),
				SourceType:  "role",
				SourceName:  "Owner",
			}, func() decision.Decision {
				return resolver.HasRoleOwner(ctx, actor, resource)
			}); result.Allow {
				return result
			}

			// Source: role - Editor
			if result := cache.QueryOr(ctx, cache.CacheKey{
				ActorKey:    actor.ToucanKey(),
				Resource:    "repository",
				ResourceKey: resource.ToucanKey(),
				SourceType:  "role",
				SourceName:  "Editor",
			}, func() decision.Decision {
				return resolver.HasRoleEditor(ctx, actor, resource)
			}); result.Allow {
				return result
			}

			// Source: role - Viewer
			if result := cache.QueryOr(ctx, cache.CacheKey{
				ActorKey:    actor.ToucanKey(),
				Resource:    "repository",
				ResourceKey: resource.ToucanKey(),
				SourceType:  "role",
				SourceName:  "Viewer",
			}, func() decision.Decision {
				return resolver.HasRoleViewer(ctx, actor, resource)
			}); result.Allow {
				return result
			}

		case RepositoryPermissionPush:
			// Source: role - Owner
			if result := cache.QueryOr(ctx, cache.CacheKey{
				ActorKey:    actor.ToucanKey(),
				Resource:    "repository",
				ResourceKey: resource.ToucanKey(),
				SourceType:  "role",
				SourceName:  "Owner",
			}, func() decision.Decision {
				return resolver.HasRoleOwner(ctx, actor, resource)
			}); result.Allow {
				return result
			}

			// Source: role - Editor
			if result := cache.QueryOr(ctx, cache.CacheKey{
				ActorKey:    actor.ToucanKey(),
				Resource:    "repository",
				ResourceKey: resource.ToucanKey(),
				SourceType:  "role",
				SourceName:  "Editor",
			}, func() decision.Decision {
				return resolver.HasRoleEditor(ctx, actor, resource)
			}); result.Allow {
				return result
			}

		case RepositoryPermissionDelete:
			// Source: role - Owner
			if result := cache.QueryOr(ctx, cache.CacheKey{
				ActorKey:    actor.ToucanKey(),
				Resource:    "repository",
				ResourceKey: resource.ToucanKey(),
				SourceType:  "role",
				SourceName:  "Owner",
			}, func() decision.Decision {
				return resolver.HasRoleOwner(ctx, actor, resource)
			}); result.Allow {
				return result
			}

		case RepositoryPermissionSnakeCase:
			// Source: role - Owner
			if result := cache.QueryOr(ctx, cache.CacheKey{
				ActorKey:    actor.ToucanKey(),
				Resource:    "repository",
				ResourceKey: resource.ToucanKey(),
				SourceType:  "role",
				SourceName:  "Owner",
			}, func() decision.Decision {
				return resolver.HasRoleOwner(ctx, actor, resource)
			}); result.Allow {
				return result
			}

		}
	}

	return decision.False("unmatched")
}

func (a Authorizer) FilterRepository(ctx context.Context, actor *models.User, action RepositoryPermission, resources []*models.Repository) ([]*models.Repository, error) {
	if !action.Valid() {
		return nil, ErrInvalidRepositoryPermission
	}

	var allowedResolvers []*models.Repository
	for _, resource := range resources {
		result := a.AuthorizeRepository(ctx, actor, action, resource)
		if result.Allow {
			allowedResolvers = append(allowedResolvers, resource)
		}
	}

	return allowedResolvers, nil
}

func (a Authorizer) AuthorizeUser(ctx context.Context, actor *models.User, action UserPermission, resource *models.User) decision.Decision {
	resolver := a.User()

	if !action.Valid() {
		return decision.Error(ErrInvalidUserPermission)
	}

	if resource != nil && actor != nil {
		switch action {
		case UserPermissionRead:
			// Source: role - Admin
			if result := cache.QueryOr(ctx, cache.CacheKey{
				ActorKey:    actor.ToucanKey(),
				Resource:    "user",
				ResourceKey: resource.ToucanKey(),
				SourceType:  "role",
				SourceName:  "Admin",
			}, func() decision.Decision {
				return resolver.HasRoleAdmin(ctx, actor, resource)
			}); result.Allow {
				return result
			}

			// Source: role - Self
			if result := cache.QueryOr(ctx, cache.CacheKey{
				ActorKey:    actor.ToucanKey(),
				Resource:    "user",
				ResourceKey: resource.ToucanKey(),
				SourceType:  "role",
				SourceName:  "Self",
			}, func() decision.Decision {
				return resolver.HasRoleSelf(ctx, actor, resource)
			}); result.Allow {
				return result
			}

			// Source: role - Viewer
			if result := cache.QueryOr(ctx, cache.CacheKey{
				ActorKey:    actor.ToucanKey(),
				Resource:    "user",
				ResourceKey: resource.ToucanKey(),
				SourceType:  "role",
				SourceName:  "Viewer",
			}, func() decision.Decision {
				return resolver.HasRoleViewer(ctx, actor, resource)
			}); result.Allow {
				return result
			}

		case UserPermissionWrite:
			// Source: role - Admin
			if result := cache.QueryOr(ctx, cache.CacheKey{
				ActorKey:    actor.ToucanKey(),
				Resource:    "user",
				ResourceKey: resource.ToucanKey(),
				SourceType:  "role",
				SourceName:  "Admin",
			}, func() decision.Decision {
				return resolver.HasRoleAdmin(ctx, actor, resource)
			}); result.Allow {
				return result
			}

			// Source: role - Self
			if result := cache.QueryOr(ctx, cache.CacheKey{
				ActorKey:    actor.ToucanKey(),
				Resource:    "user",
				ResourceKey: resource.ToucanKey(),
				SourceType:  "role",
				SourceName:  "Self",
			}, func() decision.Decision {
				return resolver.HasRoleSelf(ctx, actor, resource)
			}); result.Allow {
				return result
			}

		case UserPermissionDelete:
			// Source: role - Admin
			if result := cache.QueryOr(ctx, cache.CacheKey{
				ActorKey:    actor.ToucanKey(),
				Resource:    "user",
				ResourceKey: resource.ToucanKey(),
				SourceType:  "role",
				SourceName:  "Admin",
			}, func() decision.Decision {
				return resolver.HasRoleAdmin(ctx, actor, resource)
			}); result.Allow {
				return result
			}

		}
	}

	return decision.False("unmatched")
}

func (a Authorizer) FilterUser(ctx context.Context, actor *models.User, action UserPermission, resources []*models.User) ([]*models.User, error) {
	if !action.Valid() {
		return nil, ErrInvalidUserPermission
	}

	var allowedResolvers []*models.User
	for _, resource := range resources {
		result := a.AuthorizeUser(ctx, actor, action, resource)
		if result.Allow {
			allowedResolvers = append(allowedResolvers, resource)
		}
	}

	return allowedResolvers, nil
}

func (a Authorizer) AuthorizeGlobal(ctx context.Context, actor *models.User, action GlobalPermission) decision.Decision {
	resolver := a.Global()

	if !action.Valid() {
		return decision.Error(ErrInvalidGlobalPermission)
	}

	switch action {
	case GlobalPermissionReadAllUsers:
		// Source: role - Admin
		if result := resolver.HasRoleAdmin(ctx, actor); result.Allow {
			return result
		}

	case GlobalPermissionWriteAllUsers:
		// Source: role - Admin
		if result := resolver.HasRoleAdmin(ctx, actor); result.Allow {
			return result
		}

	case GlobalPermissionReadAllProfiles:
		// Source: attribute - ProfilesArePublic
		if result := resolver.HasAttributeProfilesArePublic(ctx); result.Allow {
			return result
		}

	}
	return decision.False("unmatched")
}

// Authorizer
type Authorizer struct {
	Resolver
}

func (a Authorizer) Authorize(ctx context.Context, actor *models.User, permission string, resource any) decision.Decision {
	switch resource.(type) {
	case *models.Repository:
		perm, err := ParseRepositoryPermission(permission)
		if err == nil {
			return a.AuthorizeRepository(ctx, actor, perm, resource.(*models.Repository))
		}
	case *models.User:
		perm, err := ParseUserPermission(permission)
		if err == nil {
			return a.AuthorizeUser(ctx, actor, perm, resource.(*models.User))
		}
	}

	return decision.False("unmatched")
}

func NewAuthorizer(resolver Resolver) *Authorizer {
	return &Authorizer{Resolver: resolver}
}
