package api

import (
	"fmt"
	"os"

	"github.com/endigma/toucan/codegen"
	"github.com/endigma/toucan/schema"
)

func Generate(schema *schema.Schema, output *codegen.OutputConfig) error {
	if err := schema.Validate(); err != nil {
		return fmt.Errorf("failed to validate schema: %w", err)
	}

	generator := codegen.NewGenerator(schema, output)

	// Delete all files in the output directory.
	err := wipeDir(output.Path)
	if err != nil {
		return err
	}

	err = generator.Generate()
	if err != nil {
		return fmt.Errorf("failed to generate code: %w", err)
	}

	return nil
}

// Delete all files in a directory.
func wipeDir(dir string) error {
	if err := os.RemoveAll(dir); err != nil {
		return fmt.Errorf("failed to remove directory %q: %w", dir, err)
	}

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory %q: %w", dir, err)
	}

	return nil
}
