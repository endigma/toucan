package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/endigma/toucan/api"
	"github.com/endigma/toucan/codegen/config"
	"github.com/endigma/toucan/codegen/spec"
)

func main() {
	cfg, err := config.LoadConfig("toucan.hcl")
	if err != nil {
		log.Fatal(err)
	}

	spec, err := spec.FromConfig(cfg)
	if err != nil {
		log.Fatal(err)
	}

	err = api.Generate(spec)
	if err != nil {
		log.Fatal(err)
	}
}
