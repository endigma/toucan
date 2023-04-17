package test

import (
	"context"
	"testing"

	"github.com/endigma/toucan/_examples/basic/gen/toucan"
	"github.com/endigma/toucan/_examples/basic/models"
	"github.com/endigma/toucan/_examples/basic/resolvers"
	"github.com/stretchr/testify/assert"
)

func TestPolicy(t *testing.T) {
	ctx := context.Background()

	authorizer := toucan.NewAuthorizer(resolvers.NewResolver())

	tom, jerry := models.NewUser("Tom"), models.NewUser("Jerry")

	// users := []*models.User{tom, jerry}

	google := models.NewRepository("Google", true)
	facebook := models.NewRepository("Facebook", false)

	// repos := []*models.Repository{google, facebook}

	assert.True(t,
		authorizer.Authorize(tom, "read", google),
		"Tom should be able to read Google",
	)

	assert.True(t,
		authorizer.Authorize(tom, "read", facebook),
		"Tom should be able to read Facebook",
	)

	assert.True(t, authorizer.Authorize(jerry, "read", google), "Jerry should be able to read Google")
	assert.False(t, authorizer.Authorize(jerry, "read", facebook), "Jerry should not be able to read Facebook")
}
