package bucketstorage

import (
	"context"
	"fmt"

	"github.com/DanielTitkov/antibruteforce-microservice/internal/bucket"
)

type BucketStorage struct {
	rubrics []string
	M       map[string]map[string]bucket.Bucket
}

type BucketArgs struct {
	ctx      context.Context
	rate     int // to 1 timespan unit
	timespan int // seconds
}

func New(rubrics []string) (*BucketStorage, error) {
	m := make(map[string]map[string]bucket.Bucket)
	for _, s := range rubrics {
		m[s] = make(map[string]bucket.Bucket)
	}
	return &BucketStorage{rubrics: rubrics, M: m}, nil
}

func (bs *BucketStorage) Resolve(rubric, arg string, ba BucketArgs) (bool, error) {
	m, ok := bs.M[rubric]
	if !ok {
		return false, fmt.Errorf("Rubric is not present: %s", rubric)
	}
	b, ok := m[arg]
	if !ok {
		b, err := bucket.New(ba.ctx, ba.rate, ba.timespan)
		if err != nil {
			return false, err
		}
		m[arg] = b
	}
	// TODO: check if bucket is alive!
	b, _ = m[arg]
	res := b.Add()
	return res, nil
}
