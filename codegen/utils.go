package codegen

import (
	"github.com/dave/jennifer/jen"
	"github.com/endigma/toucan/spec"
	"github.com/iancoleman/strcase"
)

// String to lowerCamelCase
// technically a misnomer, but it's simpler to remember
func camel(s string) string {
	return strcase.ToLowerCamel(s)
}

// String to snake_case
func snake(s string) string {
	return strcase.ToSnake(s)
}

// String to PascalCase
func pascal(s string) string {
	return strcase.ToCamel(s)
}

func paramsForAuthorizer(actor spec.QualifierSpec, resource spec.ResourceSpec) []jen.Code {
	return []jen.Code{
		jen.Id("ctx").Qual("context", "Context"),
		jen.Id("actor").Op("*").Qual(actor.Path, actor.Name),
		jen.Id("action").Id(pascal(resource.Name) + "Permission"),
		jen.Id("resource").Op("*").Qual(resource.Model.Path, resource.Model.Name),
	}
}

func CallPermissionSource(source spec.PermissionSource) (string, *jen.Statement) {
	switch source.Type {
	case "role":
		return "HasRole", jen.Call(jen.Id("ctx"), jen.Id("actor"), jen.Id(source.Name), jen.Id("resource"))
	case "attribute":
		return "HasAttribute", jen.Call(jen.Id("ctx"), jen.Id(source.Name), jen.Id("resource"))
	}

	return "", jen.Null()

	// return Call(Id(source.CallName()), source.CallParams())
}
