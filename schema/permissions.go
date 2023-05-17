package schema

import (
	"github.com/iancoleman/strcase"
	"github.com/samber/lo"
)

type PermissionType string

const (
	PermissionTypeRole      PermissionType = "role"
	PermissionTypeAttribute PermissionType = "attribute"
)

type PermissionSource struct {
	Type PermissionType // role, attribute
	Name string
}

func GetRoleSources(permission string, roles []RoleSchema) []PermissionSource {
	sources := []PermissionSource{}

	for _, role := range roles {
		if lo.Contains(role.Permissions, permission) {
			sources = append(sources, PermissionSource{
				Type: PermissionTypeRole,
				Name: strcase.ToCamel(strcase.ToCamel(role.Name)),
			})
		}
	}

	return sources
}

func GetAttributeSources(permission string, attributes []AttributeSchema) []PermissionSource {
	sources := []PermissionSource{}

	for _, attr := range attributes {
		if lo.Contains(attr.Permissions, permission) {
			sources = append(sources, PermissionSource{
				Type: PermissionTypeAttribute,
				Name: strcase.ToCamel(strcase.ToCamel(attr.Name)),
			})
		}
	}

	return sources
}

func GetPermissionSources(permission string, attributes []AttributeSchema, roles []RoleSchema) []PermissionSource {
	return append(GetRoleSources(permission, roles), GetAttributeSources(permission, attributes)...)
}

func (resource ResourceSchema) GetRoleSources(permission string) []PermissionSource {
	return GetRoleSources(permission, resource.Roles)
}

func (resource ResourceSchema) GetAttributeSources(permission string) []PermissionSource {
	return GetAttributeSources(permission, resource.Attributes)
}

// GetPermissionSources returns all sources of a permission.
func (resource ResourceSchema) GetPermissionSources(permission string) []PermissionSource {
	return append(resource.GetRoleSources(permission), resource.GetAttributeSources(permission)...)
}
