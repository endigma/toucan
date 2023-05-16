package cache

import (
	"context"
	"strings"
	"sync"

	"github.com/endigma/toucan/decision"
)

type ctxKey struct{}

type Cache struct {
	mu sync.RWMutex
	m  map[string]decision.Decision
}

type CacheKey struct {
	ActorKey    string
	Resource    string
	ResourceKey string
	SourceType  string
	SourceName  string
}

func (c *CacheKey) string() string {
	var sb strings.Builder
	sb.WriteString(c.ActorKey)
	sb.WriteByte('$')
	sb.WriteString(c.Resource)
	sb.WriteByte('$')
	sb.WriteString(c.ResourceKey)
	sb.WriteByte('$')
	sb.WriteString(c.SourceType)
	sb.WriteByte('$')
	sb.WriteString(c.SourceName)
	sb.WriteByte('$')
	return sb.String()
}

func (c *Cache) Query(key CacheKey) (decision.Decision, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	dec, ok := c.m[key.string()]
	return dec, ok
}

func (c *Cache) Insert(key CacheKey, dec decision.Decision) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.m[key.string()] = dec
}

func (c *Cache) QueryOr(key CacheKey, fallback func() decision.Decision) decision.Decision {
	dec, ok := c.Query(key)
	if ok {
		return dec
	}
	dec = fallback()
	c.Insert(key, dec)
	c.m[key.string()] = dec
	return dec
}

func newCache() *Cache {
	return &Cache{m: map[string]decision.Decision{}}
}

func NewContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxKey{}, newCache())
}

func FromContext(ctx context.Context) *Cache {
	c, _ := ctx.Value(ctxKey{}).(*Cache)
	return c
}

func QueryOr(ctx context.Context, key CacheKey, fallback func() decision.Decision) decision.Decision {
	if c := FromContext(ctx); c != nil {
		return c.QueryOr(key, fallback)
	}
	return fallback()
}
