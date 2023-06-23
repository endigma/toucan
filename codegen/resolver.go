package codegen

import (
	"fmt"

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
		).Params(Bool(), Error()),
		Id("HasAttribute").Params(
			Id("ctx").Qual("context", "Context"),
			Id("resource").Any(),
			Id("attribute").Id("Attribute"),
		).Params(Bool(), Error()),
	)

	file.Line().Type().Id("ResolverFuncs").Struct(
		Id("Role").Func().Params(
			Id("ctx").Id("context").Dot("Context"),
			Id("actor").Op("*").Qual(gen.Schema.Actor.Path, gen.Schema.Actor.Name),
			Id("resource").Any(), Id("role").Id("Role"),
		).Add(Params(Bool(), Error())),
		Id("Attribute").Func().Params(
			Id("ctx").Id("context").Dot("Context"),
			Id("resource").Any(),
			Id("attribute").Id("Attribute"),
		).Add(Params(Bool(), Error())),
	)

	file.Line().Func().Params(Id("fs").Id("ResolverFuncs")).Id("HasRole").Params(
		Id("ctx").Id("context").Dot("Context"),
		Id("actor").Op("*").Qual(gen.Schema.Actor.Path, gen.Schema.Actor.Name),
		Id("resource").Any(),
		Id("role").Id("Role"),
	).Add(Params(Bool(), Error())).Block(
		Return().Id("fs").Dot("Role").Call(Id("ctx"), Id("actor"), Id("resource"), Id("role")),
	)

	file.Line().Func().Params(Id("fs").Id("ResolverFuncs")).Id("HasAttribute").Params(
		Id("ctx").Id("context").Dot("Context"),
		Id("resource").Any(),
		Id("attribute").Id("Attribute"),
	).Add(Params(Bool(), Error())).Block(
		Return().Id("fs").Dot("Attribute").Call(Id("ctx"), Id("resource"), Id("attribute")),
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
	).Params(Bool(), Error()).Block(
		Switch(Id("role")).BlockFunc(func(group *Group) {
			for _, resource := range gen.Schema.Resources {
				for _, role := range resource.Roles {
					group.Case(Id("Role" + pascal(resource.Name) + pascal(role.Name))).BlockFunc(func(group *Group) {
						if resource.Model == nil {
							group.If(Id("resource").Op("!=").Nil()).Block(
								Return(List(False(), Qual("fmt", "Errorf").Call(
									Lit("HasRole: invalid resource type %T, wanted nil"), Id("resource"),
								))),
							)
						} else {
							group.List(Id(resource.Name), Id("ok")).Op(":=").
								Id("resource").Assert(Op("*").Qual(resource.Model.Path, resource.Model.Name))
							group.If(Op("!").Id("ok")).Block(
								Return(List(False(), Qual("fmt", "Errorf").Call(
									Lit(fmt.Sprintf(
										"HasRole: invalid resource type %%T, wanted *%s.%s",
										resource.Model.Path,
										resource.Model.Name,
									)), Id("resource"),
								))),
							)
							group.If(Id(resource.Name).Op("==").Nil()).Block(
								Return(List(False(), Qual("fmt", "Errorf").Call(
									Lit(fmt.Sprintf("HasRole: got nil %s", resource.Name)),
								))),
							)
						}
						group.Return(Id("r").Dot("root").Dot(pascal(resource.Name)).Call().Dot("HasRole"+pascal(role.Name)).Call(
							Id("ctx"),
							Id("actor"),
							Do(func(s *Statement) {
								if resource.Model != nil {
									s.Id(resource.Name)
								}
							}),
						))
					})
				}
			}
		}),
		Return(List(False(), Qual("fmt", "Errorf").Call(
			Lit("HasRole: unmatched: %s: %w"), Id("role"), Id("Deny"),
		))),
	)

	file.Line().Func().Params(
		Id("r").Id("resolver"),
	).Id("HasAttribute").Params(
		Id("ctx").Qual("context", "Context"),
		Id("resource").Any(),
		Id("attribute").Id("Attribute"),
	).Params(Bool(), Error()).Block(
		Switch(Id("attribute")).BlockFunc(func(group *Group) {
			for _, resource := range gen.Schema.Resources {
				for _, attr := range resource.Attributes {
					group.Case(Id("Attribute" + pascal(resource.Name) + pascal(attr.Name))).BlockFunc(func(group *Group) {
						if resource.Model != nil {
							group.List(Id(resource.Name), Id("ok")).Op(":=").
								Id("resource").Assert(Op("*").Qual(resource.Model.Path, resource.Model.Name))
							group.If(Op("!").Id("ok")).Block(
								Return(List(False(), Qual("fmt", "Errorf").Call(
									Lit(
										fmt.Sprintf(
											"HasAttribute: invalid resource type %%T, wanted *%s.%s",
											resource.Model.Path,
											resource.Model.Name,
										),
									), Id("resource"),
								))),
							)
							group.If(Id(resource.Name).Op("==").Nil()).Block(
								Return(List(False(), Qual("fmt", "Errorf").Call(
									Lit(fmt.Sprintf("HasRole: got nil %s", resource.Name)),
								))),
							)
						} else {
							group.If(Id("resource").Op("!=").Nil()).Block(
								Return(List(False(), Qual("fmt", "Errorf").Call(
									Lit("HasAttribute: invalid resource type %T, wanted nil"), Id("resource"),
								))),
							)
						}
						group.Return(Id("r").Dot("root").Dot(pascal(resource.Name)).Call().Dot("HasAttribute"+pascal(attr.Name)).Call(
							Id("ctx"),
							Do(func(s *Statement) {
								if resource.Model != nil {
									s.Id(resource.Name)
								}
							}),
						))
					})
				}
			}
		}),
		Return(False(), Qual("fmt", "Errorf").Call(
			Lit("HasAttribute: unmatched: %s: %w"), Id("attribute"), Id("Deny")),
		),
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
				}).Add(Params(Bool(), Error()))
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
				}).Add(Params(Bool(), Error()))
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
