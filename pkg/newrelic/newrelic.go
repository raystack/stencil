package newrelic

//go:generate mockery --name=Service -r --case underscore --with-expecter --structname NewRelic  --filename=newrelic.go --output=./mocks

import (
	"context"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type Service interface {
	StartGenericSegment(context.Context, string) func()
}
type NewRelic struct {
}

func (nr *NewRelic) StartGenericSegment(ctx context.Context, name string) func() {
	txn := newrelic.FromContext(ctx)
	if txn == nil {
		return func() {}
	}
	gs := newrelic.Segment{
		Name: name,
	}
	gs.StartTime = txn.StartSegmentNow()
	return gs.End
}
