package stencil

import (
	"io"
	"sync"
	"time"
)

type loaderFunc func(string) (*Resolver, error)
type timer struct {
	ticker *time.Ticker
	done   chan bool
}

func (t *timer) Close() error {
	t.ticker.Stop()
	t.done <- true
	return nil
}

func setInterval(d time.Duration, f func(), waitForReader <-chan bool) io.Closer {
	ticker := time.NewTicker(d)
	done := make(chan bool)
	go (func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				// wait for access
				refresh := <-waitForReader
				if refresh {
					f()
				}
			}
		}
	})()
	return &timer{ticker: ticker, done: done}
}

type store struct {
	autoRefresh bool
	timer       io.Closer
	access      chan bool
	loader      loaderFunc
	url         string
	data        *Resolver
	lock        sync.RWMutex
}

func (s *store) refresh() {
	val, err := s.loader(s.url)
	if err == nil {
		s.lock.Lock()
		defer s.lock.Unlock()
		s.data = val
	}
}

func (s *store) notify() {
	select {
	case s.access <- true:
	default:
	}
}

func (s *store) getResolver() (*Resolver, bool) {
	s.notify()
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.data, s.data != nil
}

func (s *store) Close() {
	close(s.access)
	if s.timer != nil {
		s.timer.Close()
	}
}

func newStore(url string, options Options) (*store, error) {
	loader := options.RefreshStrategy.getLoader(options)
	s := &store{loader: loader, access: make(chan bool), url: url, autoRefresh: options.AutoRefresh}
	val, err := loader(url)
	if err != nil {
		return s, err
	}
	s.data = val
	if options.AutoRefresh {
		s.timer = setInterval(options.RefreshInterval, s.refresh, s.access)
	}
	return s, nil
}
