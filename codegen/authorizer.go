package codegen

import (
	. "github.com/dave/jennifer/jen"
	"github.com/endigma/toucan/spec"
	"github.com/samber/lo"
)

func generateResourceAuthorizer(g *Group, actor spec.QualifierSpec, resource spec.ResourceSpec) error {
	g.Comment("authorizer for resource `" + resource.Name + "`").Line()

	g.Func().
		Params(
			Id("a").Id("Authorizer"),
		).
		Id("Authorize" + pascal(resource.Name)).
		Params(
			paramsForAuthorizer(actor, resource)...,
		).Bool().
		BlockFunc(func(g *Group) {
			g.If(Op("!").Id("action").Dot("Valid").Call()).Block(
				Return(False()),
			)

			g.Switch(Id("action")).BlockFunc(func(g *Group) {
				for _, perm := range resource.Permissions {
					sources := resource.GetPermissionSources(perm)

					g.Case(Id(pascal(resource.Name) + "Permission" + pascal(perm))).Block(
						Return(lo.Reduce(sources, func(s *Statement, source spec.PermissionSource, n int) *Statement {
							resolver, params := CallPermissionSource(source)
							call := Id("a").Dot("resolver").Dot(pascal(resource.Name)).Call().Dot(resolver).Add(params)

							if n == 0 {
								return s.Add(call)
							} else {
								return s.Op("||").Line().Add(call)
							}
						}, &Statement{})),
					)

				}
				g.Default().Block(Return(False()))
			})
		})

	return nil
}

func generateGlobalAuthorizer(g *Group, actor spec.QualifierSpec, resources []spec.ResourceSpec) error {
	g.Comment("Global authorizer")
	g.Type().Id("Authorizer").StructFunc(func(g *Group) {
		g.Id("resolver").Id("Resolver")
	})

	g.Func().Params(Id("a").Id("Authorizer")).Id("Authorize").Params(
		Id("ctx").Qual("context", "Context"),
		Id("actor").Op("*").Qual(actor.Path, actor.Name),
		Id("permission").String(),
		Id("resource").Any(),
	).Bool().BlockFunc(func(g *Group) {
		g.Switch(Id("resource").Assert(Type())).BlockFunc(func(g *Group) {
			for _, resource := range resources {
				g.Case(Op("*").Qual(resource.Model.Path, resource.Model.Name)).Block(
					List(Id("perm"), Id("err")).Op(":=").Id("Parse"+pascal(resource.Name)+"Permission").Call(Id("permission")),
					If(Id("err").Op("!=").Nil()).Block(
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

		g.Return(False())
	})

	g.Line()

	g.Func().Id("NewAuthorizer").Params(Id("resolver").Id("Resolver")).Id("Authorizer").Block(
		Return(Id("Authorizer").Values(Dict{
			Id("resolver"): Id("resolver"),
		})),
	)

	return nil
}
