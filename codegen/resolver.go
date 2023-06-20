package codegen

import (
	. "github.com/dave/jennifer/jen"
	"github.com/endigma/toucan/schema"
)

//nolint:gocognit
func (gen *Generator) generateResolverTypes(file *File) {
	file.Line().Type().Id("Resolver").Interface(
		Id("HasRole").Params(
			Id("ctx").Qual("context", "Context"),
			Id("actor").Op("*").Qual(gen.Schema.Actor.Path, gen.Schema.Actor.Name),
			Id("resource").Any(),
			Id("role").Id("Role"),
		).Qual("github.com/endigma/toucan/decision", "Decision"),
		Id("HasAttribute").Params(
			Id("ctx").Qual("context", "Context"),
			Id("resource").Any(),
			Id("attribute").Id("Attribute"),
		).Qual("github.com/endigma/toucan/decision", "Decision"),
	)

	file.Line().Type().Id("resolver").Struct(
		Id("root").Id("ResolverRoot"),
	)

	file.Line().Func().Params(
		Id("r").Id("resolver"),
	).Id("HasRole").Params(
		Id("ctx").Qual("context", "Context"),
		Id("actor").Op("*").Qual(gen.Schema.Actor.Path, gen.Schema.Actor.Name),
		Id("resource").Any(),
		Id("role").Id("Role"),
	).Qual("github.com/endigma/toucan/decision", "Decision").Block(
		Switch(Id("role")).BlockFunc(func(group *Group) {
			for _, resource := range gen.Schema.Resources {
				for _, role := range resource.Roles {
					group.Case(Id("Role" + pascal(resource.Name) + pascal(role.Name))).Block(
						Return(Id("r").Dot("root").Dot(pascal(resource.Name)).Call().Dot("HasRole"+pascal(role.Name)).Call(
							Id("ctx"),
							Id("actor"),
							Do(func(s *Statement) {
								if resource.Model != nil {
									s.Id("resource").Assert(Op("*").Qual(resource.Model.Path, resource.Model.Name))
								}
							}),
						)),
					)
				}
			}
			group.Default().Block(
				Return(Qual("github.com/endigma/toucan/decision", "False").Call(
					Lit("unmatched in HasRole: ").Op("+").String().Call(Id("role"))),
				),
			)
		}),
	)

	file.Line().Func().Params(
		Id("r").Id("resolver"),
	).Id("HasAttribute").Params(
		Id("ctx").Qual("context", "Context"),
		Id("resource").Any(),
		Id("attribute").Id("Attribute"),
	).Qual("github.com/endigma/toucan/decision", "Decision").Block(
		Switch(Id("attribute")).BlockFunc(func(group *Group) {
			for _, resource := range gen.Schema.Resources {
				for _, attr := range resource.Attributes {
					group.Case(Id("Attribute" + pascal(resource.Name) + pascal(attr.Name))).Block(
						Return(Id("r").Dot("root").Dot(pascal(resource.Name)).Call().Dot("HasAttribute"+pascal(attr.Name)).Call(
							Id("ctx"),
							Do(func(s *Statement) {
								if resource.Model != nil {
									s.Id("resource").Assert(Op("*").Qual(resource.Model.Path, resource.Model.Name))
								}
							}),
						)),
					)
				}
			}
			group.Default().Block(
				Return(Qual("github.com/endigma/toucan/decision", "False").Call(
					Lit("unmatched in HasAttribute: ").Op("+").String().Call(Id("attribute"))),
				),
			)
		}),
	)

	file.Line().Func().Id("NewResolver").Params(
		Id("root").Id("ResolverRoot"),
	).Id("Resolver").Block(
		Return(Id("resolver").Values(Dict{
			Id("root"): Id("root"),
		})),
	)
}

func (gen *Generator) generateResourceResolver(file *File, resource schema.ResourceSchema) {
	file.Comment("Resolver for resource `" + resource.Name + "`")

	// Generate resolver interface
	file.Type().Id(pascal(resource.Name) + "Resolver").InterfaceFunc(func(group *Group) {
		// Role resolver
		if len(resource.Roles) > 0 {
			for _, role := range resource.Roles {
				group.Id("HasRole" + pascal(role.Name)).ParamsFunc(func(group *Group) {
					group.Id("ctx").Qual("context", "Context")
					group.Id("actor").Op("*").Qual(gen.Schema.Actor.Path, gen.Schema.Actor.Name)
					if resource.Model != nil {
						group.Id("resource").Op("*").Qual(resource.Model.Path, resource.Model.Name)
					}
				}).Add(RuntimeDecision())
			}
		}

		// Attribute resolver
		if len(resource.Attributes) > 0 {
			for _, attribute := range resource.Attributes {
				group.Id("HasAttribute" + pascal(attribute.Name)).ParamsFunc(func(group *Group) {
					group.Id("ctx").Qual("context", "Context")
					if resource.Model != nil {
						group.Id("resource").Op("*").Qual(resource.Model.Path, resource.Model.Name)
					}
				}).Add(RuntimeDecision())
			}
		}
	})
}

func (gen *Generator) generateResolverRoot(group *Group) {
	group.Comment("Root Resolver")
	group.Type().Id("ResolverRoot").InterfaceFunc(func(group *Group) {
		for _, resource := range gen.Schema.Resources {
			group.Id(pascal(resource.Name)).Params().Id(pascal(resource.Name) + "Resolver")
		}
	})
}
