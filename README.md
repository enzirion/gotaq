# Gotaq

Gotaqは、メモリ上にタスクをキューイングし、非同期に順次実行できるシンプルなTask Queueライブラリです。

## インストール

```
go get enzirion/gotaq
```

## 使い方サンプル

```go
package main

import (
	"context"
	"fmt"
	"time"

	"enzirion/gotaq"
)

func main() {
	gq := gotaq.New()
	gq.Start()

	q := gotaq.Queue{
		Exec: func(ctx context.Context) (any, error) {
			fmt.Println("Exec start")
			time.Sleep(1 * time.Second)
			fmt.Println("Exec end")
			return "result", nil
		},
		Successed: func(ctx context.Context, result any) (any, error) {
			fmt.Println("Successed:", result)
			return nil, nil
		},
		Failed: func(ctx context.Context, err error) {
			fmt.Println("Failed:", err)
		},
		Finaly: func(ctx context.Context) error {
			fmt.Println("Finaly")
			return nil
		},
		DoRetry: func(ctx context.Context, attempts int) bool {
			return attempts < 3 // 3回までリトライする
		},
	}

	gq.Enqueue(q)

	// Stopを読んでもタスクが完了しているかは保証されないため注意
	gq.Stop()
}
```

## 機能
- タスクの非同期実行
- 成功/失敗/最終処理/リトライのコールバック

## ライセンス
MIT
