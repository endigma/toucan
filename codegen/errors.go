package codegen

import (
	. "github.com/dave/jennifer/jen"
)

func (gen *Generator) generateErrors(file *File) {
	file.Var().Id("Allow").Op("=").Qual("errors", "New").Call(Lit("allow")).Line()
	file.Var().Id("Deny").Op("=").Qual("errors", "New").Call(Lit("deny")).Line()

	file.Comment("AuthorizerResult takes an error, which may be Allow, Deny, or neither.")
	file.Comment("If err wraps Allow or Deny, it returns whether err was Allow and an nil error.")
	file.Comment("Otherwise, it returns false and err.")
	file.Func().Id("AuthorizerResult").Params(Err().Error()).Params(Bool(), Error()).Block(
		If(Qual("errors", "Is").Call(Err(), Id("Allow"))).Block(
			Return(True(), Nil()),
		),
		If(Qual("errors", "Is").Call(Err(), Id("Deny"))).Block(
			Return(False(), Nil()),
		),
		Return(False(), Err()),
	)
}
