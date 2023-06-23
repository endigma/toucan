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
}

//nolint:gocognit,cyclop
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
			Do(func(s *Statement) {
				if resource.Model != nil {
					s.Id("resource").Op("*").Qual(resource.Model.Path, resource.Model.Name)
				}
			}),
		).Error().
		BlockFunc(func(group *Group) {
			if resource.Model != nil {
				group.If(Id("resource").Op("==").Nil()).Block(
					Return(Qual("fmt", "Errorf").Call(
						Lit(fmt.Sprintf("authorize %s: resource is nil", resource.Name)),
					)),
				)
			}

			group.Var().Id("cancel").Func().Params()
			group.Id("ctx").Op(",").Id("cancel").Op("=").Qual("context", "WithCancel").Call(Id("ctx"))
			group.Defer().Id("cancel").Call().Line()

			group.Id("results").Op(":=").Make(Chan().Struct(
				Id("allow").Bool(),
				Id("source").String(),
				Id("error").Error(),
			))

			group.Var().Id("wg").Add(ConcWaitGroup()).Line()

			if len(resource.Attributes) > 0 {
				group.Switch(Id("action")).BlockFunc(func(group *Group) {
					for i, permission := range resource.Permissions {
						sources := resource.GetAttributeSources(permission)
						if len(sources) == 0 {
							continue
						}

						if i > 0 {
							group.Line()
						}

						generateAuthorizerCase(group, resource, permission, sources)
					}
				})
			}

			group.Line()

			if len(resource.Roles) > 0 {
				group.If(Id("actor").Op("!=").Nil()).Block(
					Switch(Id("action")).BlockFunc(func(group *Group) {
						for i, permission := range resource.Permissions {
							sources := resource.GetRoleSources(permission)
							if len(sources) == 0 {
								continue
							}

							if i > 0 {
								group.Line()
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
	sources []schema.PermissionSource,
) {
	group.Case(Id("Permission" + pascal(resource.Name) + pascal(perm))).
		BlockFunc(
			func(group *Group) {
				for i, source := range sources {
					if i > 0 {
						group.Line()
					}
					group.Commentf("Source: %s - %s", source.Type, source.Name)
					group.Id("wg").Dot("Go").Call(Func().Params().Block(
						List(Id("allow"), Err()).Op(":=").Id("a").Dot("resolver").Do(func(s *Statement) {
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
							Id(pascal(string(source.Type))+pascal(resource.Name)+pascal(source.Name)),
						),
						Id("results").Op("<-").Struct(
							Id("allow").Bool(),
							Id("source").String(),
							Id("error").Error(),
						).Values(Dict{
							Id("allow"):  Id("allow"),
							Id("source"): Lit(fmt.Sprintf("%s %s", source.Name, source.Type)),
							Id("error"):  Err(),
						}),
					))
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
					Do(func(s *Statement) {
						if resource.Model != nil {
							s.List(Id("resource"), Op("ok")).Op(":=").Id("resource").
								Assert(Op("*").Qual(resource.Model.Path, resource.Model.Name))
							s.Line()
							s.If(Op("!").Id("ok")).Block(
								Return(Qual("fmt", "Errorf").Call(
									Lit(fmt.Sprintf(
										"authorize: invalid resource type %%T for %s, wanted *%s.%s",
										resource.Name, resource.Model.Path, resource.Model.Name,
									)),
									Id("resource"),
								)),
							)
						} else {
							s.If(Id("resource").Op("!=").Nil()).Block(
								Return(Qual("fmt", "Errorf").Call(
									Lit("authorize: invalid resource type %T, wanted nil"),
									Id("resource"),
								)),
							)
						}
					}),
					Return(
						Id("a").
							Dot("authorize"+pascal(resource.Name)).
							Call(
								Id("ctx"), Id("actor"), Id("permission"),
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
