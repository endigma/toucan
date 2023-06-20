package codegen

import (
	"fmt"

	. "github.com/dave/jennifer/jen"
	"github.com/endigma/toucan/schema"
)

var (
	ConcWaitGroup   = func() *Statement { return Qual("github.com/sourcegraph/conc", "WaitGroup") }
	RuntimeDecision = func() *Statement { return Qual("github.com/endigma/toucan/decision", "Decision") }
	RuntimeTrue     = func() *Statement { return Qual("github.com/endigma/toucan/decision", "True") }
	RuntimeFalse    = func() *Statement { return Qual("github.com/endigma/toucan/decision", "False") }
	RuntimeError    = func() *Statement { return Qual("github.com/endigma/toucan/decision", "Error") }
)

func (gen *Generator) generateAuthorizerTypes(file *File) {
	file.Line().Type().Id("Authorizer").Interface(
		Id("Authorize").Params(
			Id("ctx").Qual("context", "Context"),
			Id("actor").Op("*").Qual(gen.Schema.Actor.Path, gen.Schema.Actor.Name),
			Id("permission").String(),
			Id("resourceType").String(),
			Id("resource").Interface(),
		).Qual("github.com/endigma/toucan/decision", "Decision"),
	)

	file.Line().Type().Id("AuthorizerFunc").Func().Params(
		Id("ctx").Qual("context", "Context"),
		Id("actor").Op("*").Qual(gen.Schema.Actor.Path, gen.Schema.Actor.Name),
		Id("permission").String(),
		Id("resourceType").String(),
		Id("resource").Interface(),
	).Qual("github.com/endigma/toucan/decision", "Decision")

	file.Line().Func().Params(
		Id("af").Id("AuthorizerFunc"),
	).Id("Authorize").Params(
		Id("ctx").Qual("context", "Context"),
		Id("actor").Op("*").Qual(gen.Schema.Actor.Path, gen.Schema.Actor.Name),
		Id("permission").String(),
		Id("resourceType").String(),
		Id("resource").Interface(),
	).Qual("github.com/endigma/toucan/decision", "Decision").Block(
		Return(Id("af").Call(
			Id("ctx"),
			Id("actor"),
			Id("permission"),
			Id("resourceType"),
			Id("resource"),
		)),
	)
}

