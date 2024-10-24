package assets

import (
	"context"
	"github.com/hashicorp/golang-lru/v2/expirable"
)

type lookupCache struct {
	baseLookup Lookup
	lru        *expirable.LRU[Params, Result]
}

func (l lookupCache) Lookup(ctx context.Context, params Params) (Result, error) {
	if cached, ok := l.lru.Get(params); ok {
		return cached, nil
	}
	result, err := l.baseLookup.Lookup(ctx, params)
	if err == nil {
		l.lru.Add(params, result)
	}
	return result, err
}
