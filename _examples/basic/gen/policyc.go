package main

import (
	"github.com/endigma/toucan/api"
	"github.com/endigma/toucan/codegen"
	"github.com/endigma/toucan/schema"
)

func main() {
	loadedSchema, err := schema.LoadSchema("policy/schema/*.hcl")
	if err != nil {
		panic(err)
	}

	// You can modify the schema after loading it.
	loadedSchema.Actor = schema.Model{
		Path: "github.com/endigma/toucan/_examples/basic/models",
		Name: "User",
	}

	err = api.Generate(loadedSchema, &codegen.OutputConfig{
		Path:    "./gen/toucan",
		Package: "policy",
	})
	if err != nil {
		panic(err)
	}

	return
}
