// Copyright (C) 2021 Creditor Corp. Group.
// See LICENSE for copying information.

package thelooper

import (
	"context"
	"errors"
	"time"
)

// ErrCtxCancelled is error returned by loop when context cancelled.
var ErrCtxCancelled = errors.New("context cancelled")

// Loop implements a controllable recurring event.
type Loop struct {
	interval time.Duration
	ticker   *time.Ticker
	stop     chan struct{}
}

// NewLoop creates a new loop with the specified interval.
func NewLoop(interval time.Duration) *Loop {
	loop := &Loop{}
	loop.SetInterval(interval)
	return loop
}

// SetInterval allows to change the interval before starting.
func (loop *Loop) SetInterval(interval time.Duration) {
	loop.interval = interval
}

// SetNextTickDuration allows to change the next tick duration after starting.
func (loop *Loop) SetNextTickDuration(nextTick time.Time) {
	loop.ticker = time.NewTicker(nextTick.Sub(time.Now().UTC()))
}

// Run runs the specified in an interval.
// Every interval `fn` is started.
func (loop *Loop) Run(ctx context.Context, fn func(ctx context.Context) error) error {
	loop.stop = make(chan struct{})
	defer close(loop.stop)

	loop.ticker = time.NewTicker(loop.interval)
	defer loop.ticker.Stop()

	if err := fn(ctx); err != nil {
		return err
	}
	for {
		select {
		case <-ctx.Done():
			return ErrCtxCancelled
		case <-loop.stop:
			return nil

		case <-loop.ticker.C:
			if err := fn(ctx); err != nil {
				return err
			}
		}
	}
}

// Close closes the loop.
func (loop *Loop) Close() {
	<-loop.stop
}
