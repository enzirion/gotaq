package gotaq

import (
	"context"
)

type Gotaq struct {
	ch   chan Queue
	stop chan struct{}
}

// New は新しいGotaqインスタンスを生成します。
func New() *Gotaq {
	return &Gotaq{
		ch:   make(chan Queue, 100), // バッファサイズは適宜調整
		stop: make(chan struct{}),
	}
}

// Start を呼ぶと、GotaqはEnqueueされたタスクの監視・実行を開始します。
func (g *Gotaq) Start() {
	go func() {
		for {
			select {
			case q := <-g.ch:
				go g.handleQueue(q)
			case <-g.stop:
				return
			}
		}
	}()
}

// handleQueue はタスクの実行・リトライ・コールバック呼び出しを行います。
func (g *Gotaq) handleQueue(q Queue) {
	ctx := context.Background()
	attempts := 0
	q.normalize()
	for {
		result, err := q.Exec(ctx)
		if err == nil {
			if err := q.Successed(ctx, result); err != nil {
				q.Failover(ctx, err)
			}
			break
		} else {
			if err := q.Failed(ctx, err); err != nil {
				q.Failover(ctx, err)
			}
			if q.DoRetry(ctx, attempts) {
				attempts++
				continue
			}
			break
		}
	}
	if err := q.Finaly(ctx); err != nil {
		q.Failover(ctx, err)
	}
}

// Enqueue は新しいタスクをキューに追加します。
func (g *Gotaq) Enqueue(q Queue) {
	g.ch <- q
}

// GracefulStop を呼ぶと、全てのタスクが完了するまで待機し、その後停止します。
func (g *Gotaq) GracefulStop() {
	close(g.stop)
}
