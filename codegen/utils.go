package codegen

import (
	"github.com/dave/jennifer/jen"
	"github.com/endigma/toucan/schema"
	"github.com/iancoleman/strcase"
)

// String to lowerCamelCase
// technically a misnomer, but it's simpler to remember.
func camel(s string) string {
	return strcase.ToLowerCamel(s)
}

// String to snake_case.
func snake(s string) string {
	return strcase.ToSnake(s)
}

// String to PascalCase.
func pascal(s string) string {
	return strcase.ToCamel(s)
}

func paramsForAuthorizer(actor schema.Model, resource schema.ResourceSchema) func(*jen.Group) {
	return func(group *jen.Group) {
		group.Id("ctx").Qual("context", "Context")
		group.Id("actor").Op("*").Qual(actor.Path, actor.Name)
		group.Id("action").Id(pascal(resource.Name) + "Permission")

		if resource.Model != nil {
			group.Id("resource").Op("*").Qual(resource.Model.Path, resource.Model.Name)
		}
	}
}

func CallGlobalSource(source schema.PermissionSource) (string, *jen.Statement) {
	switch source.Type {
	case schema.PermissionTypeRole:
		return "HasRole" + pascal(source.Name), jen.Call(jen.Id("ctx"), jen.Id("actor"))
	case schema.PermissionTypeAttribute:
		return "HasAttribute" + pascal(source.Name), jen.Call(jen.Id("ctx"))
	}

	return "", jen.Null()
}

func CallPermissionSource(source schema.PermissionSource) (string, *jen.Statement) {
	switch source.Type {
	case schema.PermissionTypeRole:
		return "HasRole" + pascal(source.Name), jen.Call(jen.Id("ctx"), jen.Id("actor"), jen.Id("resource"))
	case schema.PermissionTypeAttribute:
		return "HasAttribute" + pascal(source.Name), jen.Call(jen.Id("ctx"), jen.Id("resource"))
	}

	return "", jen.Null()
}
