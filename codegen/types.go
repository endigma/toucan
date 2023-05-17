package codegen

import (
	. "github.com/dave/jennifer/jen"
	"github.com/endigma/toucan/schema"
)

func (gen *Generator) generateResourceTypes(file *File, resource schema.ResourceSchema) {
	// Generate permissions enum
	if len(resource.Permissions) > 0 {
		enumGen := newEnumGenerator(resource.Name+"Permission", resource.Permissions, EnumGeneratorFeatures{})
		enumGen.Generate(file.Group)
	}

	file.Line()
}

type enumGenerator struct {
	enumName   string
	parserName string
	namesArray string
	namesMap   string
	errInvalid string
	errNil     string

	values []string

	features EnumGeneratorFeatures
}

type EnumGeneratorFeatures struct {
	MarshalerUnmarshaler bool
	ScannerValuer        bool
	StringHelper         bool
	ValidHelper          bool
}

func newEnumGenerator(name string, values []string, features EnumGeneratorFeatures) *enumGenerator {
	return &enumGenerator{
		enumName:   pascal(name),
		parserName: "Parse" + pascal(name),
		namesArray: camel(name) + "Names",
		namesMap:   camel(name) + "Map",
		errInvalid: "ErrInvalid" + pascal(name),
		errNil:     "ErrNil" + pascal(name),

		values: values,

		features: features,
	}
}

func (gen *enumGenerator) Generate(group *Group) {
	group.Type().Id(gen.enumName).String()

	// Constants
	group.Const().DefsFunc(func(group *Group) {
		for _, value := range gen.values {
			group.Id(gen.enumName + pascal(value)).Id(pascal(gen.enumName)).Op("=").Lit(snake(value))
		}
	}).Line()

	group.Var().Defs(
		Id(gen.errInvalid).
			Op("=").
			Qual("fmt", "Errorf").
			Call(
				Lit("not a valid "+gen.enumName+", try [%s]"),
				Qual("strings", "Join").
					Call(Id(gen.namesArray), Lit(", ")),
			),
		Id(gen.errNil).
			Op("=").
			Qual("errors", "New").
			Call(
				Lit("value is nil"),
			),
	).Line()

	// Definitions
	group.Var().Defs(
		Id(gen.namesMap).Op("=").Map(String()).Id(gen.enumName).Values(DictFunc(func(d Dict) {
			for _, value := range gen.values {
				d[Lit(snake(value))] = Id(gen.enumName + pascal(value))
			}
		})),
		Id(gen.namesArray).Op("=").Index().String().ValuesFunc(func(group *Group) {
			for _, value := range gen.values {
				group.String().Parens(Id(gen.enumName + pascal(value)))
			}
		}),
	).Line()

	// Parsing/validation
	gen.generateValidHelper(group)
	gen.generateParser(group)

	// Feature Flags
	if gen.features.MarshalerUnmarshaler {
		gen.generateMarshalText(group)
		gen.generateUnmarshalText(group)
	}

	if gen.features.ScannerValuer {
		gen.generateScanner(group)
		gen.generateValuer(group)
	}

	if gen.features.StringHelper {
		gen.generateToStringHelper(group)
	}
}

func (gen *enumGenerator) generateParser(group *Group) {
	group.Func().Id(gen.parserName).Params(Id("s").String()).Params(Id(gen.enumName), Error()).Block(
		If(
			List(Id("x"), Id("ok")).Op(":=").Id(gen.namesMap).Index(Id("s")),
			Id("ok"),
		).Block(Return(Id("x"), Nil())),
		Line().Comment("Try to parse from snake case"),
		If(
			List(Id("x"), Id("ok")).Op(":=").Id(gen.namesMap).
				Index(Qual("github.com/iancoleman/strcase", "ToSnake").Call(Id("s"))),
			Id("ok"),
		).Block(Return(Id("x"), Nil())),
		Line(),
		Return(
			Lit(""),
			Qual("fmt", "Errorf").Call(Lit("%s is %w"), Id("s"), Id("ErrInvalid"+gen.enumName)),
		),
	).Line()
}

func (gen *enumGenerator) generateMarshalText(group *Group) {
	group.Func().Params(Id("s").Id(gen.enumName)).Id("MarshalText").Params().Params(Index().Byte(), Error()).Block(
		Return(Index().Byte().Parens(String().Parens(Id("s"))), Nil()),
	).Line()
}

func (gen *enumGenerator) generateUnmarshalText(group *Group) {
	group.Func().Params(Id("s").Op("*").Id(gen.enumName)).
		Id("UnmarshalText").Params(Id("data").Index().Byte()).Error().
		Block(
			List(Id("x"), Id("err")).Op(":=").Id(gen.parserName).Call(String().Parens(Id("data"))),
			If(Id("err").Op("!=").Nil()).Block(
				Return(Id("err")),
			),
			Line(),
			Id("*s").Op("=").Id("x"),
			Return(Nil()),
		).
		Line()
}

func (gen *enumGenerator) generateScanner(group *Group) {
	group.Func().Params(Id("s").Op("*").Id(gen.enumName)).Id("Scan").
		Params(Id("value").Any()).Params(Id("err").Error()).
		Block(
			If(Id("value").Op("==").Nil()).Block(
				Op("*").Id("s").Op("=").Id(gen.enumName).Call(Lit("")),
				Return(Nil()),
			).Line(),

			Switch(Id("v").Op(":=").Id("value").Assert(Type())).Block(
				Case(Id("string")).Block(
					List(Op("*").Id("s"), Id("err")).Op("=").Id(gen.parserName).Call(Id("v")),
				),
				Case(Id("[]byte")).Block(
					List(Op("*").Id("s"), Id("err")).Op("=").Id(gen.parserName).Call(String().Parens(Id("v"))),
				),
				Case(Id(gen.enumName)).Block(
					Op("*").Id("s").Op("=").Id("v"),
				),
				Case(Op("*").Id(gen.enumName)).Block(
					If(Id("v").Op("==").Nil()).Block(
						Return(Id("")),
					),
					Id("*s").Op("=").Op("*").Id("v"),
				),
				Case(Op("*").String()).Block(
					If(Id("v").Op("==").Nil()).Block(
						Return(Id("")),
					),
					List(Op("*").Id("s"), Id("err")).Op("=").Id(gen.parserName).Call(Op("*").Id("v")),
				),
				Default().Block(
					Return(Qual("errors", "New").Call(Lit("invalid type for "+gen.enumName))),
				),
			),
			Line().Return(),
		).Line()
}

func (gen *enumGenerator) generateValuer(group *Group) {
	group.Func().
		Params(Id("s").
			Id(gen.enumName),
		).
		Id("Value").
		Params().
		Params(
			Qual("database/sql/driver", "Value"),
			Error(),
		).Block(
		Return(
			String().
				Parens(
					Id("s"),
				), Nil(),
		),
	)
}

func (gen *enumGenerator) generateValidHelper(group *Group) {
	group.Func().Params(Id("s").Id(gen.enumName)).Id("Valid").Params().Bool().Block(
		List(Id("_"), Id("err")).Op(":=").Id(gen.parserName).Call(String().Parens(Id("s"))),
		Return(Id("err").Op("==").Nil()),
	).Line()
}

func (gen *enumGenerator) generateToStringHelper(group *Group) {
	group.Func().Params(Id("s").Id(gen.enumName)).Id("String").Params().String().Block(
		Return(String().Parens(Id("s"))),
	).Line()
}
