package codegen

import (
	. "github.com/dave/jennifer/jen"
)

func (gen *Generator) generateErrors(file *File) {
	file.Var().Id("Allow").Op("=").Qual("errors", "New").Call(Lit("allow"))
	file.Var().Id("Deny").Op("=").Qual("errors", "New").Call(Lit("deny"))
}
