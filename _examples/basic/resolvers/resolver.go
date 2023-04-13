package resolvers

import (
	"github.com/endigma/toucan/_examples/basic/gen/toucan"
)

type Resolver struct{}

var _ (toucan.Resolver) = (*Resolver)(nil)

func NewResolver() *Resolver {
	return &Resolver{}
}
