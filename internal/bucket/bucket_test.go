package bucket

import (
	"context"
	"testing"
	"time"
)

func TestBucketOverfill(t *testing.T) {
	ctx := context.Background()
	b, err := New(ctx, 10, 60)
	if err != nil {
		t.Errorf("bucket creating failed with error: %v", err)
	}
	for i := 0; i < 20; i++ { // test will fail if this number is bigger, because leak manages to get one value from the channel
		res := b.Add()
		if i < 10 && res != true {
			t.Errorf("expeted true, got %v on iteration %d", res, i)
		}
		if i >= 10 && res != false {
			t.Errorf("expeted false, got %v on iteration %d", res, i)
		}
	}
}

func TestBucketLeakingResult(t *testing.T) {
	ctx := context.Background()
	b, err := New(ctx, 100, 1)
	if err != nil {
		t.Errorf("bucket creating failed with error: %v", err)
	}

	for j := 0; j < 2; j++ {
		for i := 0; i < 90; i++ {
			res := b.Add()
			if res != true {
				t.Errorf("expeted true, got %v on iteration %d-%d", res, i, j)
			}
		}
		time.Sleep(time.Second)
	}
}

func TestBuckerLeakingQueueLen(t *testing.T) {
	ctx := context.Background()
	b, err := New(ctx, 100, 1)
	if err != nil {
		t.Errorf("bucket creating failed with error: %v", err)
	}

	for i := 0; i < 100; i++ {
		_ = b.Add()
	}
	if res := len(b.queue); res != 100 {
		t.Errorf("expected len of 100, got %d", res)
	}
	time.Sleep(time.Second + time.Millisecond*50)
	if res := len(b.queue); res != 0 {
		t.Errorf("expected empty queue, got %d", res)
	}
}

func TestBucketIsAlive(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	b, err := New(ctx, 100, 1)
	if err != nil {
		t.Errorf("bucket creating failed with error: %v", err)
	}
	if res := b.IsAlive(); res != true {
		t.Errorf("expected true before cancelation, got %v", res)
	}
	cancel()
	if res := b.IsAlive(); res != false {
		t.Errorf("expected false after cancelation, got %v", res)
	}
}
