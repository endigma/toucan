package codegen

import (
	"path/filepath"

	. "github.com/dave/jennifer/jen"
	"github.com/endigma/toucan/codegen/spec"
	"github.com/iancoleman/strcase"
	"github.com/samber/lo"
)

type SpecGenerator struct {
	*spec.Spec
}

func NewSpecGenerator(spec *spec.Spec) *SpecGenerator {
	return &SpecGenerator{spec}
}

func (gen *SpecGenerator) Generate() error {
	for _, resource := range gen.Resources {
		if err := gen.generateResource(resource); err != nil {
			return err
		}
	}

	if err := gen.generateGlobalAuthorizer(); err != nil {
		return err
	}

	if err := gen.generateGlobalResolver(); err != nil {
		return err
	}

	return nil
}

func (gen *SpecGenerator) generateGlobalResolver() error {
	f := gen.newFile()

	f.Comment("Global resolver")

	f.Type().Id("Resolver").InterfaceFunc(func(g *Group) {
		for _, resource := range gen.Resources {
			g.Id(resource.CamelName()).Params().Id(resource.CamelName() + "Resolver")
		}
	})

	return gen.saveFile(f, "resolver")
}

func (gen *SpecGenerator) generateGlobalAuthorizer() error {
	f := gen.newFile()

	f.Comment("Global authorizer")

	f.Type().Id("Authorizer").StructFunc(func(g *Group) {
		g.Id("resolver").Id("Resolver")
	})

	f.Func().Params(Id("a").Id("Authorizer")).Id("Authorize").Params(
		Id("actor").Op("*").Qual(gen.Actor.Path, gen.Actor.Name),
		Id("permission").String(),
		Id("resource").Any(),
	).Bool().BlockFunc(func(g *Group) {
		g.Switch(Id("resource").Assert(Type())).BlockFunc(func(g *Group) {
			for _, resource := range gen.Resources {
				g.Case(Op("*").Qual(resource.Model.Path, resource.Model.Name)).Block(
					Return(
						Id("a").
							Dot("Authorize"+resource.CamelName()).
							Call(
								Id("actor"),
								Id("permission"),
								Id("resource").
									Assert(
										Op("*").Qual(resource.Model.Path, resource.Model.Name),
									),
							),
					),
				)
			}
			g.Default().Return(False())
		})
	})

	f.Line()

	f.Func().Id("NewAuthorizer").Params(Id("resolver").Id("Resolver")).Id("Authorizer").Block(
		Return(Id("Authorizer").Values(Dict{
			Id("resolver"): Id("resolver"),
		})),
	)

	return gen.saveFile(f, "authorizer")
}

func (gen *SpecGenerator) newFile() *File {
	f := NewFile(gen.Output.Package)
	f.PackageComment("Code generated by toucan, DO NOT EDIT.")
	f.Line()

	return f
}

func (gen *SpecGenerator) saveFile(file *File, name string) error {
	return file.Save(filepath.Join(gen.Output.Path, name+".go"))
}

func (gen *SpecGenerator) generateResource(resource spec.ResourceSpec) error {
	f := gen.newFile()

	permissionType := func() *Statement { return Id(resource.CamelName() + "Permission") }
	resolverType := Id(resource.CamelName() + "Resolver")
	authorizerType := Id("Authorize" + resource.CamelName())

	formatPermission := func(permission string) Code {
		return Id(resource.CamelName() + "Permission" + strcase.ToCamel(permission))
	}

	formatRoleResolver := func(role spec.RoleSpec) Code {
		return Id("HasRole" + strcase.ToCamel(role.Name))
	}

	formatAttrResolver := func(attr spec.AttributeSpec) Code {
		return Id("HasAttribute" + strcase.ToCamel(attr.Name))
	}

	f.Type().Add(permissionType()).String()
	f.Const().DefsFunc(func(g *Group) {
		for _, perm := range resource.Permissions {
			g.Add(formatPermission(perm)).Add(permissionType()).Op("=").Lit(perm)
		}
	})

	f.Commentf("Resolver for resource %s", resource.Name)
	f.Type().Add(resolverType).InterfaceFunc(func(g *Group) {
		for _, role := range resource.Roles {
			g.Add(formatRoleResolver(role)).Params(
				Id("actor").Op("*").Qual(gen.Actor.Path, gen.Actor.Name),
				Id("resource").Op("*").Qual(resource.Model.Path, resource.Model.Name),
			).Bool()
		}

		g.Line()

		for _, attribute := range resource.Attributes {
			g.Add(formatAttrResolver(attribute)).Params(
				Id("resource").Op("*").Qual(resource.Model.Path, resource.Model.Name),
			).Bool()
		}
	})

	f.Line()

	f.Func().Params(
		Id("a").Id("Authorizer"),
	).Add(authorizerType).Params(
		Id("actor").Op("*").Qual(gen.Actor.Path, gen.Actor.Name),
		Id("permission").String(),
		Id("resource").Op("*").Qual(resource.Model.Path, resource.Model.Name),
	).Bool().BlockFunc(func(g *Group) {
		g.Switch(permissionType().Call(Id("permission"))).BlockFunc(func(g *Group) {
			for _, permission := range resource.Permissions {
				g.Case(formatPermission(permission)).BlockFunc(func(g *Group) {
					sources := getPermissionSources(resource, permission)

					g.Return(lo.Reduce(sources, func(s *Statement, source PermissionSource, n int) *Statement {
						call := Id("a").Dot("resolver").Dot(resource.CamelName()).Call().Dot(source.CallName()).Add(source.CallParams())

						if n == 0 {
							return s.Add(call)
						} else {
							return s.Op("||").Line().Add(call)
						}
					}, &Statement{}))
				})
			}

			g.Default().Block(Return(False()))
		})
	})

	f.Line()

	return gen.saveFile(f, strcase.ToSnake(resource.Name))
}

type PermissionSource struct {
	Type string // role, attribute
	Name string
}

func (p PermissionSource) CallName() string {
	switch p.Type {
	case "role":
		return "HasRole" + strcase.ToCamel(p.Name)
	case "attribute":
		return "HasAttribute" + strcase.ToCamel(p.Name)
	}

	return ""
}

func (p PermissionSource) CallParams() *Statement {
	switch p.Type {
	case "role":
		return Call(Id("actor"), Id("resource"))
	case "attribute":
		return Call(Id("resource"))
	}

	return Null()
}

func getPermissionSources(resource spec.ResourceSpec, permission string) []PermissionSource {
	var sources []PermissionSource = lo.Union(
		lo.FilterMap(resource.Attributes, func(attr spec.AttributeSpec, _ int) (PermissionSource, bool) {
			if lo.Contains(attr.Permissions, permission) {
				return PermissionSource{
					Type: "attribute",
					Name: strcase.ToCamel(attr.Name),
				}, true
			}

			return PermissionSource{}, false
		}),
		lo.FilterMap(resource.Roles, func(role spec.RoleSpec, _ int) (PermissionSource, bool) {
			if lo.Contains(role.Permissions, permission) {
				return PermissionSource{
					Type: "role",
					Name: strcase.ToCamel(role.Name),
				}, true
			}

			return PermissionSource{}, false
		}),
	)

	return sources
}