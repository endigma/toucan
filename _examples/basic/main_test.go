package test

import (
	"context"
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

	authorizer := toucan.NewAuthorizer(resolvers.NewResolver())

	assert.True(t, authorizer.Repository().HasRoleEditor(ctx, jerry, google).Allow)

	// Define test cases
	testCases := []struct {
		name     string
		user     *models.User
		action   string
		repo     *models.Repository
		expected bool
	}{
		{
			name:     "Tom can read Facebook",
			user:     tom,
			action:   "read",
			repo:     facebook,
			expected: true,
		},
		{
			name:     "Tom can delete Facebook",
			user:     tom,
			action:   "delete",
			repo:     facebook,
			expected: true,
		},
		{
			name:     "Jerry can read Google",
			user:     jerry,
			action:   "read",
			repo:     google,
			expected: true,
		},
		{
			name:     "Graham can read Facebook",
			user:     graham,
			action:   "read",
			repo:     facebook,
			expected: true,
		},
		{
			name:     "Jerry cannot read Facebook",
			user:     jerry,
			action:   "read",
			repo:     facebook,
			expected: false,
		},
		{
			name:     "Graham cannot delete Facebook",
			user:     graham,
			action:   "delete",
			repo:     facebook,
			expected: false,
		},
		{
			name:     "Jerry cannot delete Google",
			user:     jerry,
			action:   "delete",
			repo:     google,
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := authorizer.Authorize(ctx, tc.user, tc.action, tc.repo)
			assert.Equal(t, result.Allow, tc.expected)
		})
	}
}

func TestFilter(t *testing.T) {
	ctx := context.Background()

	google := models.NewRepository("Google", true)
	facebook := models.NewRepository("Facebook", false)

	tom := models.NewUser("Tom", false, models.RepositoryRole{Role: "owner", Repo: facebook.ID})

	authorizer := toucan.NewAuthorizer(resolvers.NewResolver())

	allRepos := []*models.Repository{facebook, google}

	readRepos, err := authorizer.FilterRepository(ctx, tom, toucan.RepositoryPermissionRead, allRepos)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(readRepos), "there should be two repositories")

	writeRepos, err := authorizer.FilterRepository(ctx, tom, toucan.RepositoryPermissionPush, allRepos)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(writeRepos), "there should be one repository")
}
