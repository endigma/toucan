package api

import (
	"os"

	"github.com/endigma/toucan/codegen"
	"github.com/endigma/toucan/spec"
)

func Generate(spec *spec.Spec) error {
	generator := codegen.NewGenerator(spec)

	err := wipeOutputDir(spec)
	if err != nil {
		return err
	}

	return generator.Generate()
}

func Validate(spec *spec.Spec) error {
	return spec.Validate()
}

// Delete all files in output directory
func wipeOutputDir(spec *spec.Spec) error {
	if err := os.RemoveAll(spec.Output.Path); err != nil {
		return err
	}

	if err := os.MkdirAll(spec.Output.Path, os.ModePerm); err != nil {
		return err
	}

	return nil
}