func (gen *Generator) generateResourceAuthorizer(file *File, resource schema.ResourceSchema) {
	file.Line().Func().
		Params(
			Id("a").Id("authorizer"),
		).
		Id("authorize" + pascal(resource.Name)).
		ParamsFunc(paramsForAuthorizer(gen.Schema.Actor, resource)).Add(RuntimeDecision()).
		BlockFunc(func(group *Group) {
			group.If(Op("!").Id("action").Dot("Valid").Call()).Block(
				Return(
					RuntimeError().Call(Id(fmt.Sprintf("ErrInvalid%sPermission", pascal(resource.Name)))),
				),
			).Line()

			if resource.Model != nil {
				group.If(Id("resource").Op("==").Nil()).Block(
					Return().Add(RuntimeFalse()).Call(Lit("unmatched")),
				)
			}

			group.Var().Id("cancel").Func().Params()
			group.Id("ctx").Op(",").Id("cancel").Op("=").Qual("context", "WithCancel").Call(Id("ctx"))
			group.Defer().Id("cancel").Call().Line()

			group.Id("results").Op(":=").Make(Chan().Add(RuntimeDecision())).Line()
			group.Var().Id("wg").Add(ConcWaitGroup()).Line()

			if len(resource.Attributes) > 0 {
				group.Switch(Id("action")).BlockFunc(func(group *Group) {
					for _, permission := range resource.Permissions {
						sources := resource.GetAttributeSources(permission)
						if len(sources) == 0 {
							continue
						}

						generateAuthorizerCase(group, resource, permission, sources)
					}
				})
			}

			if len(resource.Roles) > 0 {
				group.If(Id("actor").Op("!=").Nil()).Block(
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

			group.Go().Func().Params().Block(
				Id("wg").Dot("Wait").Call(),
				Close(Id("results")),
			).Call().Line()

			group.Var().Id("allowReason").Add(String())
			group.Var().Id("denyReasons").Index().Add(String())

			group.For(Id("result").Op(":=").Range().Id("results")).Block(
				If(Id("result").Dot("Reason")).Op("==").Lit("").Block(
					Id("result").Dot("Reason").Op("=").Lit("unspecified"),
				),
				If(Id("result").Dot("Allow")).Block(
					Id("cancel").Call(),
					Id("allowReason").Op("=").Id("result").Dot("Reason"),
				).Else().Block(
					Id("denyReasons").Op("=").Append(Id("denyReasons"), Id("result").Dot("Reason")),
				),
			).Line()

			group.If(Id("allowReason").Op("!=").Lit("")).Block(
				Return(RuntimeTrue().Call(Id("allowReason"))),
			).Else().Block(
				Id("result").Op(":=").Add(RuntimeFalse()).Call(Qual("strings", "Join").Call(Id("denyReasons"), Lit(", "))),
				If(Id("result").Dot("Reason")).Op("==").Lit("").Block(
					Id("result").Dot("Reason").Op("=").Lit("unspecified"),
				),
				Return(Id("result")),
			)
		})

	file.Line()
}

func generateAuthorizerCase(
	group *Group,
	resource schema.ResourceSchema,
	perm string,
	sources []schema.PermissionSource,
) {
	group.Case(Id(pascal(resource.Name) + "Permission" + pascal(perm))).
		BlockFunc(
			func(group *Group) {
				for _, source := range sources {
					group.Commentf("Source: %s - %s", source.Type, source.Name)
					group.Id("wg").Dot("Go").Call(Func().Params().Block(
						Id("results").Op("<-").Id("a").Dot("resolver").Do(func(s *Statement) {
							switch source.Type {
							case schema.PermissionTypeRole:
								s.Dot("HasRole")
							case schema.PermissionTypeAttribute:
								s.Dot("HasAttribute")
							}
						}).Call(
							Id("ctx"),
							Do(func(s *Statement) {
								if source.Type == schema.PermissionTypeRole {
									s.Id("actor")
								}
							}),
							Do(func(s *Statement) {
								if resource.Model != nil {
									s.Id("resource")
								} else {
									s.Nil()
								}
							}),
							Lit(resource.Name),
							Lit(source.Name),
						),
					))
					group.Line()
				}
			})
}

func (gen *Generator) generateAuthorizerRoot(group *Group) {
	group.Comment("Authorizer")
	group.Type().Id("authorizer").StructFunc(func(group *Group) {
		group.Id("resolver").Id("Resolver")
	})

	group.Func().Params(Id("a").Id("authorizer")).Id("Authorize").Params(
		Id("ctx").Qual("context", "Context"),
		Id("actor").Op("*").Qual(gen.Schema.Actor.Path, gen.Schema.Actor.Name),
		Id("permission").String(),
		Id("resourceType").String(),
		Id("resource").Any(),
	).Add(RuntimeDecision()).BlockFunc(func(group *Group) {
		group.Switch(Id("resourceType")).BlockFunc(func(group *Group) {
			for _, resource := range gen.Schema.Resources {
				group.Case(Lit(resource.Name)).Block(
					List(Id("perm"), Id("err")).Op(":=").Id("Parse"+pascal(resource.Name)+"Permission").Call(Id("permission")),
					Do(func(s *Statement) {
						if resource.Model != nil {
							s.List(Id("resource"), Op("_")).Op(":=").Id("resource").
								Assert(Op("*").Qual(resource.Model.Path, resource.Model.Name))
						}
					}),
					If(Id("err").Op("!=").Nil()).Block(
						Return(RuntimeError().Call(Id("err"))),
					),
					Return(
						Id("a").
							Dot("authorize"+pascal(resource.Name)).
							Call(
								Id("ctx"), Id("actor"), Id("perm"),
								Do(func(s *Statement) {
									if resource.Model != nil {
										s.Id("resource")
									}
								}),
							),
					),
				)
			}
		}).Line()

		group.Return(RuntimeFalse().Call(Lit("unmatched")))
	})

	group.Line()

	group.Func().Id("NewAuthorizer").Params(Id("resolver").Id("Resolver")).Id("Authorizer").Block(
		Return(Id("authorizer").Values(Dict{
			Id("resolver"): Id("resolver"),
		})),
	)
}
