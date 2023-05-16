package codegen

import (
	"fmt"

	"github.com/dave/jennifer/jen"
	. "github.com/dave/jennifer/jen"
	"github.com/endigma/toucan/schema"
)

var (
	RuntimeDecision = func() *Statement { return Qual("github.com/endigma/toucan/decision", "Decision") }
	RuntimeTrue     = func() *Statement { return Qual("github.com/endigma/toucan/decision", "True") }
	RuntimeFalse    = func() *Statement { return Qual("github.com/endigma/toucan/decision", "False") }
	RuntimeError    = func() *Statement { return Qual("github.com/endigma/toucan/decision", "Error") }

	RuntimeCache    = func() *Statement { return Qual("github.com/endigma/toucan/cache", "Cache") }
	RuntimeCacheKey = func() *Statement { return Qual("github.com/endigma/toucan/cache", "CacheKey") }
	RuntimeQueryOr  = func() *Statement { return Qual("github.com/endigma/toucan/cache", "QueryOr") }
)

func (gen *Generator) generateResourceAuthorizer(file *File, resource schema.ResourceSchema) {
	file.Func().
		Params(
			Id("a").Id("Authorizer"),
		).
		Id("Authorize" + pascal(resource.Name)).
		Params(
			paramsForAuthorizer(gen.Schema.Actor, resource)...,
		).Add(RuntimeDecision()).
		BlockFunc(func(group *Group) {
			group.Id("resolver").Op(":=").Id("a").Dot(pascal(resource.Name)).Call().Line()

			group.If(Op("!").Id("action").Dot("Valid").Call()).Block(
				Return(
					RuntimeError().Call(Id(fmt.Sprintf("ErrInvalid%sPermission", pascal(resource.Name)))),
				),
			).Line()

			if len(resource.Attributes) > 0 {
				group.If(Id("resource").Op("!=").Nil()).BlockFunc(
					func(group *Group) {
						group.Switch(Id("action")).BlockFunc(func(group *Group) {
							for _, permission := range resource.Permissions {
								sources := resource.GetAttributeSources(permission)
								if len(sources) == 0 {
									continue
								}

								generateAuthorizerCase(group, resource.Name, permission, sources)
							}
						})
					},
				).Line()
			}

			if len(resource.Roles) > 0 {
				group.If(Id("resource").Op("!=").Nil().Op("&&").Id("actor").Op("!=").Nil()).Block(
					Switch(Id("action")).BlockFunc(func(group *Group) {
						for _, permission := range resource.Permissions {
							sources := resource.GetRoleSources(permission)
							if len(sources) == 0 {
								continue
							}

							generateAuthorizerCase(group, resource.Name, permission, sources)
						}
					}),
				).Line()
			}

			group.Return(RuntimeFalse().Call(Lit("unmatched")))
		})

	file.Line()
}

func (gen *Generator) generateResourceFilter(file *File, resource schema.ResourceSchema) {
	file.Func().
		Params(
			Id("a").Id("Authorizer"),
		).
		Id("Filter" + pascal(resource.Name)).
		Params(
			paramsForFilter(gen.Schema.Actor, resource)...,
		).Add(Params(Index().Op("*").Qual(resource.Model.Path, resource.Model.Name), Error())).
		BlockFunc(func(group *Group) {
			group.If(Op("!").Id("action").Dot("Valid").Call()).Block(
				Return(
					Nil(), Id(fmt.Sprintf("ErrInvalid%sPermission", pascal(resource.Name))),
				),
			).Line()

			group.Var().Id("allowedResolvers").Index().Op("*").Qual(resource.Model.Path, resource.Model.Name)
			group.For(jen.List(Id("_"), Id("resource")).Op(":=").Range().Id("resources")).
				BlockFunc(func(group *Group) {
					group.Id("result").Op(":=").Id("a").
						Dot("Authorize"+pascal(resource.Name)).Call(Id("ctx"), Id("actor"), Id("action"), Id("resource"))

					group.If(Id("result").Dot("Allow")).Block(Id("allowedResolvers").
						Op("=").Id("append").Call(Id("allowedResolvers"), Id("resource")))
				}).Line()
			group.Return(Id("allowedResolvers"), Nil())
		})
	file.Line()
}

