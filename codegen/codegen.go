package codegen

import (
	"fmt"
	"path/filepath"

	. "github.com/dave/jennifer/jen"
	"github.com/endigma/toucan/schema"
)

type Generator struct {
	Schema *schema.Schema
	Output *OutputConfig
}

type OutputConfig struct {
	Path    string `validate:"required"`
	Package string `validate:"required"`
}

func NewGenerator(schema *schema.Schema, out *OutputConfig) *Generator {
	return &Generator{Output: out, Schema: schema}
}

func (gen *Generator) Generate() error {
	typesFile := gen.NewFile()
	resolverFile := gen.NewFile()
	authorizerFile := gen.NewFile()

	// Generate resources
	for _, resource := range gen.Schema.Resources {
		// Generate types
		gen.generateResourceTypes(typesFile, resource)

		// Generate resolver
		gen.generateResourceResolver(resolverFile, resource)

		// Generate authorizer
		gen.generateResourceAuthorizer(authorizerFile, resource)

		// Generate filter
		if resource.Model != nil {
			gen.generateResourceFilter(authorizerFile, resource)
		}
	}

	if err := typesFile.Save(filepath.Join(gen.Output.Path + "/types.go")); err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	// Generate resolver
	gen.generateResolverRoot(resolverFile.Group)

	if err := resolverFile.Save(filepath.Join(gen.Output.Path + "/resolvers.go")); err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	// Generate authorizer
	gen.generateAuthorizerRoot(authorizerFile.Group)

	if err := authorizerFile.Save(filepath.Join(gen.Output.Path + "/authorizers.go")); err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	return nil
}

func (gen *Generator) NewFile() *File {
	resourceFile := NewFile(gen.Output.Package)
	resourceFile.PackageComment("Code generated by toucan. DO NOT EDIT.")
	resourceFile.Line()

	return resourceFile
}
