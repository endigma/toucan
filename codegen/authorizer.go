package codegen

import (
	. "github.com/dave/jennifer/jen"
	"github.com/endigma/toucan/schema"
	"github.com/samber/lo"
)

func generateResourceAuthorizer(group *Group, actor schema.Model, resource schema.ResourceSchema) error {
	group.Comment("authorizer for resource `" + resource.Name + "`").Line()

	group.Func().
		Params(
			Id("a").Id("Authorizer"),
		).
		Id("Authorize" + pascal(resource.Name)).
		Params(
			paramsForAuthorizer(actor, resource)...,
		).Bool().
		BlockFunc(func(group *Group) {
			group.If(Op("!").Id("action").Dot("Valid").Call()).Block(
				Return(False()),
			)

			group.Switch(Id("action")).BlockFunc(func(group *Group) {
				for _, perm := range resource.Permissions {
					sources := resource.GetPermissionSources(perm)

					group.Case(Id(pascal(resource.Name) + "Permission" + pascal(perm))).Block(
						Return(lo.Reduce(sources, func(statement *Statement, source schema.PermissionSource, n int) *Statement {
							resolver, params := CallPermissionSource(source)
							call := Id("a").Dot("resolver").Dot(pascal(resource.Name)).Call().Dot(resolver).Add(params)

							if n == 0 {
								return statement.Add(call)
							} else {
								return statement.Op("||").Line().Add(call)
							}
						}, &Statement{})),
					)
				}
				group.Default().Block(Return(False()))
			})
		})

	return nil
}

func generateGlobalAuthorizer(group *Group, actor schema.Model, resources []schema.ResourceSchema) {
	group.Comment("Global authorizer")
	group.Type().Id("Authorizer").StructFunc(func(g *Group) {
		g.Id("resolver").Id("Resolver")
	})

	group.Func().Params(Id("a").Id("Authorizer")).Id("Authorize").Params(
		Id("ctx").Qual("context", "Context"),
		Id("actor").Op("*").Qual(actor.Path, actor.Name),
		Id("permission").String(),
		Id("resource").Any(),
	).Bool().BlockFunc(func(group *Group) {
		group.Switch(Id("resource").Assert(Type())).BlockFunc(func(group *Group) {
			for _, resource := range resources {
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

		group.Return(False())
	})

	group.Line()

	group.Func().Id("NewAuthorizer").Params(Id("resolver").Id("Resolver")).Op("*").Id("Authorizer").Block(
		Return(Op("&").Id("Authorizer").Values(Dict{
			Id("resolver"): Id("resolver"),
		})),
	)
}
