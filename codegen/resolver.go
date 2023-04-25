package codegen

import (
	. "github.com/dave/jennifer/jen"
	"github.com/endigma/toucan/schema"
)

func generateResourceResolver(group *Group, actor schema.Model, resource schema.ResourceSchema) error {
	group.Comment("Resolver for resource `" + resource.Name + "`")

	// Generate resolver interface
	group.Type().Id(pascal(resource.Name) + "Resolver").InterfaceFunc(func(group *Group) {
		// Role resolver
		if len(resource.Roles) > 0 {
			for _, role := range resource.Roles {
				group.Id("HasRole"+pascal(role.Name)).Params(
					Id("context").Qual("context", "Context"),
					Id("actor").Op("*").Qual(actor.Path, actor.Name),
					Id("resource").Op("*").Qual(resource.Model.Path, resource.Model.Name),
				).Bool()
			}
		}

		// Attribute resolver
		if len(resource.Attributes) > 0 {
			for _, attribute := range resource.Attributes {
				group.Id("HasAttribute"+pascal(attribute.Name)).Params(
					Id("context").Qual("context", "Context"),
					Id("resource").Op("*").Qual(resource.Model.Path, resource.Model.Name),
				).Bool()
			}
		}
	})

	return nil
}

func generateGlobalResolver(group *Group, resources []schema.ResourceSchema) {
	group.Comment("Global resolver")
	group.Type().Id("Resolver").InterfaceFunc(func(group *Group) {
		for _, resource := range resources {
			group.Id(pascal(resource.Name)).Params().Id(pascal(resource.Name) + "Resolver")
		}
	})
}
