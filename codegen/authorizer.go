package codegen

import (
	"fmt"

	. "github.com/dave/jennifer/jen"
	"github.com/endigma/toucan/schema"
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

	file.Line().Func().Params(
		Id("a").Id("authorizer"),
	).Id("checkAttributes").Params(
		Id("ctx").Qual("context", "Context"),
		Id("wg").Op("*").Qual("github.com/sourcegraph/conc", "WaitGroup"),
		Id("results").Chan().Op("<-").Id("authorizerResult"),
		Id("resource").Any(),
		Id("attributes").Index().Id("Attribute"),
	).Block(
		For(List(Id("_"), Id("attribute")).Op(":=").Range().Id("attributes")).Block(
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
	)

	file.Line().Func().Params(
		Id("a").Id("authorizer"),
	).Id("checkRoles").Params(
		Id("ctx").Qual("context", "Context"),
		Id("wg").Op("*").Qual("github.com/sourcegraph/conc", "WaitGroup"),
		Id("results").Chan().Op("<-").Id("authorizerResult"),
		Id("actor").Op("*").Qual(gen.Schema.Actor.Path, gen.Schema.Actor.Name),
		Id("resource").Any(),
		Id("roles").Index().Id("Role"),
	).Block(
		For(List(Id("_"), Id("role")).Op(":=").Range().Id("roles")).Block(
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
	)
}

func (gen *Generator) generateResourceAuthorizer(file *File, resource schema.ResourceSchema) {
	file.Line().Func().
		Params(
			Id("a").Id("authorizer"),
		).
		Id("authorize"+pascal(resource.Name)).
		Params(
			Id("ctx").Qual("context", "Context"),
			Id("actor").Op("*").Qual(gen.Schema.Actor.Path, gen.Schema.Actor.Name),
			Id("action").Id("Permission"),
			Id("resource").Id("any"),
		).Error().
		BlockFunc(func(group *Group) {
			group.Var().Id("cancel").Func().Params()
			group.Id("ctx").Op(",").Id("cancel").Op("=").Qual("context", "WithCancel").Call(Id("ctx"))
			group.Defer().Id("cancel").Call().Line()

			group.Id("results").Op(":=").Make(Chan().Id("authorizerResult"))

			group.Var().Id("wg").Add(ConcWaitGroup()).Line()

			if len(resource.Attributes) > 0 {
				group.Switch(Id("action")).BlockFunc(func(group *Group) {
					for _, permission := range resource.Permissions {
						sources := resource.GetAttributeSources(permission)
						if len(sources) == 0 {
							continue
						}

						generateAuthorizerCase(group, resource, permission, schema.PermissionTypeAttribute, sources)
					}
				})
			}

			group.Line()

			if len(resource.Roles) > 0 {
				group.If(Id("actor").Op("!=").Nil()).Block(
					Switch(Id("action")).BlockFunc(func(group *Group) {
						for _, permission := range resource.Permissions {
							sources := resource.GetRoleSources(permission)
							if len(sources) == 0 {
								continue
							}

							generateAuthorizerCase(group, resource, permission, schema.PermissionTypeRole, sources)
						}
					}),
				).Line()
			}

			group.Go().Func().Params().Block(
				Id("wg").Dot("Wait").Call(),
				Close(Id("results")),
			).Call().Line()

			group.Var().Id("denyReasons").Index().String()

			group.For(List(Id("result")).Op(":=").Range().Id("results")).Block(
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
						Lit(fmt.Sprintf("authorize %s: %%w: has %%s", resource.Name)),
						Id("Allow"),
						Id("result").Dot("source"),
					)),
				),
				Id("denyReasons").Op("=").Append(
					Id("denyReasons"),
					Qual("fmt", "Sprintf").Call(Lit("%s"), Id("result").Dot("source")),
				),
			)

			group.Return(Qual("fmt", "Errorf").Call(
				Lit(fmt.Sprintf("authorize %s: %%w: missing %%s", resource.Name)),
				Id("Deny"),
				Qual("strings", "Join").Call(Id("denyReasons"), Lit(", ")),
			))
		})

	file.Line()
}

func generateAuthorizerCase(
	group *Group,
	resource schema.ResourceSchema,
	perm string,
	sourceType schema.PermissionType,
	sources []schema.PermissionSource,
) {
	group.Case(Id("Permission" + pascal(resource.Name) + pascal(perm))).Block(
		Id("a").Do(func(s *Statement) {
			switch sourceType {
			case schema.PermissionTypeAttribute:
				s.Dot("checkAttributes")
			case schema.PermissionTypeRole:
				s.Dot("checkRoles")
			}
		}).Call(
			Id("ctx"),
			Op("&").Id("wg"),
			Id("results"),
			Do(func(s *Statement) {
				if sourceType == schema.PermissionTypeRole {
					s.Id("actor")
				}
			}),
			Id("resource"),
			Index().Id(pascal(string(sourceType))).BlockFunc(func(g *Group) {
				for _, source := range sources {
					g.Id(pascal(string(sourceType)) + pascal(resource.Name) + pascal(source.Name)).Op(",")
				}
			}),
		),
	)
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
	).Add(Error()).BlockFunc(func(group *Group) {
		group.Switch(Id("permission")).BlockFunc(func(group *Group) {
			for _, resource := range gen.Schema.Resources {
				group.CaseFunc(func(group *Group) {
					for i, permission := range resource.Permissions {
						group.
							Do(func(s *Statement) {
								if i > 0 {
									s.Line()
								}
							}).
							Id("Permission" + pascal(resource.Name) + pascal(permission))
					}
				}).Block(
					Return(Id("a").Dot("authorize"+pascal(resource.Name)).Call(
						Id("ctx"), Id("actor"), Id("permission"), Id("resource"),
					)),
				)
			}
		}).Line()

		group.Return(Qual("fmt", "Errorf").Call(
			Lit("invalid permission %s"),
			Id("permission"),
		))
	})

	group.Line()

	group.Func().Id("NewAuthorizer").Params(Id("resolver").Id("Resolver")).Id("Authorizer").Block(
		Return(Id("authorizer").Values(Dict{
			Id("resolver"): Id("resolver"),
		})),
	)
}
