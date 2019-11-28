package bucket

import (
	"time"
)

// Bucket implements the leaky bucket logic
type Bucket struct {
	queue    chan struct{}
	kill     chan struct{}
	rate     int // to 1 timespan unit
	timespan int // seconds
}

func (b *Bucket) leak() {
	sleeptime := 1000 * b.timespan / b.rate
	for {
		select {
		case <-b.queue:
			time.Sleep(time.Duration(sleeptime) * time.Millisecond)
		case <-b.kill:
			return
		}
	}
}

// New function returns new bucket and starts leaking
func New(rate, timespan int) (Bucket, error) {
	q, k := make(chan struct{}, rate), make(chan struct{}, rate)
	b := Bucket{queue: q, kill: k, rate: rate, timespan: timespan}
	go b.leak()
	return b, nil
}

// Add method adds one object to bucket's queue
func (b *Bucket) Add() bool {
	select {
	case b.queue <- struct{}{}:
		return true
	default:
		return false
	}
}

// Kill method stops bucket leaking routine
func (b *Bucket) Kill() error {
	close(b.kill)
	return nil
}
