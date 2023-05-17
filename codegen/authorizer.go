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

	RuntimeCache    = func() *Statement { return Qual("github.com/endigma/toucan/cache", "Cache") }
	RuntimeCacheKey = func() *Statement { return Qual("github.com/endigma/toucan/cache", "CacheKey") }
	RuntimeQuery    = func() *Statement { return Qual("github.com/endigma/toucan/cache", "Query") }
)

func (gen *Generator) generateResourceAuthorizer(file *File, resource schema.ResourceSchema) {
	file.Func().
		Params(
			Id("a").Id("Authorizer"),
		).
		Id("Authorize" + pascal(resource.Name)).
		ParamsFunc(paramsForAuthorizer(gen.Schema.Actor, resource)).Add(RuntimeDecision()).
		BlockFunc(func(group *Group) {
			group.Id("resolver").Op(":=").Id("a").Dot(pascal(resource.Name)).Call().Line()

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

						generateAuthorizerCase(group, resource.Name, permission, sources, resource.Model != nil)
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

							generateAuthorizerCase(group, resource.Name, permission, sources, resource.Model != nil)
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

func generateAuthorizerCase(group *Group, name string, perm string, sources []schema.PermissionSource, hasModel bool) {
	group.Case(Id(pascal(name) + "Permission" + pascal(perm))).
		BlockFunc(
			func(group *Group) {
				for _, source := range sources {
					group.Commentf("Source: %s - %s", source.Type, source.Name)
					group.Id("wg").Dot("Go").Call(Func().Params().Block(
						Id("results").Op("<-").Add(RuntimeQuery()).Call(
							Id("ctx"),
							Add(RuntimeCacheKey()).BlockFunc(func(group *Group) {
								group.Id("ActorKey").Op(":").Do(func(s *Statement) {
									switch source.Type {
									case "role":
										s.Id("a").Dot("CacheKey").Call(Id("actor"))
									case "attribute":
										s.Lit("")
									}
								}).Op(",")
								group.Id("Resource").Op(":").Lit(name).Op(",")
								group.Id("ResourceKey").Op(":").Do(func(s *Statement) {
									if hasModel {
										s.Id("resolver").Dot("CacheKey").Call(Id("resource"))
									} else {
										s.Lit("")
									}
								}).Op(",")
								group.Id("SourceType").Op(":").Lit(string(source.Type)).Op(",")
								group.Id("SourceName").Op(":").Lit(source.Name).Op(",")
							}),
							Func().Params().Add(RuntimeDecision()).BlockFunc(func(group *Group) {
								group.Return(
									Id("resolver").
										Do(func(s *Statement) {
											switch source.Type {
											case "role":
												s.Dot("HasRole" + pascal(source.Name))
											case "attribute":
												s.Dot("HasAttribute" + pascal(source.Name))
											}
										}).
										Call(
											Id("ctx"),
											Do(func(s *Statement) {
												if source.Type == "role" {
													s.Id("actor")
												}
											}),
											Do(func(s *Statement) {
												if hasModel {
													s.Id("resource")
												}
											}),
										),
								)
							}),
						),
					))
					group.Line()
				}
			})
}

func (gen *Generator) generateAuthorizerRoot(group *Group) {
	group.Comment("Authorizer")
	group.Type().Id("Authorizer").StructFunc(func(group *Group) {
		group.Id("Resolver")
	})

	group.Func().Params(Id("a").Id("Authorizer")).Id("Authorize").Params(
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
							s.List(Id("resource"), Op("_")).Op(":=").Id("resource").Assert(Op("*").Qual(resource.Model.Path, resource.Model.Name))
						}
					}),
					If(Id("err").Op("!=").Nil()).Block(
						Return(RuntimeError().Call(Id("err"))),
					),
					Return(
						Id("a").
							Dot("Authorize"+pascal(resource.Name)).
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

	group.Func().Id("NewAuthorizer").Params(Id("resolver").Id("Resolver")).Op("*").Id("Authorizer").Block(
		Return(Op("&").Id("Authorizer").Values(Dict{
			Id("Resolver"): Id("resolver"),
		})),
	)
}
