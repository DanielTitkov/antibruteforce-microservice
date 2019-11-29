package bucketstorage

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestOverfillInStorage(t *testing.T) {
	rubrics := []string{"ip", "login"}
	bs, err := New(rubrics, 0)
	if err != nil {
		t.Errorf("storage creation failed: %v", err)
	}

	for i := 0; i < 20; i++ {
		ctx := context.Background()
		ba := BucketArgs{ctx, 10, 60}
		res, err := bs.Resolve("ip", "123.123.123.123", ba)
		if err != nil {
			t.Errorf("error occured during resolving at iteration %d: %v", i, err)
		}
		if i < 10 && res != true {
			t.Errorf("expeted true, got %v on iteration %d", res, i)
		}
		if i >= 10 && res != false {
			t.Errorf("expeted false, got %v on iteration %d", res, i)
		}
	}
}

func TestLeakingResultInMultirubricStorage(t *testing.T) {
	rubrics := []string{"ip", "login"}
	bs, err := New(rubrics, 0)
	if err != nil {
		t.Errorf("storage creation failed: %v", err)
	}

	for j := 0; j < 2; j++ { // two times
		for i := 0; i < 90; i++ {
			ctx := context.Background()
			ba := BucketArgs{ctx, 100, 1}
			for _, r := range rubrics {
				for _, a := range []string{"123.123", "321.321", "999.999"} {
					res, err := bs.Resolve(r, a, ba)
					if err != nil {
						t.Errorf("error occured during resolving at iteration %d-%d rub %s arg %s: %v", i, j, r, a, err)
					}
					if res != true {
						t.Errorf("expeted true, got %v on iteration %d-%d rub %s arg %s", res, i, j, r, a)
					}
				}
			}
		}
		time.Sleep(time.Second) // wait till buckets' queues leak out
	}
}

func TestDeadBucketManualClean(t *testing.T) {
	rubrics := []string{"ip"}
	bs, err := New(rubrics, 0)
	if err != nil {
		t.Errorf("storage creation failed: %v", err)
	}
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	ba := BucketArgs{ctx, 100, 1}
	_, _ = bs.Resolve("ip", "foo", ba)
	if _, ok := bs.M["ip"]["foo"]; !ok {
		t.Errorf("expected bucket to exist after resolve, got %v", ok)
	}
	cancel()
	err = bs.Clean()
	if err != nil {
		t.Errorf("error occured during clean")
	}
	if _, ok := bs.M["ip"]["foo"]; ok {
		t.Errorf("expected bucket to be absent after clean, got %v", ok)
	}
}

func TestDeadBucketAutoClean(t *testing.T) {
	rubrics := []string{"ip"}
	bs, err := New(rubrics, 10)
	if err != nil {
		t.Errorf("storage creation failed: %v", err)
	}
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	ba := BucketArgs{ctx, 100, 1}
	bs.Resolve("ip", "foo", ba)
	if _, ok := bs.M["ip"]["foo"]; !ok {
		t.Errorf("expected bucket to exist after resolve, got %v", ok)
	}
	cancel()
	time.Sleep(time.Millisecond * 20)
	if _, ok := bs.M["ip"]["foo"]; ok {
		t.Errorf("expected bucket to be absent after clean, got %v", ok)
	}
}

func TestOverfillWithConcurrency(t *testing.T) {
	rubrics := []string{"ip", "login"}
	bs, err := New(rubrics, 0)
	if err != nil {
		t.Errorf("storage creation failed: %v", err)
	}
	var passed uint64
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		ctx := context.Background()
		ba := BucketArgs{ctx, 10, 60}
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			res, err := bs.Resolve("ip", "123.123.123.123", ba)
			if err != nil {
				t.Errorf("error occured during resolving at iteration %d: %v", i, err)
			}
			if res {
				atomic.AddUint64(&passed, 1)
			}
		}(i)
	}
	wg.Wait()

	if passed != 10 && passed != 11 { // 10 or 11 because one of elements may leak before all routines do their stuff
		t.Errorf("expected 10 or 11, got %d", passed)
	}
}

func TestLeakingWithConcurrency(t *testing.T) {
	rubrics := []string{"ip", "login"}
	bs, err := New(rubrics, 0)
	if err != nil {
		t.Errorf("storage creation failed: %v", err)
	}
	for j := 0; j < 2; j++ { // two times
		var wg sync.WaitGroup
		for i := 0; i < 90; i++ {
			ctx := context.Background()
			ba := BucketArgs{ctx, 100, 1}
			for _, r := range rubrics {
				for _, a := range []string{"123.123", "321.321", "999.999"} {
					wg.Add(1)
					go func(r, a string) {
						defer wg.Done()
						res, err := bs.Resolve(r, a, ba)
						if err != nil {
							t.Errorf("error occured during resolving at iteration %d-%d rub %s arg %s: %v", i, j, r, a, err)
						}
						if res != true {
							t.Errorf("expeted true, got %v on iteration %d-%d rub %s arg %s", res, i, j, r, a)
						}
					}(r, a)
				}
			}
		}
		wg.Wait()
		time.Sleep(time.Second) // wait till buckets' queues leak out
	}
}
