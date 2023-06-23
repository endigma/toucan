package codegen

import (
	"github.com/iancoleman/strcase"
)

// String to lowerCamelCase
// technically a misnomer, but it's simpler to remember.
func camel(s string) string {
	return strcase.ToLowerCamel(s)
}

// String to PascalCase.
func pascal(s string) string {
	return strcase.ToCamel(s)
}
