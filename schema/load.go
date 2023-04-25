package schema

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alecthomas/hcl/v2"
	"github.com/go-playground/mold/v4/modifiers"
	"github.com/imdario/mergo"
)

func LoadSchema(resourceCfgGlob string) (*Schema, error) {
	filenames, err := filepath.Glob(resourceCfgGlob)
	if err != nil {
		return nil, fmt.Errorf("failed to glob schema files: %w", err)
	}

	var schema Schema

	for _, filename := range filenames {
		fileSchema, err := LoadSchemaFile(filename)
		if err != nil {
			return nil, fmt.Errorf("failed to load schema file: %w", err)
		}

		if err := mergo.Merge(&schema, fileSchema, mergo.WithOverride, mergo.WithAppendSlice); err != nil {
			return nil, fmt.Errorf("failed to merge schemas: %w", err)
		}
	}

	return &schema, nil
}

func ReadSchemaFile(data []byte) (*Schema, error) {
	var schema Schema

	err := hcl.Unmarshal(data, &schema)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Use mold to transform the config
	m := modifiers.New()

	err = m.Struct(context.Background(), &schema)
	if err != nil {
		return nil, fmt.Errorf("failed to transform config: %w", err)
	}

	return &schema, nil
}

func LoadSchemaFile(path string) (*Schema, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("unable to read config: %w", err)
	}

	schema, err := ReadSchemaFile(b)
	if err != nil {
		return nil, err
	}

	return schema, err
}
