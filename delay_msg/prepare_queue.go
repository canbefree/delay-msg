package delay_msg

import "context"

// 预备消费队列
type PrepareQueue interface {
	Pushlish(ctx context.Context, job *Job) error
	Subscribe(ctx context.Context) chan *Job
}

var MyPrepareQueue PrepareQueue
