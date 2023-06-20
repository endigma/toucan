// Code generated by toucan. DO NOT EDIT.
package toucan

import (
	"context"
	models "github.com/endigma/toucan/_examples/basic/models"
	cache "github.com/endigma/toucan/cache"
	decision "github.com/endigma/toucan/decision"
	conc "github.com/sourcegraph/conc"
	"strings"
)

type Authorizer interface {
	Authorize(ctx context.Context, actor *models.User, permission string, resourceType string, resource interface{}) decision.Decision
}

type AuthorizerFunc func(ctx context.Context, actor *models.User, permission string, resourceType string, resource interface{}) decision.Decision

func (af AuthorizerFunc) Authorize(ctx context.Context, actor *models.User, permission string, resourceType string, resource interface{}) decision.Decision {
	return af(ctx, actor, permission, resourceType, resource)
}

func (a authorizer) authorizeGlobal(ctx context.Context, actor *models.User, action GlobalPermission) decision.Decision {
	resolver := a.resolver.Global()

	if !action.Valid() {
		return decision.Error(ErrInvalidGlobalPermission)
	}

	var cancel func()
	ctx, cancel = context.WithCancel(ctx)
	defer cancel()

	results := make(chan decision.Decision)

	var wg conc.WaitGroup

	switch action {
	case GlobalPermissionReadAllProfiles:
		// Source: attribute - ProfilesArePublic
		wg.Go(func() {
			results <- cache.Query(ctx, cache.CacheKey{
				ActorKey:    "",
				Resource:    "global",
				ResourceKey: "",
				SourceType:  "attribute",
				SourceName:  "ProfilesArePublic",
			}, func() decision.Decision {
				return resolver.HasAttributeProfilesArePublic(ctx)
			})
		})

	}
	if actor != nil {
		switch action {
		case GlobalPermissionReadAllUsers:
			// Source: role - Admin
			wg.Go(func() {
				results <- cache.Query(ctx, cache.CacheKey{
					ActorKey:    a.resolver.CacheKey(actor),
					Resource:    "global",
					ResourceKey: "",
					SourceType:  "role",
					SourceName:  "Admin",
				}, func() decision.Decision {
					return resolver.HasRoleAdmin(ctx, actor)
				})
			})

		case GlobalPermissionWriteAllUsers:
			// Source: role - Admin
			wg.Go(func() {
				results <- cache.Query(ctx, cache.CacheKey{
					ActorKey:    a.resolver.CacheKey(actor),
					Resource:    "global",
					ResourceKey: "",
					SourceType:  "role",
					SourceName:  "Admin",
				}, func() decision.Decision {
					return resolver.HasRoleAdmin(ctx, actor)
				})
			})

		}
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var allowReason string
	var denyReasons []string
	for result := range results {
		if result.Reason == "" {
			result.Reason = "unspecified"
		}
		if result.Allow {
			cancel()
			allowReason = result.Reason
		} else {
			denyReasons = append(denyReasons, result.Reason)
		}
	}

	if allowReason != "" {
		return decision.True(allowReason)
	} else {
		result := decision.False(strings.Join(denyReasons, ", "))
		if result.Reason == "" {
			result.Reason = "unspecified"
		}
		return result
	}
}

func (a authorizer) authorizeRepository(ctx context.Context, actor *models.User, action RepositoryPermission, resource *models.Repository) decision.Decision {
	resolver := a.resolver.Repository()

	if !action.Valid() {
		return decision.Error(ErrInvalidRepositoryPermission)
	}

	if resource == nil {
		return decision.False("unmatched")
	}
	var cancel func()
	ctx, cancel = context.WithCancel(ctx)
	defer cancel()

	results := make(chan decision.Decision)

	var wg conc.WaitGroup

	switch action {
	case RepositoryPermissionRead:
		// Source: attribute - Public
		wg.Go(func() {
			results <- cache.Query(ctx, cache.CacheKey{
				ActorKey:    "",
				Resource:    "repository",
				ResourceKey: resolver.CacheKey(resource),
				SourceType:  "attribute",
				SourceName:  "Public",
			}, func() decision.Decision {
				return resolver.HasAttributePublic(ctx, resource)
			})
		})

	}
	if actor != nil {
		switch action {
		case RepositoryPermissionRead:
			// Source: role - Owner
			wg.Go(func() {
				results <- cache.Query(ctx, cache.CacheKey{
					ActorKey:    a.resolver.CacheKey(actor),
					Resource:    "repository",
					ResourceKey: resolver.CacheKey(resource),
					SourceType:  "role",
					SourceName:  "Owner",
				}, func() decision.Decision {
					return resolver.HasRoleOwner(ctx, actor, resource)
				})
			})

			// Source: role - Editor
			wg.Go(func() {
				results <- cache.Query(ctx, cache.CacheKey{
					ActorKey:    a.resolver.CacheKey(actor),
					Resource:    "repository",
					ResourceKey: resolver.CacheKey(resource),
					SourceType:  "role",
					SourceName:  "Editor",
				}, func() decision.Decision {
					return resolver.HasRoleEditor(ctx, actor, resource)
				})
			})

			// Source: role - Viewer
			wg.Go(func() {
				results <- cache.Query(ctx, cache.CacheKey{
					ActorKey:    a.resolver.CacheKey(actor),
					Resource:    "repository",
					ResourceKey: resolver.CacheKey(resource),
					SourceType:  "role",
					SourceName:  "Viewer",
				}, func() decision.Decision {
					return resolver.HasRoleViewer(ctx, actor, resource)
				})
			})

		case RepositoryPermissionPush:
			// Source: role - Owner
			wg.Go(func() {
				results <- cache.Query(ctx, cache.CacheKey{
					ActorKey:    a.resolver.CacheKey(actor),
					Resource:    "repository",
					ResourceKey: resolver.CacheKey(resource),
					SourceType:  "role",
					SourceName:  "Owner",
				}, func() decision.Decision {
					return resolver.HasRoleOwner(ctx, actor, resource)
				})
			})

			// Source: role - Editor
			wg.Go(func() {
				results <- cache.Query(ctx, cache.CacheKey{
					ActorKey:    a.resolver.CacheKey(actor),
					Resource:    "repository",
					ResourceKey: resolver.CacheKey(resource),
					SourceType:  "role",
					SourceName:  "Editor",
				}, func() decision.Decision {
					return resolver.HasRoleEditor(ctx, actor, resource)
				})
			})

		case RepositoryPermissionDelete:
			// Source: role - Owner
			wg.Go(func() {
				results <- cache.Query(ctx, cache.CacheKey{
					ActorKey:    a.resolver.CacheKey(actor),
					Resource:    "repository",
					ResourceKey: resolver.CacheKey(resource),
					SourceType:  "role",
					SourceName:  "Owner",
				}, func() decision.Decision {
					return resolver.HasRoleOwner(ctx, actor, resource)
				})
			})

		case RepositoryPermissionSnakeCase:
			// Source: role - Owner
			wg.Go(func() {
				results <- cache.Query(ctx, cache.CacheKey{
					ActorKey:    a.resolver.CacheKey(actor),
					Resource:    "repository",
					ResourceKey: resolver.CacheKey(resource),
					SourceType:  "role",
					SourceName:  "Owner",
				}, func() decision.Decision {
					return resolver.HasRoleOwner(ctx, actor, resource)
				})
			})

		}
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var allowReason string
	var denyReasons []string
	for result := range results {
		if result.Reason == "" {
			result.Reason = "unspecified"
		}
		if result.Allow {
			cancel()
			allowReason = result.Reason
		} else {
			denyReasons = append(denyReasons, result.Reason)
		}
	}

	if allowReason != "" {
		return decision.True(allowReason)
	} else {
		result := decision.False(strings.Join(denyReasons, ", "))
		if result.Reason == "" {
			result.Reason = "unspecified"
		}
		return result
	}
}

func (a authorizer) authorizeUser(ctx context.Context, actor *models.User, action UserPermission, resource *models.User) decision.Decision {
	resolver := a.resolver.User()

	if !action.Valid() {
		return decision.Error(ErrInvalidUserPermission)
	}

	if resource == nil {
		return decision.False("unmatched")
	}
	var cancel func()
	ctx, cancel = context.WithCancel(ctx)
	defer cancel()

	results := make(chan decision.Decision)

	var wg conc.WaitGroup

	if actor != nil {
		switch action {
		case UserPermissionRead:
			// Source: role - Admin
			wg.Go(func() {
				results <- cache.Query(ctx, cache.CacheKey{
					ActorKey:    a.resolver.CacheKey(actor),
					Resource:    "user",
					ResourceKey: resolver.CacheKey(resource),
					SourceType:  "role",
					SourceName:  "Admin",
				}, func() decision.Decision {
					return resolver.HasRoleAdmin(ctx, actor, resource)
				})
			})

			// Source: role - Self
			wg.Go(func() {
				results <- cache.Query(ctx, cache.CacheKey{
					ActorKey:    a.resolver.CacheKey(actor),
					Resource:    "user",
					ResourceKey: resolver.CacheKey(resource),
					SourceType:  "role",
					SourceName:  "Self",
				}, func() decision.Decision {
					return resolver.HasRoleSelf(ctx, actor, resource)
				})
			})

			// Source: role - Viewer
			wg.Go(func() {
				results <- cache.Query(ctx, cache.CacheKey{
					ActorKey:    a.resolver.CacheKey(actor),
					Resource:    "user",
					ResourceKey: resolver.CacheKey(resource),
					SourceType:  "role",
					SourceName:  "Viewer",
				}, func() decision.Decision {
					return resolver.HasRoleViewer(ctx, actor, resource)
				})
			})

		case UserPermissionWrite:
			// Source: role - Admin
			wg.Go(func() {
				results <- cache.Query(ctx, cache.CacheKey{
					ActorKey:    a.resolver.CacheKey(actor),
					Resource:    "user",
					ResourceKey: resolver.CacheKey(resource),
					SourceType:  "role",
					SourceName:  "Admin",
				}, func() decision.Decision {
					return resolver.HasRoleAdmin(ctx, actor, resource)
				})
			})

			// Source: role - Self
			wg.Go(func() {
				results <- cache.Query(ctx, cache.CacheKey{
					ActorKey:    a.resolver.CacheKey(actor),
					Resource:    "user",
					ResourceKey: resolver.CacheKey(resource),
					SourceType:  "role",
					SourceName:  "Self",
				}, func() decision.Decision {
					return resolver.HasRoleSelf(ctx, actor, resource)
				})
			})

		case UserPermissionDelete:
			// Source: role - Admin
			wg.Go(func() {
				results <- cache.Query(ctx, cache.CacheKey{
					ActorKey:    a.resolver.CacheKey(actor),
					Resource:    "user",
					ResourceKey: resolver.CacheKey(resource),
					SourceType:  "role",
					SourceName:  "Admin",
				}, func() decision.Decision {
					return resolver.HasRoleAdmin(ctx, actor, resource)
				})
			})

		}
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var allowReason string
	var denyReasons []string
	for result := range results {
		if result.Reason == "" {
			result.Reason = "unspecified"
		}
		if result.Allow {
			cancel()
			allowReason = result.Reason
		} else {
			denyReasons = append(denyReasons, result.Reason)
		}
	}

	if allowReason != "" {
		return decision.True(allowReason)
	} else {
		result := decision.False(strings.Join(denyReasons, ", "))
		if result.Reason == "" {
			result.Reason = "unspecified"
		}
		return result
	}
}

// Authorizer
type authorizer struct {
	resolver Resolver
}

func (a authorizer) Authorize(ctx context.Context, actor *models.User, permission string, resourceType string, resource any) decision.Decision {
	switch resourceType {
	case "global":
		perm, err := ParseGlobalPermission(permission)
		if err != nil {
			return decision.Error(err)
		}
		return a.authorizeGlobal(ctx, actor, perm)
	case "repository":
		perm, err := ParseRepositoryPermission(permission)
		resource, _ := resource.(*models.Repository)
		if err != nil {
			return decision.Error(err)
		}
		return a.authorizeRepository(ctx, actor, perm, resource)
	case "user":
		perm, err := ParseUserPermission(permission)
		resource, _ := resource.(*models.User)
		if err != nil {
			return decision.Error(err)
		}
		return a.authorizeUser(ctx, actor, perm, resource)
	}

	return decision.False("unmatched")
}

func NewAuthorizer(resolver Resolver) Authorizer {
	return authorizer{resolver: resolver}
}
