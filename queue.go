package gotaq

import "context"

// Queue は、Gotaqキューに登録する1つのタスクとそのコールバック群を表します。
type Queue struct {
	// Exec はメインの実行関数です。error==nilの場合はSuccessedが呼ばれます。
	Exec func(ctx context.Context) (any, error)
	// Successed はExecが正常終了した場合に呼ばれます。
	Successed func(ctx context.Context, result any) error
	// Failed はExecがエラーを返した場合に呼ばれます。
	Failed func(ctx context.Context, err error) error
	// Finaly はExecの成否にかかわらず必ず最後に呼ばれます。
	Finaly func(ctx context.Context) error
	// DoRetry はExecがエラーを返した場合に呼ばれ、trueを返すとリトライします。attemptsはリトライ回数です。
	DoRetry func(ctx context.Context, attempts int) bool

	// Failover はSucceed,Failed,Finalyがエラーを返した場合に呼ばれます。
	Failover func(ctx context.Context, err error)
}

// normalize は、Queueの各フィールドがnilの場合にデフォルト値を設定します。
func (q *Queue) normalize() {
	if q.Exec == nil {
		// デフォルトのExecは何もしない、ないよりマシなのでデフォ関数を入れておく
		q.Exec = func(ctx context.Context) (any, error) { return nil, nil }
	}
	if q.Successed == nil {
		q.Successed = func(ctx context.Context, result any) error { return nil }
	}
	if q.Failed == nil {
		q.Failed = func(ctx context.Context, err error) error { return nil }
	}
	if q.Finaly == nil {
		q.Finaly = func(ctx context.Context) error { return nil }
	}
	if q.DoRetry == nil {
		q.DoRetry = func(ctx context.Context, attempts int) bool { return false }
	}
	if q.Failover == nil {
		q.Failover = func(ctx context.Context, err error) {}
	}
}
