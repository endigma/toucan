// Code generated by toucan. DO NOT EDIT.
package toucan

import (
	"context"
	models "github.com/endigma/toucan/_examples/basic/models"
	decision "github.com/endigma/toucan/decision"
	conc "github.com/sourcegraph/conc"
	"strings"
)

type Authorizer interface {
	Authorize(ctx context.Context, actor *models.User, permission Permission, resource any) decision.Decision
}

type AuthorizerFunc func(ctx context.Context, actor *models.User, permission Permission, resource any) decision.Decision

func (af AuthorizerFunc) Authorize(ctx context.Context, actor *models.User, permission Permission, resource any) decision.Decision {
	return af(ctx, actor, permission, resource)
}

func (a authorizer) authorizeGlobal(ctx context.Context, actor *models.User, action Permission) decision.Decision {
	var cancel func()
	ctx, cancel = context.WithCancel(ctx)
	defer cancel()

	results := make(chan decision.Decision)

	var wg conc.WaitGroup

	switch action {

	case PermissionGlobalReadAllProfiles:
		// Source: attribute - profiles_are_public
		wg.Go(func() {
			results <- a.resolver.HasAttribute(ctx, nil, AttributeGlobalProfilesArePublic)
		})
	}

	if actor != nil {
		switch action {
		case PermissionGlobalReadAllUsers:
			// Source: role - admin
			wg.Go(func() {
				results <- a.resolver.HasRole(ctx, actor, nil, RoleGlobalAdmin)
			})

		case PermissionGlobalWriteAllUsers:
			// Source: role - admin
			wg.Go(func() {
				results <- a.resolver.HasRole(ctx, actor, nil, RoleGlobalAdmin)
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

func (a authorizer) authorizeRepository(ctx context.Context, actor *models.User, action Permission, resource *models.Repository) decision.Decision {
	if resource == nil {
		return decision.False("unmatched")
	}
	var cancel func()
	ctx, cancel = context.WithCancel(ctx)
	defer cancel()

	results := make(chan decision.Decision)

	var wg conc.WaitGroup

	switch action {
	case PermissionRepositoryRead:
		// Source: attribute - public
		wg.Go(func() {
			results <- a.resolver.HasAttribute(ctx, resource, AttributeRepositoryPublic)
		})
	}

	if actor != nil {
		switch action {
		case PermissionRepositoryRead:
			// Source: role - owner
			wg.Go(func() {
				results <- a.resolver.HasRole(ctx, actor, resource, RoleRepositoryOwner)
			})

			// Source: role - editor
			wg.Go(func() {
				results <- a.resolver.HasRole(ctx, actor, resource, RoleRepositoryEditor)
			})

			// Source: role - viewer
			wg.Go(func() {
				results <- a.resolver.HasRole(ctx, actor, resource, RoleRepositoryViewer)
			})

		case PermissionRepositoryPush:
			// Source: role - owner
			wg.Go(func() {
				results <- a.resolver.HasRole(ctx, actor, resource, RoleRepositoryOwner)
			})

			// Source: role - editor
			wg.Go(func() {
				results <- a.resolver.HasRole(ctx, actor, resource, RoleRepositoryEditor)
			})

		case PermissionRepositoryDelete:
			// Source: role - owner
			wg.Go(func() {
				results <- a.resolver.HasRole(ctx, actor, resource, RoleRepositoryOwner)
			})

		case PermissionRepositorySnakeCase:
			// Source: role - owner
			wg.Go(func() {
				results <- a.resolver.HasRole(ctx, actor, resource, RoleRepositoryOwner)
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

func (a authorizer) authorizeUser(ctx context.Context, actor *models.User, action Permission, resource *models.User) decision.Decision {
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
		case PermissionUserRead:
			// Source: role - admin
			wg.Go(func() {
				results <- a.resolver.HasRole(ctx, actor, resource, RoleUserAdmin)
			})

			// Source: role - self
			wg.Go(func() {
				results <- a.resolver.HasRole(ctx, actor, resource, RoleUserSelf)
			})

			// Source: role - viewer
			wg.Go(func() {
				results <- a.resolver.HasRole(ctx, actor, resource, RoleUserViewer)
			})

		case PermissionUserWrite:
			// Source: role - admin
			wg.Go(func() {
				results <- a.resolver.HasRole(ctx, actor, resource, RoleUserAdmin)
			})

			// Source: role - self
			wg.Go(func() {
				results <- a.resolver.HasRole(ctx, actor, resource, RoleUserSelf)
			})

		case PermissionUserDelete:
			// Source: role - admin
			wg.Go(func() {
				results <- a.resolver.HasRole(ctx, actor, resource, RoleUserAdmin)
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

func (a authorizer) Authorize(ctx context.Context, actor *models.User, permission Permission, resource any) decision.Decision {
	switch permission {
	case PermissionGlobalReadAllUsers,
		PermissionGlobalWriteAllUsers,
		PermissionGlobalReadAllProfiles:
		return a.authorizeGlobal(ctx, actor, permission)
	case PermissionRepositoryRead,
		PermissionRepositoryPush,
		PermissionRepositoryDelete,
		PermissionRepositorySnakeCase:
		resource, _ := resource.(*models.Repository)
		return a.authorizeRepository(ctx, actor, permission, resource)
	case PermissionUserRead,
		PermissionUserWrite,
		PermissionUserDelete:
		resource, _ := resource.(*models.User)
		return a.authorizeUser(ctx, actor, permission, resource)
	}

	return decision.False("unmatched")
}

func NewAuthorizer(resolver Resolver) Authorizer {
	return authorizer{resolver: resolver}
}
