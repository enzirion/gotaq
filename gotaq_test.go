package gotaq_test

import (
	"context"
	"errors"
	"sync"
	"testing"

	"enzirion/gotaq"
)

func Test_EnqueueExecutesTask_Success(t *testing.T) {
	gq := gotaq.New()
	gq.Start()

	called := false

	var wg sync.WaitGroup
	wg.Add(1)

	q := gotaq.Queue{
		Exec: func(ctx context.Context) (any, error) {
			t.Log("called Exec")
			return true, nil
		},
		Successed: func(ctx context.Context, result any) error {
			called = true
			t.Log("called Successed")
			if result == true {
				t.Log("go true success")
				return nil
			}

			return errors.New("successed error")
		},
		Finaly: func(ctx context.Context) error {
			t.Log("called Finaly")
			wg.Done()
			return nil
		},
		Failover: func(ctx context.Context, err error) {
			t.Fatalf("Expected no error, got %v", err)
		},
	}

	gq.Enqueue(q)
	wg.Wait()

	if !called {
		t.Fatal("Exec was not called by Enqueue/Start")
	}
}

func Test_EnqueueExecutesTask_Failed(t *testing.T) {
	gq := gotaq.New()
	gq.Start()

	called := false

	var wg sync.WaitGroup
	wg.Add(1)

	q := gotaq.Queue{
		Exec: func(ctx context.Context) (any, error) {
			t.Log("called Exec")
			return nil, errors.New("exec error")
		},
		Failed: func(ctx context.Context, err error) error {
			called = true
			t.Log("called Failed")
			if err != nil {
				t.Log("got error")
				return nil
			}

			return errors.New("successed error")
		},
		Finaly: func(ctx context.Context) error {
			t.Log("called Finaly")
			wg.Done()
			return nil
		},
		Failover: func(ctx context.Context, err error) {
			t.Fatalf("Expected no error, got %v", err)
		},
	}

	gq.Enqueue(q)
	wg.Wait()

	if !called {
		t.Fatal("Exec was not called by Enqueue/Start")
	}
}
