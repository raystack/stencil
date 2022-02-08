package stencil

import (
	"github.com/goburrow/cache"
)

type store struct {
	cache.LoadingCache
}

func (s *store) getResolver(key string) (*Resolver, bool) {
	val, err := s.Get(key)
	if err != nil {
		return nil, false
	}
	return val.(*Resolver), true
}

func newStore(urls []string, loadingCache cache.LoadingCache) (*store, error) {
	s := &store{loadingCache}
	for _, url := range urls {
		if _, err := loadingCache.Get(url); err != nil {
			return s, err
		}
	}
	return s, nil
}
