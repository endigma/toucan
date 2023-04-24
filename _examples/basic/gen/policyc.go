package main

import (
	"github.com/endigma/toucan/config"
)

func main() {
	_, err := config.LoadSchema("policy/schema/*")
	if err != nil {
		panic(err)
	}

	// cfg, err := config.LoadConfig("toucan.hcl")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// spec, err := spec.FromConfig(cfg)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// err = api.Generate(spec, &cfg.)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	return
}
