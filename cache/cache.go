package cache

import (
	"context"
	"strings"
	"sync"

	"github.com/endigma/toucan/decision"
	"golang.org/x/sync/singleflight"
)

type ctxKey struct{}

type Cache struct {
	mu  sync.RWMutex
	m   map[string]decision.Decision
	sfg singleflight.Group
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
	write := func(s string) {
		if strings.ContainsRune(s, 0) {
			panic("contained null byte")
		}
		sb.WriteString(s)
		sb.WriteByte(0)
	}
	write(c.ActorKey)
	write(c.Resource)
	write(c.ResourceKey)
	write(c.SourceType)
	write(c.SourceName)
	return sb.String()
}

func (c *Cache) query(keystr string) (decision.Decision, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	dec, ok := c.m[keystr]
	return dec, ok
}

func (c *Cache) insert(keystr string, dec decision.Decision) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.m[keystr] = dec
}

func (c *Cache) Query(key CacheKey, fallback func() decision.Decision) decision.Decision {
	keystr := key.string()
	dec, ok := c.query(keystr)
	if ok {
		return dec
	}

	v, err, _ := c.sfg.Do(keystr, func() (interface{}, error) {
		return fallback(), nil
	})
	if err != nil {
		// probably can't happen?
		dec = decision.Error(err)
	} else {
		dec = v.(decision.Decision)
	}

	c.insert(keystr, dec)
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

func Query(ctx context.Context, key CacheKey, fallback func() decision.Decision) decision.Decision {
	if c := FromContext(ctx); c != nil {
		return c.Query(key, fallback)
	}
	return fallback()
}
