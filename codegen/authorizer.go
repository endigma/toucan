package codegen

import (
	. "github.com/dave/jennifer/jen"
	"github.com/samber/lo"
)

var ConcWaitGroup = func() *Statement { return Qual("github.com/sourcegraph/conc", "WaitGroup") }

func (gen *Generator) generateAuthorizerTypes(file *File) {
	file.Line().Type().Id("Authorizer").Interface(
		Id("Authorize").Params(
			Id("ctx").Qual("context", "Context"),
			Id("actor").Op("*").Qual(gen.Schema.Actor.Path, gen.Schema.Actor.Name),
			Id("permission").Id("Permission"),
			Id("resource").Any(),
		).Error(),
	)

	file.Line().Type().Id("AuthorizerFunc").Func().Params(
		Id("ctx").Qual("context", "Context"),
		Id("actor").Op("*").Qual(gen.Schema.Actor.Path, gen.Schema.Actor.Name),
		Id("permission").Id("Permission"),
		Id("resource").Any(),
	).Error()

	file.Line().Func().Params(
		Id("af").Id("AuthorizerFunc"),
	).Id("Authorize").Params(
		Id("ctx").Qual("context", "Context"),
		Id("actor").Op("*").Qual(gen.Schema.Actor.Path, gen.Schema.Actor.Name),
		Id("permission").Id("Permission"),
		Id("resource").Any(),
	).Error().Block(
		Return(Id("af").Call(
			Id("ctx"),
			Id("actor"),
			Id("permission"),
			Id("resource"),
		)),
	)

	file.Line().Type().Id("authorizerResult").Struct(
		Id("allow").Bool(),
		Id("source").String(),
		Id("error").Error(),
	)

	file.Line().Var().Id("authorizerData").Op("=").Map(Id("Permission")).Struct(
		Id("Attributes").Index().Id("Attribute"),
		Id("Roles").Index().Id("Role"),
	).Values(DictFunc(
		func(d Dict) {
			for _, resource := range gen.Schema.Resources {
				for _, permission := range resource.Permissions {
					d[Id("Permission"+pascal(resource.Name)+pascal(permission))] = Values(Dict{
						Id("Attributes"): Index().Id("Attribute").ValuesFunc(func(group *Group) {
							for _, attribute := range resource.Attributes {
								if lo.Contains(attribute.Permissions, permission) {
									group.Id("Attribute" + pascal(resource.Name) + pascal(attribute.Name))
								}
							}
						}),
						Id("Roles"): Index().Id("Role").ValuesFunc(func(group *Group) {
							for _, role := range resource.Roles {
								if lo.Contains(role.Permissions, permission) {
									group.Id("Role" + pascal(resource.Name) + pascal(role.Name))
								}
							}
						}),
					})
				}
			}
		},
	))
}

func (gen *Generator) generateAuthorizerRoot(group *Group) {
	group.Comment("Authorizer")
	group.Type().Id("authorizer").StructFunc(func(group *Group) {
		group.Id("resolver").Id("Resolver")
	})

	group.Func().Params(Id("a").Id("authorizer")).Id("Authorize").Params(
		Id("ctx").Qual("context", "Context"),
		Id("actor").Op("*").Qual(gen.Schema.Actor.Path, gen.Schema.Actor.Name),
		Id("permission").Id("Permission"),
		Id("resource").Any(),
	).Add(Error()).Block(
		Var().Id("cancel").Func().Params(),
		Id("ctx").Op(",").Id("cancel").Op("=").Qual("context", "WithCancel").Call(Id("ctx")),
		Defer().Id("cancel").Call(),
		Line(),
		Id("results").Op(":=").Make(Chan().Id("authorizerResult")),
		Var().Id("wg").Add(ConcWaitGroup()),
		Line(),
		List(Op("authorizerData"), Id("ok")).Op(":=").Id("authorizerData").Index(Id("permission")),
		If(Op("!").Id("ok")).Block(
			Return(Qual("fmt", "Errorf").Call(Lit("invalid permission %s"), Id("permission"))),
		),
		Line(),
		For(List(Id("_"), Id("attribute")).Op(":=").Range().Id("authorizerData").Dot("Attributes")).Block(
			Id("attribute").Op(":=").Id("attribute"),
			Id("wg").Dot("Go").Call(Func().Params().Block(
				List(Id("allow"), Id("err")).Op(":=").Id("a").Dot("resolver").Dot("HasAttribute").Call(
					Id("ctx"),
					Id("resource"),
					Id("attribute"),
				),
				Id("results").Op("<-").Id("authorizerResult").Values(Dict{
					Id("allow"):  Id("allow"),
					Id("source"): Qual("fmt", "Sprintf").Call(Lit("attribute %s"), Id("attribute")),
					Id("error"):  Err(),
				}),
			)),
		),
		If(Id("actor").Op("!=").Nil()).Block(
			For(List(Id("_"), Id("role")).Op(":=").Range().Id("authorizerData").Dot("Roles")).Block(
				Id("role").Op(":=").Id("role"),
				Id("wg").Dot("Go").Call(Func().Params().Block(
					List(Id("allow"), Id("err")).Op(":=").Id("a").Dot("resolver").Dot("HasRole").Call(
						Id("ctx"),
						Id("actor"),
						Id("resource"),
						Id("role"),
					),
					Id("results").Op("<-").Id("authorizerResult").Values(Dict{
						Id("allow"):  Id("allow"),
						Id("source"): Qual("fmt", "Sprintf").Call(Lit("role %s"), Id("role")),
						Id("error"):  Err(),
					}),
				)),
			),
		),
		Line(),
		Go().Func().Params().Block(
			Id("wg").Dot("Wait").Call(),
			Close(Id("results")),
		).Call(),
		Line(),
		Var().Id("denyReasons").Index().String(),
		For(List(Id("result")).Op(":=").Range().Id("results")).Block(
			If(Qual("errors", "Is").Call(Id("result").Dot("error"), Qual("context", "Canceled"))).Block(
				Continue(),
			),
			If(Id("result").Dot("error").Op("!=").Nil()).Block(
				Id("cancel").Call(),
				For(Range().Id("results")).Block(
					Comment("drain channel"),
				),
				Return(Id("result").Dot("error")),
			),
			If(Id("result").Dot("allow")).Block(
				Id("cancel").Call(),
				For(Range().Id("results")).Block(
					Comment("drain channel"),
				),
				Return(Qual("fmt", "Errorf").Call(
					Lit("authorize %s: %w: has %s"),
					Id("permission"),
					Id("Allow"),
					Id("result").Dot("source"),
				)),
			),
			Id("denyReasons").Op("=").Append(
				Id("denyReasons"),
				Qual("fmt", "Sprintf").Call(Lit("%s"), Id("result").Dot("source")),
			),
		),
		Line(),
		Return(Qual("fmt", "Errorf").Call(
			Lit("authorize %s: %w: missing %s"),
			Id("permission"),
			Id("Deny"),
			Qual("strings", "Join").Call(Id("denyReasons"), Lit(", ")),
		)),
	)

	group.Line()

	group.Func().Id("NewAuthorizer").Params(Id("resolver").Id("Resolver")).Id("Authorizer").Block(
		Return(Id("authorizer").Values(Dict{
			Id("resolver"): Id("resolver"),
		})),
	)
}
