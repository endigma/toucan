package schema

import (
	"github.com/iancoleman/strcase"
	"github.com/samber/lo"
)

type PermissionSource struct {
	Type string // role, attribute
	Name string
}

// GetPermissionSources returns all sources of a permission.
func (resource ResourceSchema) GetPermissionSources(permission string) []PermissionSource {
	sources := []PermissionSource{}

	for _, attr := range resource.Attributes {
		if lo.Contains(attr.Permissions, permission) {
			sources = append(sources, PermissionSource{
				Type: "attribute",
				Name: strcase.ToCamel(resource.Name + "Attribute" + strcase.ToCamel(attr.Name)),
			})
		}
	}

	for _, role := range resource.Roles {
		if lo.Contains(role.Permissions, permission) {
			sources = append(sources, PermissionSource{
				Type: "role",
				Name: strcase.ToCamel(resource.Name + "Role" + strcase.ToCamel(role.Name)),
			})
		}
	}

	return sources
}
