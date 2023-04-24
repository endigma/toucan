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
			group.Id("HasRole").Params(
				Id("context").Qual("context", "Context"),
				Id("actor").Op("*").Qual(actor.Path, actor.Name),
				Id("role").Id(pascal(resource.Name+"Role")),
				Id("resource").Op("*").Qual(resource.Model.Path, resource.Model.Name),
			).Bool()
		}

		// Attribute resolver
		if len(resource.Attributes) > 0 {
			group.Id("HasAttribute").Params(
				Id("context").Qual("context", "Context"),
				Id("attribute").Id(pascal(resource.Name+"Attribute")),
				Id("resource").Op("*").Qual(resource.Model.Path, resource.Model.Name),
			).Bool()
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
