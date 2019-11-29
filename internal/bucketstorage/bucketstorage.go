package bucketstorage

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/DanielTitkov/antibruteforce-microservice/internal/bucket"
)

type BucketStorage struct {
	mx      sync.RWMutex
	rubrics []string
	M       map[string]map[string]bucket.Bucket
}

type BucketArgs struct {
	ctx      context.Context
	rate     int // to 1 timespan unit
	timespan int // seconds
}

func New(rubrics []string, clean int) (*BucketStorage, error) {
	m := make(map[string]map[string]bucket.Bucket)
	for _, s := range rubrics {
		m[s] = make(map[string]bucket.Bucket)
	}
	bs := BucketStorage{rubrics: rubrics, M: m}
	if clean > 0 {
		bs.runClean(clean)
	}
	return &bs, nil
}

func (bs *BucketStorage) Resolve(rubric, arg string, ba BucketArgs) (bool, error) {
	m, ok := bs.M[rubric]
	if !ok {
		return false, fmt.Errorf("Rubric is not present: %s", rubric)
	}
	bs.mx.Lock()
	b, ok := m[arg]
	if !ok || !b.IsAlive() {
		b, err := bucket.New(ba.ctx, ba.rate, ba.timespan)
		if err != nil {
			return false, err
		}
		m[arg] = b
	}
	b, _ = m[arg]
	res := b.Add()
	bs.mx.Unlock()
	return res, nil
}

func (bs *BucketStorage) Clean() error {
	bs.mx.Lock()
	for _, v := range bs.M {
		for ki, vi := range v {
			if !vi.IsAlive() {
				delete(v, ki)
			}
		}
	}
	bs.mx.Unlock()
	return nil
}

func (bs *BucketStorage) runClean(sleepMS int) {
	go func() {
		for {
			bs.Clean()
			time.Sleep(time.Millisecond * time.Duration(sleepMS))
		}
	}()
}
