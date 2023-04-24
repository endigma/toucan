package schema_test

import (
	"testing"

	"github.com/endigma/toucan/schema"
	"github.com/stretchr/testify/assert"
)

//nolint:funlen
func TestSpec(t *testing.T) {
	t.Parallel()

	t.Run("valid", func(t *testing.T) {
		t.Parallel()

		schema := schema.Schema{
			Actor: schema.Model{"github.com/endigma/toucan/_examples/basic", "User"},
			Resources: []schema.ResourceSchema{
				{
					Name:        "Post",
					Model:       schema.Model{"github.com/endigma/toucan/_examples/basic", "Post"},
					Permissions: []string{"Read", "Write"},
					Roles: []schema.RoleSchema{
						{
							Name:        "Admin",
							Permissions: []string{"Read", "Write"},
						},
					},
				},
			},
		}

		assert.NoError(t, schema.Validate())
	})

	t.Run("invalid names", func(t *testing.T) {
		t.Parallel()

		schema := schema.Schema{
			Actor: schema.Model{"github.com/endigma/toucan/_examples/basic", "user"},
			Resources: []schema.ResourceSchema{
				{
					Name:        "Post",
					Model:       schema.Model{"github.com/endigma/toucan/_examples/basic", "9ost"},
					Permissions: []string{"read", "write"},
					Roles: []schema.RoleSchema{
						{
							Name:        "admin",
							Permissions: []string{"read", "write"},
						},
					},
				},
			},
		}

		assert.Error(t, schema.Validate())
	})

	t.Run("invalid paths", func(t *testing.T) {
		t.Parallel()

		schema := schema.Schema{
			Actor: schema.Model{"asdasd343//112.$asd", "User"},
			Resources: []schema.ResourceSchema{
				{
					Name:        "Post",
					Model:       schema.Model{"github.com/edd333", "Post"},
					Permissions: []string{"read", "write"},
					Roles: []schema.RoleSchema{
						{
							Name:        "admin",
							Permissions: []string{"read", "write"},
						},
					},
				},
			},
		}

		assert.Error(t, schema.Validate())
	})

	t.Run("duplicate resource names", func(t *testing.T) {
		t.Parallel()

		resource := schema.ResourceSchema{
			Name:        "Post",
			Model:       schema.Model{"github.com/endigma/toucan/_examples/basic", "Post"},
			Permissions: []string{"read", "write"},
			Roles: []schema.RoleSchema{
				{
					Name:        "admin",
					Permissions: []string{"read", "write"},
				},
			},
			Attributes: []schema.AttributeSchema{},
		}

		schema := schema.Schema{
			Actor: schema.Model{"github.com/endigma/toucan/_examples/basic", "User"},
			Resources: []schema.ResourceSchema{
				resource, resource,
			},
		}

		assert.Error(t, schema.Validate())
	})

	t.Run("duplicate role names", func(t *testing.T) {
		t.Parallel()

		role := schema.RoleSchema{
			Name:        "admin",
			Permissions: []string{"read", "write"},
		}

		schema := schema.Schema{
			Actor: schema.Model{"github.com/endigma/toucan/_examples/basic", "User"},
			Resources: []schema.ResourceSchema{
				{
					Name:        "Post",
					Model:       schema.Model{"github.com/endigma/toucan/_examples/basic", "Post"},
					Permissions: []string{"read", "write"},
					Roles:       []schema.RoleSchema{role, role},
				},
			},
		}

		assert.Error(t, schema.Validate())
	})

	t.Run("duplicate attribute names", func(t *testing.T) {
		t.Parallel()

		attr := schema.AttributeSchema{
			Name:        "Attribute",
			Permissions: []string{"read", "write"},
		}

		schema := schema.Schema{
			Actor: schema.Model{"github.com/endigma/toucan/_examples/basic", "User"},
			Resources: []schema.ResourceSchema{
				{
					Name:        "Post",
					Model:       schema.Model{"github.com/endigma/toucan/_examples/basic", "Post"},
					Permissions: []string{"read", "write"},
					Attributes:  []schema.AttributeSchema{attr, attr},
				},
			},
		}

		assert.Error(t, schema.Validate())
	})

	t.Run("duplicate permission names", func(t *testing.T) {
		t.Parallel()

		schema := schema.Schema{
			Actor: schema.Model{"github.com/endigma/toucan/_examples/basic", "User"},
			Resources: []schema.ResourceSchema{
				{
					Name:        "Post",
					Model:       schema.Model{"github.com/endigma/toucan/_examples/basic", "Post"},
					Permissions: []string{"read", "read"},
				},
			},
		}

		assert.Error(t, schema.Validate())
	})

	t.Run("unmatched permissions", func(t *testing.T) {
		t.Parallel()

		schema := schema.Schema{
			Actor: schema.Model{"github.com/endigma/toucan/_examples/basic", "User"},
			Resources: []schema.ResourceSchema{
				{
					Name:        "Post",
					Model:       schema.Model{"github.com/endigma/toucan/_examples/basic", "Post"},
					Permissions: []string{"read", "write", "read"},
					Roles: []schema.RoleSchema{
						{
							Name:        "admin",
							Permissions: []string{"eat_doritos"},
						},
					},
				},
			},
		}

		assert.Error(t, schema.Validate())
	})

	t.Run("unmatched permissions", func(t *testing.T) {
		t.Parallel()

		schema := schema.Schema{
			Actor: schema.Model{"github.com/endigma/toucan/_examples/basic", "User"},
			Resources: []schema.ResourceSchema{
				{
					Name:        "Post",
					Model:       schema.Model{"github.com/endigma/toucan/_examples/basic", "Post"},
					Permissions: []string{"read", "write", "read"},
					Attributes: []schema.AttributeSchema{
						{
							Name:        "public",
							Permissions: []string{"eat_doritos"},
						},
					},
				},
			},
		}

		assert.Error(t, schema.Validate())
	})

	t.Run("reserved names", func(t *testing.T) {
		t.Parallel()

		schema := schema.Schema{
			Actor: schema.Model{"github.com/endigma/toucan/_examples/basic", "User"},
			Resources: []schema.ResourceSchema{
				{
					Name:  "Authorizer",
					Model: schema.Model{"github.com/endigma/toucan/_examples/basic", "Post"},
				},
			},
		}

		assert.Error(t, schema.Validate())
	})

	t.Run("role with no perms", func(t *testing.T) {
		t.Parallel()

		schema := schema.Schema{
			Actor: schema.Model{"github.com/endigma/toucan/_examples/basic", "User"},
			Resources: []schema.ResourceSchema{
				{
					Name:  "Post",
					Model: schema.Model{"github.com/endigma/toucan/_examples/basic", "Post"},
					Roles: []schema.RoleSchema{
						{
							Name: "admin",
						},
					},
				},
			},
		}

		assert.Error(t, schema.Validate())
	})

	t.Run("attribute with no perms", func(t *testing.T) {
		t.Parallel()

		schema := schema.Schema{
			Actor: schema.Model{"github.com/endigma/toucan/_examples/basic", "User"},
			Resources: []schema.ResourceSchema{
				{
					Name:  "Post",
					Model: schema.Model{"github.com/endigma/toucan/_examples/basic", "Post"},
					Attributes: []schema.AttributeSchema{
						{
							Name: "admin",
						},
					},
				},
			},
		}

		assert.Error(t, schema.Validate())
	})
}
