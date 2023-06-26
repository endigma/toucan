package codegen

import (
	. "github.com/dave/jennifer/jen"
	"github.com/endigma/toucan/schema"
	"github.com/samber/lo"
)

func (gen *Generator) generateAttributeEnums(file *File) {
	roleNames := lo.FlatMap(gen.Schema.Resources, func(resource schema.ResourceSchema, _ int) []string {
		return lo.Map(
			resource.Roles,
			func(schema schema.RoleSchema, _ int) string { return resource.Name + "_" + schema.Name },
		)
	})

	roleValues := lo.FlatMap(gen.Schema.Resources, func(resource schema.ResourceSchema, _ int) []string {
		return lo.Map(
			resource.Roles,
			func(schema schema.RoleSchema, _ int) string { return resource.Name + "." + schema.Name },
		)
	})

	roleGen := newEnumGenerator("Role", roleNames, roleValues, EnumGeneratorFeatures{})
	roleGen.Generate(file.Group)

	attributeNames := lo.FlatMap(gen.Schema.Resources, func(resource schema.ResourceSchema, _ int) []string {
		return lo.Map(
			resource.Attributes,
			func(schema schema.AttributeSchema, _ int) string { return resource.Name + "_" + schema.Name },
		)
	})

	attributeValues := lo.FlatMap(gen.Schema.Resources, func(resource schema.ResourceSchema, _ int) []string {
		return lo.Map(
			resource.Attributes,
			func(schema schema.AttributeSchema, _ int) string { return resource.Name + "." + schema.Name },
		)
	})

	attributeGen := newEnumGenerator("Attribute", attributeNames, attributeValues, EnumGeneratorFeatures{})
	attributeGen.Generate(file.Group)
}

func (gen *Generator) generatePermissionEnum(file *File) {
	permissionNames := lo.FlatMap(gen.Schema.Resources, func(resource schema.ResourceSchema, _ int) []string {
		return lo.Map(
			resource.Permissions,
			func(name string, _ int) string {
				return resource.Name + "_" + name
			},
		)
	})

	permissionValues := lo.FlatMap(gen.Schema.Resources, func(resource schema.ResourceSchema, _ int) []string {
		return lo.Map(
			resource.Permissions,
			func(name string, _ int) string {
				return resource.Name + "." + name
			},
		)
	})

	permissionGen := newEnumGenerator("Permission", permissionNames, permissionValues, EnumGeneratorFeatures{})
	permissionGen.Generate(file.Group)
}

type enumGenerator struct {
	enumName   string
	parserName string
	namesArray string
	namesMap   string
	errInvalid string
	errNil     string

	names  []string
	values []string

	features EnumGeneratorFeatures
}

type EnumGeneratorFeatures struct {
	MarshalerUnmarshaler bool
	ScannerValuer        bool
	StringHelper         bool
	ValidHelper          bool
}

func newEnumGenerator(name string, names, values []string, features EnumGeneratorFeatures) *enumGenerator {
	return &enumGenerator{
		enumName:   pascal(name),
		parserName: "Parse" + pascal(name),
		namesArray: camel(name) + "Names",
		namesMap:   camel(name) + "Map",
		errInvalid: "ErrInvalid" + pascal(name),
		errNil:     "ErrNil" + pascal(name),

		names:  names,
		values: values,

		features: features,
	}
}

func (gen *enumGenerator) Generate(group *Group) {
	group.Type().Id(gen.enumName).String()

	// Constants
	group.Const().DefsFunc(func(group *Group) {
		for i, name := range gen.names {
			group.Id(gen.enumName + pascal(name)).Id(pascal(gen.enumName)).Op("=").Lit(gen.values[i])
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
			for i, name := range gen.names {
				d[Lit(gen.values[i])] = Id(gen.enumName + pascal(name))
			}
		})),
		Id(gen.namesArray).Op("=").Index().String().ValuesFunc(func(group *Group) {
			for _, name := range gen.names {
				group.String().Parens(Id(gen.enumName + pascal(name)))
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