func generateAuthorizerCase(group *Group, name string, perm string, sources []schema.PermissionSource) {
	group.Case(Id(pascal(name) + "Permission" + pascal(perm))).
		BlockFunc(
			func(group *Group) {
				for _, source := range sources {
					group.Commentf("Source: %s - %s", source.Type, source.Name)
					resolver, params := CallPermissionSource(source)
					group.If(
						Id("result").Op(":=").Add(RuntimeQueryOr()).Call(
							Id("ctx"),
							Add(RuntimeCacheKey()).Block(
								Id("ActorKey").Op(":").Id("actor").Dot("ToucanKey").Call().Op(","),
								Id("Resource").Op(":").Lit(name).Op(","),
								Id("ResourceKey").Op(":").Id("resource").Dot("ToucanKey").Call().Op(","),
								Id("SourceType").Op(":").Lit(source.Type).Op(","),
								Id("SourceName").Op(":").Lit(source.Name).Op(","),
							),
							Func().Params().Add(RuntimeDecision()).BlockFunc(func(group *jen.Group) {
								group.Return(Id("resolver").Dot(resolver).Add(params))
							}),
						).Op(";").Id("result").Dot("Allow").Block(Return(Id("result"))),
					)
					group.Line()
				}
			})
}

func (gen *Generator) generateGlobalAuthorizer(file *File) {
	file.Func().
		Params(
			Id("a").Id("Authorizer"),
		).
		Id("AuthorizeGlobal").
		Params(
			Id("ctx").Qual("context", "Context"),
			Id("actor").Op("*").Qual(gen.Schema.Actor.Path, gen.Schema.Actor.Name),
			Id("action").Id("GlobalPermission"),
		).Add(RuntimeDecision()).
		BlockFunc(func(group *Group) {
			group.Id("resolver").Op(":=").Id("a").Dot("Global").Call().Line()

			group.If(Op("!").Id("action").Dot("Valid").Call()).Block(
				Return(
					RuntimeError().Call(Id("ErrInvalidGlobalPermission")),
				),
			).Line()

			group.Switch(Id("action")).BlockFunc(func(group *Group) {
				for _, permission := range gen.Schema.Global.Permissions {
					sources := schema.GetPermissionSources(permission, gen.Schema.Global.Attributes, gen.Schema.Global.Roles)
					if len(sources) == 0 {
						continue
					}

					group.Case(Id("GlobalPermission" + pascal(permission))).
						BlockFunc(
							func(group *Group) {
								for _, source := range sources {
									resolver, params := CallGlobalSource(source)

									group.Commentf("Source: %s - %s", source.Type, source.Name)
									group.If(Id("result").Op(":=").Id("resolver").Dot(resolver).Add(params), Id("result").Dot("Allow")).Block(
										Return(Id("result")),
									)
									group.Line()
								}
							})
				}
			})

			group.Return(RuntimeFalse().Call(Lit("unmatched")))
		})

	file.Line()
}

func (gen *Generator) generateAuthorizerRoot(group *Group) {
	group.Comment("Authorizer")
	group.Type().Id("Authorizer").StructFunc(func(g *Group) {
		g.Id("Resolver")
	})

	group.Func().Params(Id("a").Id("Authorizer")).Id("Authorize").Params(
		Id("ctx").Qual("context", "Context"),
		Id("actor").Op("*").Qual(gen.Schema.Actor.Path, gen.Schema.Actor.Name),
		Id("permission").String(),
		Id("resource").Any(),
	).Add(RuntimeDecision()).BlockFunc(func(group *Group) {
		group.Switch(Id("resource").Assert(Type())).BlockFunc(func(group *Group) {
			for _, resource := range gen.Schema.Resources {
				group.Case(Op("*").Qual(resource.Model.Tuple())).Block(
					List(Id("perm"), Id("err")).Op(":=").Id("Parse"+pascal(resource.Name)+"Permission").Call(Id("permission")),
					If(Id("err").Op("==").Nil()).Block(
						Return(
							Id("a").
								Dot("Authorize"+pascal(resource.Name)).
								Call(
									Id("ctx"),
									Id("actor"),
									Id("perm"),
									Id("resource").
										Assert(
											Op("*").Qual(resource.Model.Path, resource.Model.Name),
										),
								),
						),
					),
				)
			}
		}).Line()

		group.Return(RuntimeFalse().Call(Lit("unmatched")))
	})

	group.Line()

	group.Func().Id("NewAuthorizer").Params(Id("resolver").Id("Resolver")).Op("*").Id("Authorizer").Block(
		Return(Op("&").Id("Authorizer").Values(Dict{
			Id("Resolver"): Id("resolver"),
		})),
	)
}
