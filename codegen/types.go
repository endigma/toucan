package codegen

import (
	. "github.com/dave/jennifer/jen"
	"github.com/endigma/toucan/spec"
	"github.com/samber/lo"
)

func generateResourceTypes(g *Group, resource spec.ResourceSpec) error {
	// Generate permissions enum
	if len(resource.Permissions) > 0 {
		err := generateStringEnum(g, resource.Name+"Permission", resource.Permissions)
		if err != nil {
			return err
		}
	}

	// Generate roles enum
	if len(resource.Roles) > 0 {
		err := generateStringEnum(g, resource.Name+"Role", lo.Map(resource.Roles, func(role spec.RoleSpec, _ int) string {
			return role.Name
		}))
		if err != nil {
			return err
		}
	}

	// Generate attributes enum
	if len(resource.Attributes) > 0 {
		err := generateStringEnum(g, resource.Name+"Attribute", lo.Map(resource.Attributes, func(attribute spec.AttributeSpec, _ int) string {
			return attribute.Name
		}))
		if err != nil {
			return err
		}
	}

	return nil
}

func generateStringEnum(g *Group, name string, values []string) error {
	enumName := pascal(name)
	parserName := "Parse" + enumName
	namesFunc := enumName + "Names"
	valuesFunc := enumName + "Values"
	namesArray := camel(name) + "Names"
	namesMap := camel(name) + "Map"
	errInvalid := "ErrInvalid" + enumName
	errNil := "ErrNil" + enumName

	g.Comment("Enum " + enumName)

	g.Type().Id(enumName).String()

	// ToString helper
	g.Func().Params(Id("s").Id(enumName)).Id("String").Params().String().Block(
		Return(String().Parens(Id("s"))),
	).Line()

	// Valid helper
	g.Func().Params(Id("s").Id(enumName)).Id("Valid").Params().Bool().Block(
		List(Id("_"), Id("err")).Op(":=").Id(parserName).Call(String().Parens(Id("s"))),
		Return(Id("err").Op("==").Nil()),
	).Line()

	// Invalid value error
	g.Var().
		Id(errInvalid).
		Op("=").
		Qual("fmt", "Errorf").
		Call(
			Lit("not a valid "+name+", try [%s]"),
			Qual("strings", "Join").
				Call(Id(namesArray), Lit(", ")),
		)

	// Null ptr error
	g.Var().
		Id(errNil).
		Op("=").
		Qual("errors", "New").
		Call(
			Lit("value is nil"),
		)

	// Constants
	g.Const().DefsFunc(func(g *Group) {
		for _, value := range values {
			g.Id(enumName + pascal(value)).Id(pascal(name)).Op("=").Lit(snake(value))
		}
	})

	// Names array and map
	g.Var().Id(namesArray).Op("=").Index().String().ValuesFunc(func(g *Group) {
		for _, value := range values {
			g.String().Parens(Id(enumName + pascal(value)))
		}
	})

	g.Var().Id(namesMap).Op("=").Map(String()).Id(enumName).Values(DictFunc(func(d Dict) {
		for _, value := range values {
			d[Lit(snake(value))] = Id(enumName + pascal(value))
		}
	}))

	// Names and Values functions
	g.Func().Id(namesFunc).Params().Index().String().Block(
		Id("tmp").Op(":=").Make(Index().String(), Len(Id(namesArray))),
		Copy(Id("tmp"), Id(namesArray)),
		Return(Id("tmp")),
	).Line()

	g.Func().Id(valuesFunc).Params().Index().Id(enumName).Block(
		Return(Index().Id(enumName).ValuesFunc(func(g *Group) {
			for _, value := range values {
				g.Id(enumName + pascal(value))
			}
		})),
	).Line()

	// Parsing
	g.Func().Id(parserName).Params(Id("s").String()).Params(Id(enumName), Error()).Block(
		If(
			List(Id("x"), Id("ok")).Op(":=").Id(namesMap).Index(Id("s")),
			Id("ok"),
		).Block(Return(Id("x"), Nil())),
		Line().Comment("Try to parse from snake case"),
		If(
			List(Id("x"), Id("ok")).Op(":=").Id(namesMap).Index(Qual("github.com/iancoleman/strcase", "ToSnake").Call(Id("s"))),
			Id("ok"),
		).Block(Return(Id("x"), Nil())),
		Line(),
		Return(
			Id(enumName).Call(Lit("")),
			Qual("fmt", "Errorf").Call(Lit("%s is %w"), Id("s"), Id("ErrInvalid"+enumName)),
		),
	).Line()

	g.Func().Id("MustParse"+enumName).Params(Id("s").String()).Id(enumName).Block(
		List(Id("x"), Id("err")).Op(":=").Id(parserName).Call(Id("s")),
		If(Id("err").Op("!=").Nil()).Block(
			Panic(Id("err")),
		),
		Line(),
		Return(Id("x")),
	).Line()

	// Text marshalling
	g.Func().Params(Id("s").Id(enumName)).Id("MarshalText").Params().Params(Index().Byte(), Error()).Block(
		Return(Index().Byte().Parens(String().Parens(Id("s"))), Nil()),
	).Line()

	g.Func().Params(Id("s").Op("*").Id(enumName)).Id("UnmarshalText").Params(Id("data").Index().Byte()).Error().Block(
		List(Id("x"), Id("err")).Op(":=").Id(parserName).Call(String().Parens(Id("data"))),
		If(Id("err").Op("!=").Nil()).Block(
			Return(Id("err")),
		),
		Line(),
		Id("*s").Op("=").Id("x"),
		Return(Nil()),
	).Line()

	// Scanner interface
	g.Func().Params(Id("s").Op("*").Id(enumName)).Id("Scan").Params(Id("value").Any()).Params(Id("err").Error()).Block(
		If(Id("value").Op("==").Nil()).Block(
			Op("*").Id("s").Op("=").Id(enumName).Call(Lit("")),
			Return(Nil()),
		).Line(),

		Switch(Id("v").Op(":=").Id("value").Assert(Type())).Block(
			Case(Id("string")).Block(
				List(Op("*").Id("s"), Id("err")).Op("=").Id(parserName).Call(Id("v")),
			),
			Case(Id("[]byte")).Block(
				List(Op("*").Id("s"), Id("err")).Op("=").Id(parserName).Call(String().Parens(Id("v"))),
			),
			Case(Id(enumName)).Block(
				Op("*").Id("s").Op("=").Id("v"),
			),
			Case(Op("*").Id(enumName)).Block(
				If(Id("v").Op("==").Nil()).Block(
					Return(Id("")),
				),
				Id("*s").Op("=").Op("*").Id("v"),
			),
			Case(Op("*").String()).Block(
				If(Id("v").Op("==").Nil()).Block(
					Return(Id("")),
				),
				List(Op("*").Id("s"), Id("err")).Op("=").Id(parserName).Call(Op("*").Id("v")),
			),
			Default().Block(
				Return(Qual("errors", "New").Call(Lit("invalid type for "+enumName))),
			),
		),
		Line().Return(),
	).Line()

	// Valuer interface
	g.Func().Params(Id("s").Id(enumName)).Id("Value").Params().Params(Qual("database/sql/driver", "Value"), Error()).Block(
		Return(String().Parens(Id("s")), Nil()),
	)

	return nil
}
