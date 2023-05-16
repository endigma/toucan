package codegen

import (
	. "github.com/dave/jennifer/jen"
	"github.com/endigma/toucan/schema"
)

func (gen *Generator) generateResourceResolver(file *File, resource schema.ResourceSchema) {
	file.Comment("Resolver for resource `" + resource.Name + "`")

	// Generate resolver interface
	file.Type().Id(pascal(resource.Name) + "Resolver").InterfaceFunc(func(group *Group) {
		// Role resolver
		if len(resource.Roles) > 0 {
			for _, role := range resource.Roles {
				group.Id("HasRole"+pascal(role.Name)).Params(
					Id("ctx").Qual("context", "Context"),
					Id("actor").Op("*").Qual(gen.Schema.Actor.Path, gen.Schema.Actor.Name),
					Id("resource").Op("*").Qual(resource.Model.Path, resource.Model.Name),
				).Add(RuntimeDecision())
			}
		}

		// Attribute resolver
		if len(resource.Attributes) > 0 {
			for _, attribute := range resource.Attributes {
				group.Id("HasAttribute"+pascal(attribute.Name)).Params(
					Id("ctx").Qual("context", "Context"),
					Id("resource").Op("*").Qual(resource.Model.Path, resource.Model.Name),
				).Add(RuntimeDecision())
			}
		}
	})
}

func (gen *Generator) generateGlobalResolver(file *File) {
	file.Comment("Resolver for global scope")

	// Generate resolver interface
	file.Type().Id("GlobalResolver").InterfaceFunc(func(group *Group) {
		// Role resolver
		if len(gen.Schema.Global.Roles) > 0 {
			for _, role := range gen.Schema.Global.Roles {
				group.Id("HasRole"+pascal(role.Name)).Params(
					Id("ctx").Qual("context", "Context"),
					Id("actor").Op("*").Qual(gen.Schema.Actor.Path, gen.Schema.Actor.Name),
				).Add(RuntimeDecision())
			}
		}

		// Attribute resolver
		if len(gen.Schema.Global.Attributes) > 0 {
			for _, attribute := range gen.Schema.Global.Attributes {
				group.Id("HasAttribute" + pascal(attribute.Name)).Params(
					Id("ctx").Qual("context", "Context"),
				).Add(RuntimeDecision())
			}
		}
	})
}

func (gen *Generator) generateResolverRoot(group *Group) {
	group.Comment("Root Resolver")
	group.Type().Id("Resolver").InterfaceFunc(func(group *Group) {
		group.Id("Global").Params().Id("GlobalResolver").Line()
		for _, resource := range gen.Schema.Resources {
			group.Id(pascal(resource.Name)).Params().Id(pascal(resource.Name) + "Resolver")
		}
	})
}
