package codegen

import (
	"fmt"

	"github.com/dave/jennifer/jen"
	. "github.com/dave/jennifer/jen"
	"github.com/endigma/toucan/schema"
)

var (
	RuntimeDecision = func() *jen.Statement { return Qual("github.com/endigma/toucan/decision", "Decision") }
	RuntimeSkip     = func() *jen.Statement { return Qual("github.com/endigma/toucan/decision", "Skip") }
	RuntimeAllow    = func() *jen.Statement { return Qual("github.com/endigma/toucan/decision", "Allow") }
	RuntimeError    = func() *jen.Statement { return Qual("github.com/endigma/toucan/decision", "Error") }
)

func (gen *Generator) generateResourceAuthorizer(file *File, resource schema.ResourceSchema) error {
	file.Func().
		Params(
			Id("a").Id("Authorizer"),
		).
		Id("Authorize" + pascal(resource.Name)).
		Params(
			paramsForAuthorizer(gen.Schema.Actor, resource)...,
		).Add(RuntimeDecision()).
		BlockFunc(func(group *Group) {
			group.Id("resolver").Op(":=").Id("a").Dot("resolver").Dot(pascal(resource.Name)).Call().Line()

			group.If(Op("!").Id("action").Dot("Valid").Call()).Block(
				Return(
					RuntimeError().Call(Id(fmt.Sprintf("ErrInvalid%sPermission", pascal(resource.Name)))),
				),
			).Line()

			if len(resource.Attributes) > 0 {
				group.If(Id("resource").Op("!=").Nil()).Block(
					Switch(Id("action")).BlockFunc(func(group *Group) {
						for _, permission := range resource.Permissions {
							sources := resource.GetAttributeSources(permission)
							if len(sources) == 0 {
								continue
							}

							generateAuthorizerCase(group, resource, permission, sources)
						}
					}),
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

							generateAuthorizerCase(group, resource, permission, sources)
						}
					}),
				).Line()
			}

			group.Return(RuntimeSkip().Call(Lit("unmatched")))
		})

	file.Line()

	return nil
}

func generateAuthorizerCase(group *Group, res schema.ResourceSchema, perm string, sources []schema.PermissionSource) {
	group.Case(Id(pascal(res.Name) + "Permission" + pascal(perm))).
		BlockFunc(
			func(group *Group) {
				for _, source := range sources {
					resolver, params := CallPermissionSource(source)

					group.Commentf("Source: %s - %s", source.Type, source.Name)
					group.If(Id("result").Op(":=").Id("resolver").Dot(resolver).Add(params), Id("result").Dot("Allow")).Block(
						Return(Id("result")),
					)
					group.Line()
				}
			})

	// 	Return(lo.Reduce(sources, func(statement *Statement, source schema.PermissionSource, number int) *Statement {
	// 		resolver, params := CallPermissionSource(source)

	// 		call := Id("resolver").Dot(resolver).Add(params)

	// 		if number == 0 {
	// 			return statement.Add(call)
	// 		} else {
	// 			return statement.Op("||").Line().Add(call)
	// 		}
	// 	}, &Statement{})),
	// )
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
	).Add(RuntimeDecision()).BlockFunc(func(group *Group) {
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

		group.Return(RuntimeSkip().Call(Lit("unmatched")))
	})

	group.Line()

	group.Func().Id("NewAuthorizer").Params(Id("resolver").Id("Resolver")).Op("*").Id("Authorizer").Block(
		Return(Op("&").Id("Authorizer").Values(Dict{
			Id("resolver"): Id("resolver"),
		})),
	)
}
