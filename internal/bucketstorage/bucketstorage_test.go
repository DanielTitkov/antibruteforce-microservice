package bucketstorage

import (
	"context"
	"testing"
	"time"
)

func TestOverfillInStorage(t *testing.T) {
	rubrics := []string{"ip", "login"}
	bs, err := New(rubrics)
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
	bs, err := New(rubrics)
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
