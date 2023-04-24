package schema

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alecthomas/hcl/v2"
	"github.com/go-playground/mold/v4/modifiers"
	"github.com/spewerspew/spew"
)

func LoadSchema(resourceCfgGlob string) (*Schema, error) {
	filenames, err := filepath.Glob(resourceCfgGlob)
	if err != nil {
		return nil, fmt.Errorf("failed to glob resource configs: %w", err)
	}

	spew.Dump(filenames)

	return &Schema{}, nil
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

	// err = validate.Struct(&schema)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to validate config: %w", err)
	// }

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
