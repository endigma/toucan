package main_test

import (
	"context"
	"errors"
	"testing"

	"github.com/endigma/toucan/_examples/basic/gen/toucan"
	"github.com/endigma/toucan/_examples/basic/models"
	"github.com/endigma/toucan/_examples/basic/policy/resolvers"
	"github.com/stretchr/testify/assert"
)

func TestAuthorization(t *testing.T) {
	ctx := context.Background()

	google := models.NewRepository("Google", true)
	facebook := models.NewRepository("Facebook", false)

	tom, jerry, graham := models.NewUser("Tom", false, models.RepositoryRole{Role: "owner", Repo: facebook.ID}),
		models.NewUser("Jerry", false, models.RepositoryRole{Role: "editor", Repo: google.ID}),
		models.NewUser("Graham", false, models.RepositoryRole{Role: "viewer", Repo: facebook.ID})

	resolver := toucan.NewResolver(resolvers.NewResolver())
	authorizer := toucan.NewAuthorizer(resolver)

	isEditor, err := resolver.HasRole(ctx, jerry, google, toucan.RoleRepositoryEditor)
	assert.NoError(t, err)
	assert.True(t, isEditor)

	// Define test cases
	testCases := []struct {
		name     string
		user     *models.User
		action   toucan.Permission
		repo     *models.Repository
		expected bool
	}{
		{
			name:     "Tom can read Facebook",
			user:     tom,
			action:   toucan.PermissionRepositoryRead,
			repo:     facebook,
			expected: true,
		},
		{
			name:     "Tom can delete Facebook",
			user:     tom,
			action:   toucan.PermissionRepositoryDelete,
			repo:     facebook,
			expected: true,
		},
		{
			name:     "Jerry can read Google",
			user:     jerry,
			action:   toucan.PermissionRepositoryRead,
			repo:     google,
			expected: true,
		},
		{
			name:     "Graham can read Facebook",
			user:     graham,
			action:   toucan.PermissionRepositoryRead,
			repo:     facebook,
			expected: true,
		},
		{
			name:     "Jerry cannot read Facebook",
			user:     jerry,
			action:   toucan.PermissionRepositoryRead,
			repo:     facebook,
			expected: false,
		},
		{
			name:     "Graham cannot delete Facebook",
			user:     graham,
			action:   toucan.PermissionRepositoryDelete,
			repo:     facebook,
			expected: false,
		},
		{
			name:     "Jerry cannot delete Google",
			user:     jerry,
			action:   toucan.PermissionRepositoryDelete,
			repo:     google,
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := authorizer.Authorize(ctx, tc.user, tc.action, tc.repo)
			if tc.expected {
				assert.True(t, errors.Is(err, toucan.Allow))
				t.Log(err)
			} else {
				assert.True(t, errors.Is(err, toucan.Deny))
				t.Log(err)
			}
		})
	}
}
