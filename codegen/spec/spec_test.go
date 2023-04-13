package spec_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/endigma/toucan/codegen/spec"
)

func TestSpec(t *testing.T) {
	t.Run("valid spec", func(t *testing.T) {
		s := spec.Spec{
			Actor: spec.QualifierSpec{
				Path: "github.com/endigma/toucan/_examples/basic",
				Name: "User",
			},
			Resources: []spec.ResourceSpec{
				{
					Name: "Post",
					Model: spec.QualifierSpec{
						Path: "github.com/endigma/toucan/_examples/basic",
						Name: "Post",
					},
					Permissions: []string{"Read", "Write"},
					Roles: []spec.RoleSpec{
						{
							Name:        "Admin",
							Permissions: []string{"Read", "Write"},
						},
					},
				},
			},
			Output: spec.OutputSpec{
				Path:    "output",
				Package: "output",
			},
		}

		assert.NoError(t, s.Validate())
	})

	t.Run("invalid spec (paths)", func(t *testing.T) {
		s := spec.Spec{
			Actor: spec.QualifierSpec{
				Path: "github.com/endigma/toucan/_examples/basic",
				Name: "user",
			},
			Resources: []spec.ResourceSpec{
				{
					Name: "Post",
					Model: spec.QualifierSpec{
						Path: "github.com/endigma/toucan/_examples/basic",
						Name: "9ost",
					},
					Permissions: []string{"read", "write"},
					Roles: []spec.RoleSpec{
						{
							Name:        "admin",
							Permissions: []string{"read", "write"},
						},
					},
				},
			},
			Output: spec.OutputSpec{
				Path:    "output",
				Package: "output",
			},
		}

		assert.Error(t, s.Validate())
	})

	t.Run("invalid spec (names)", func(t *testing.T) {
		s := spec.Spec{
			Actor: spec.QualifierSpec{
				Path: "asdasd343//112.$asd",
				Name: "User",
			},
			Resources: []spec.ResourceSpec{
				{
					Name: "Post",
					Model: spec.QualifierSpec{
						Path: "github.com/edd333",
						Name: "Post",
					},
					Permissions: []string{"read", "write"},
					Roles: []spec.RoleSpec{
						{
							Name:        "admin",
							Permissions: []string{"read", "write"},
						},
					},
				},
			},
			Output: spec.OutputSpec{
				Path:    "output",
				Package: "output",
			},
		}

		assert.Error(t, s.Validate())
	})

	t.Run("invalid spec (duplicate resource names)", func(t *testing.T) {
		resource := spec.ResourceSpec{
			Name: "Post",
			Model: spec.QualifierSpec{
				Path: "github.com/endigma/toucan/_examples/basic",
				Name: "Post",
			},
			Permissions: []string{"read", "write"},
			Roles: []spec.RoleSpec{
				{
					Name:        "admin",
					Permissions: []string{"read", "write"},
				},
			},
		}

		s := spec.Spec{
			Actor: spec.QualifierSpec{
				Path: "github.com/endigma/toucan/_examples/basic",
				Name: "User",
			},
			Resources: []spec.ResourceSpec{
				resource, resource,
			},
			Output: spec.OutputSpec{
				Path:    "output",
				Package: "output",
			},
		}

		assert.Error(t, s.Validate())
	})

	t.Run("invalid spec (duplicate role names)", func(t *testing.T) {
		role := spec.RoleSpec{
			Name:        "admin",
			Permissions: []string{"read", "write"},
		}

		s := spec.Spec{
			Actor: spec.QualifierSpec{
				Path: "github.com/endigma/toucan/_examples/basic",
				Name: "User",
			},
			Resources: []spec.ResourceSpec{
				{
					Name: "Post",
					Model: spec.QualifierSpec{
						Path: "github.com/endigma/toucan/_examples/basic",
						Name: "Post",
					},
					Permissions: []string{"read", "write"},
					Roles:       []spec.RoleSpec{role, role},
				},
			},
			Output: spec.OutputSpec{
				Path:    "output",
				Package: "output",
			},
		}

		assert.Error(t, s.Validate())
	})

	t.Run("invalid spec (duplicate attribute names)", func(t *testing.T) {
		attr := spec.AttributeSpec{
			Name:        "Attribute",
			Permissions: []string{"read", "write"},
		}

		s := spec.Spec{
			Actor: spec.QualifierSpec{
				Path: "github.com/endigma/toucan/_examples/basic",
				Name: "User",
			},
			Resources: []spec.ResourceSpec{
				{
					Name: "Post",
					Model: spec.QualifierSpec{
						Path: "github.com/endigma/toucan/_examples/basic",
						Name: "Post",
					},
					Permissions: []string{"read", "write"},
					Attributes:  []spec.AttributeSpec{attr, attr},
				},
			},
			Output: spec.OutputSpec{
				Path:    "output",
				Package: "output",
			},
		}

		assert.Error(t, s.Validate())
	})

	t.Run("invalid spec (duplicate permission names)", func(t *testing.T) {
		s := spec.Spec{
			Actor: spec.QualifierSpec{
				Path: "github.com/endigma/toucan/_examples/basic",
				Name: "User",
			},
			Resources: []spec.ResourceSpec{
				{
					Name: "Post",
					Model: spec.QualifierSpec{
						Path: "github.com/endigma/toucan/_examples/basic",
						Name: "Post",
					},
					Permissions: []string{"read", "read"},
				},
			},
			Output: spec.OutputSpec{
				Path:    "output",
				Package: "output",
			},
		}

		assert.Error(t, s.Validate())
	})

	t.Run("invalid spec: unmatched permissions", func(t *testing.T) {
		s := spec.Spec{
			Actor: spec.QualifierSpec{
				Path: "github.com/endigma/toucan/_examples/basic",
				Name: "User",
			},
			Resources: []spec.ResourceSpec{
				{
					Name: "Post",
					Model: spec.QualifierSpec{
						Path: "github.com/endigma/toucan/_examples/basic",
						Name: "Post",
					},
					Permissions: []string{"read", "write", "read"},
					Roles: []spec.RoleSpec{
						{
							Name:        "admin",
							Permissions: []string{"eat_doritos"},
						},
					},
				},
			},
			Output: spec.OutputSpec{
				Path:    "output",
				Package: "output",
			},
		}

		assert.Error(t, s.Validate())
	})

	t.Run("invalid spec: unmatched permissions", func(t *testing.T) {
		s := spec.Spec{
			Actor: spec.QualifierSpec{
				Path: "github.com/endigma/toucan/_examples/basic",
				Name: "User",
			},
			Resources: []spec.ResourceSpec{
				{
					Name: "Post",
					Model: spec.QualifierSpec{
						Path: "github.com/endigma/toucan/_examples/basic",
						Name: "Post",
					},
					Permissions: []string{"read", "write", "read"},
					Attributes: []spec.AttributeSpec{
						{
							Name:        "public",
							Permissions: []string{"eat_doritos"},
						},
					},
				},
			},
			Output: spec.OutputSpec{
				Path:    "output",
				Package: "output",
			},
		}

		assert.Error(t, s.Validate())
	})

	t.Run("invalid spec: reserved names", func(t *testing.T) {
		s := spec.Spec{
			Actor: spec.QualifierSpec{
				Path: "github.com/endigma/toucan/_examples/basic",
				Name: "User",
			},
			Resources: []spec.ResourceSpec{
				{
					Name: "Authorizer",
					Model: spec.QualifierSpec{
						Path: "github.com/endigma/toucan/_examples/basic",
						Name: "Post",
					},
				},
			},
			Output: spec.OutputSpec{
				Path:    "output",
				Package: "output",
			},
		}

		assert.Error(t, s.Validate())
	})

	t.Run("invalid spec: role with no perms", func(t *testing.T) {
		s := spec.Spec{
			Actor: spec.QualifierSpec{
				Path: "github.com/endigma/toucan/_examples/basic",
				Name: "User",
			},
			Resources: []spec.ResourceSpec{
				{
					Name: "Post",
					Model: spec.QualifierSpec{
						Path: "github.com/endigma/toucan/_examples/basic",
						Name: "Post",
					},
					Roles: []spec.RoleSpec{
						{
							Name: "admin",
						},
					},
				},
			},
			Output: spec.OutputSpec{
				Path:    "output",
				Package: "output",
			},
		}

		assert.Error(t, s.Validate())
	})

	t.Run("invalid spec: attribute with no perms", func(t *testing.T) {
		s := spec.Spec{
			Actor: spec.QualifierSpec{
				Path: "github.com/endigma/toucan/_examples/basic",
				Name: "User",
			},
			Resources: []spec.ResourceSpec{
				{
					Name: "Post",
					Model: spec.QualifierSpec{
						Path: "github.com/endigma/toucan/_examples/basic",
						Name: "Post",
					},
					Attributes: []spec.AttributeSpec{
						{
							Name: "admin",
						},
					},
				},
			},
			Output: spec.OutputSpec{
				Path:    "output",
				Package: "output",
			},
		}

		assert.Error(t, s.Validate())
	})
}
