package codegen

import (
	. "github.com/dave/jennifer/jen"
	"github.com/endigma/toucan/schema"
)

func (gen *Generator) generateResourceResolver(file *File, resource schema.ResourceSchema) {
	file.Comment("Resolver for resource `" + resource.Name + "`")

	// Generate resolver interface
	file.Type().Id(pascal(resource.Name) + "Resolver").InterfaceFunc(func(group *Group) {
		if resource.Model != nil {
			group.Id("CacheKey").Params(Id("resource").Op("*").Qual(resource.Model.Tuple())).Id("string").Line()
		}
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
	group.Type().Id("Resolver").InterfaceFunc(func(group *Group) {
		group.Id("CacheKey").Params(Id("actor").Op("*").Qual(gen.Schema.Actor.Tuple())).Id("string").Line()

		for _, resource := range gen.Schema.Resources {
			group.Id(pascal(resource.Name)).Params().Id(pascal(resource.Name) + "Resolver")
		}
	})
}
