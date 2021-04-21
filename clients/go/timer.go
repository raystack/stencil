package stencil

import (
	"io"
	"time"
)

type timer struct {
	ticker *time.Ticker
	done   chan bool
}

func (t *timer) Close() error {
	t.ticker.Stop()
	t.done <- true
	return nil
}

func setInterval(d time.Duration, f func()) io.Closer {
	ticker := time.NewTicker(d)
	done := make(chan bool)
	go (func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				f()
			}
		}
	})()
	return &timer{ticker: ticker, done: done}
}
