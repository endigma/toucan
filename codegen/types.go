package codegen

import (
	. "github.com/dave/jennifer/jen"
	"github.com/endigma/toucan/schema"
)

func generateResourceTypes(group *Group, resource schema.ResourceSchema) error {
	// Generate permissions enum
	if len(resource.Permissions) > 0 {
		err := generateStringEnum(group, resource.Name+"Permission", resource.Permissions)
		if err != nil {
			return err
		}
	}

	// Generate roles enum
	// if len(resource.Roles) > 0 {
	// 	err := generateStringEnum(
	// 		group,
	// 		resource.Name+"Role",
	// 		lo.Map(resource.Roles, func(role schema.RoleSchema, _ int) string {
	// 			return role.Name
	// 		}))
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// Generate attributes enum
	// if len(resource.Attributes) > 0 {
	// 	err := generateStringEnum(group,
	// 		resource.Name+"Attribute",
	// 		lo.Map(resource.Attributes, func(attribute schema.AttributeSchema, _ int) string {
	// 			return attribute.Name
	// 		}))
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

func generateStringEnum(group *Group, name string, values []string) error {
	enumName := pascal(name)
	parserName := "Parse" + enumName
	// namesFunc := enumName + "Names"
	// valuesFunc := enumName + "Values"
	namesArray := camel(name) + "Names"
	namesMap := camel(name) + "Map"
	errInvalid := "ErrInvalid" + enumName
	errNil := "ErrNil" + enumName

	group.Comment("Enum " + enumName)

	group.Type().Id(enumName).String()

	// Constants
	group.Const().DefsFunc(func(group *Group) {
		for _, value := range values {
			group.Id(enumName + pascal(value)).Id(pascal(name)).Op("=").Lit(snake(value))
		}
	})

	// ToString helper
	generateToStringHelper(group, enumName)

	// Valid helper
	generateValidHelper(group, enumName, parserName)

	group.Var().Defs(
		Id(errInvalid).
			Op("=").
			Qual("fmt", "Errorf").
			Call(
				Lit("not a valid "+name+", try [%s]"),
				Qual("strings", "Join").
					Call(Id(namesArray), Lit(", ")),
			),
		Id(errNil).
			Op("=").
			Qual("errors", "New").
			Call(
				Lit("value is nil"),
			),
	)

	group.Line()

	group.Var().Defs(
		Id(namesMap).Op("=").Map(String()).Id(enumName).Values(DictFunc(func(d Dict) {
			for _, value := range values {
				d[Lit(snake(value))] = Id(enumName + pascal(value))
			}
		})),
		Id(namesArray).Op("=").Index().String().ValuesFunc(func(group *Group) {
			for _, value := range values {
				group.String().Parens(Id(enumName + pascal(value)))
			}
		}),
	)

	// Invalid value error
	// group.Var().

	// Null ptr error
	// group.Var().

	// Names array and map
	// group.Var().
	// Names and Values functions
	// group.Func().Id(namesFunc).Params().Index().String().Block(
	// 	Id("tmp").Op(":=").Make(Index().String(), Len(Id(namesArray))),
	// 	Copy(Id("tmp"), Id(namesArray)),
	// 	Return(Id("tmp")),
	// ).Line()

	// group.Func().Id(valuesFunc).Params().Index().Id(enumName).Block(
	// 	Return(Index().Id(enumName).ValuesFunc(func(group *Group) {
	// 		for _, value := range values {
	// 			group.Id(enumName + pascal(value))
	// 		}
	// 	})),
	// ).Line()

	// Parsing
	group.Func().Id(parserName).Params(Id("s").String()).Params(Id(enumName), Error()).Block(
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

	// group.Func().Id("MustParse"+enumName).Params(Id("s").String()).Id(enumName).Block(
	// 	List(Id("x"), Id("err")).Op(":=").Id(parserName).Call(Id("s")),
	// 	If(Id("err").Op("!=").Nil()).Block(
	// 		Panic(Id("err")),
	// 	),
	// 	Line(),
	// 	Return(Id("x")),
	// ).Line()

	// Text marshalling
	// group.Func().Params(Id("s").Id(enumName)).Id("MarshalText").Params().Params(Index().Byte(), Error()).Block(
	// 	Return(Index().Byte().Parens(String().Parens(Id("s"))), Nil()),
	// ).Line()

	// group.Func().Params(Id("s").Op("*").Id(enumName)).Id("UnmarshalText").Params(Id("data").Index().Byte()).Error().Block(
	// 	List(Id("x"), Id("err")).Op(":=").Id(parserName).Call(String().Parens(Id("data"))),
	// 	If(Id("err").Op("!=").Nil()).Block(
	// 		Return(Id("err")),
	// 	),
	// 	Line(),
	// 	Id("*s").Op("=").Id("x"),
	// 	Return(Nil()),
	// ).Line()

	// Scanner interface
	// group.Func().Params(Id("s").Op("*").Id(enumName)).Id("Scan").Params(Id("value").Any()).Params(Id("err").Error()).Block(
	// 	If(Id("value").Op("==").Nil()).Block(
	// 		Op("*").Id("s").Op("=").Id(enumName).Call(Lit("")),
	// 		Return(Nil()),
	// 	).Line(),

	// 	Switch(Id("v").Op(":=").Id("value").Assert(Type())).Block(
	// 		Case(Id("string")).Block(
	// 			List(Op("*").Id("s"), Id("err")).Op("=").Id(parserName).Call(Id("v")),
	// 		),
	// 		Case(Id("[]byte")).Block(
	// 			List(Op("*").Id("s"), Id("err")).Op("=").Id(parserName).Call(String().Parens(Id("v"))),
	// 		),
	// 		Case(Id(enumName)).Block(
	// 			Op("*").Id("s").Op("=").Id("v"),
	// 		),
	// 		Case(Op("*").Id(enumName)).Block(
	// 			If(Id("v").Op("==").Nil()).Block(
	// 				Return(Id("")),
	// 			),
	// 			Id("*s").Op("=").Op("*").Id("v"),
	// 		),
	// 		Case(Op("*").String()).Block(
	// 			If(Id("v").Op("==").Nil()).Block(
	// 				Return(Id("")),
	// 			),
	// 			List(Op("*").Id("s"), Id("err")).Op("=").Id(parserName).Call(Op("*").Id("v")),
	// 		),
	// 		Default().Block(
	// 			Return(Qual("errors", "New").Call(Lit("invalid type for "+enumName))),
	// 		),
	// 	),
	// 	Line().Return(),
	// ).Line()

	// Valuer interface
	// group.Func().
	// 	Params(Id("s").
	// 		Id(enumName),
	// 	).
	// 	Id("Value").
	// 	Params().
	// 	Params(
	// 		Qual("database/sql/driver", "Value"),
	// 		Error(),
	// 	).Block(
	// 	Return(
	// 		String().
	// 			Parens(
	// 				Id("s"),
	// 			), Nil(),
	// 	),
	// )

	return nil
}

func generateValidHelper(group *Group, enumName string, parserName string) {
	group.Func().Params(Id("s").Id(enumName)).Id("Valid").Params().Bool().Block(
		List(Id("_"), Id("err")).Op(":=").Id(parserName).Call(String().Parens(Id("s"))),
		Return(Id("err").Op("==").Nil()),
	).Line()
}

func generateToStringHelper(group *Group, enumName string) {
	group.Func().Params(Id("s").Id(enumName)).Id("String").Params().String().Block(
		Return(String().Parens(Id("s"))),
	).Line()
}
