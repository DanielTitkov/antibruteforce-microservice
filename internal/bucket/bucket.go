package bucket

import (
	"context"
	"time"
)

// Bucket implements the leaky bucket logic
type Bucket struct {
	ctx      context.Context
	queue    chan struct{}
	rate     int // to 1 timespan unit
	timespan int // seconds
}

func (b *Bucket) leak() {
	sleeptime := 1000 * b.timespan / b.rate
	for {
		select {
		case <-b.queue:
			time.Sleep(time.Duration(sleeptime) * time.Millisecond)
		case <-b.ctx.Done():
			return
		}
	}
}

// New function returns new bucket and starts leaking
func New(ctx context.Context, rate, timespan int) (Bucket, error) {
	q := make(chan struct{}, rate)
	b := Bucket{ctx: ctx, queue: q, rate: rate, timespan: timespan}
	go b.leak()
	return b, nil // maybe return pointer?
}

// Add method adds one object to bucket's queue
func (b *Bucket) Add() bool {
	select {
	case <-b.ctx.Done():
		return false // do what?
	case b.queue <- struct{}{}:
		return true
	default:
		return false
	}
}

// Kill method stops bucket leaking routine
func (b *Bucket) Kill() error {
	return nil
}

func (b *Bucket) IsAlive() bool {
	select {
	case <-b.ctx.Done():
		return false
	default:
		return true
	}
}
