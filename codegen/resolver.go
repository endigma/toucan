package codegen

import (
	. "github.com/dave/jennifer/jen"
	"github.com/endigma/toucan/spec"
)

func generateResourceResolver(g *Group, actor spec.QualifierSpec, resource spec.ResourceSpec) error {
	g.Comment("Resolver for resource `" + resource.Name + "`")

	// Generate resolver interface
	g.Type().Id(pascal(resource.Name) + "Resolver").InterfaceFunc(func(g *Group) {
		// Role resolver
		if len(resource.Roles) > 0 {
			g.Id("HasRole").Params(
				Id("context").Qual("context", "Context"),
				Id("actor").Op("*").Qual(actor.Path, actor.Name),
				Id("role").Id(pascal(resource.Name+"Role")),
				Id("resource").Op("*").Qual(resource.Model.Path, resource.Model.Name),
			).Bool()
		}

		// Attribute resolver
		if len(resource.Attributes) > 0 {
			g.Id("HasAttribute").Params(
				Id("context").Qual("context", "Context"),
				Id("attribute").Id(pascal(resource.Name+"Attribute")),
				Id("resource").Op("*").Qual(resource.Model.Path, resource.Model.Name),
			).Bool()
		}
	})

	return nil
}

func generateGlobalResolver(g *Group, actor spec.QualifierSpec, resources []spec.ResourceSpec) error {
	g.Comment("Global resolver")
	g.Type().Id("Resolver").InterfaceFunc(func(g *Group) {
		for _, resource := range resources {
			g.Id(pascal(resource.Name)).Params().Id(pascal(resource.Name) + "Resolver")
		}
	})

	return nil
}
