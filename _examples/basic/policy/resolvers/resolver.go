package resolvers

import (
	"github.com/endigma/toucan/_examples/basic/gen/toucan"
	"github.com/endigma/toucan/_examples/basic/models"
)

type Resolver struct{}

var _ (toucan.Resolver) = (*Resolver)(nil)

func NewResolver() *Resolver {
	return &Resolver{}
}

func (Resolver) CacheKey(actor *models.User) string {
	return actor.ID.String()
}
