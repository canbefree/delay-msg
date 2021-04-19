package delay_msg

import "context"

type QueueIFace interface {
	// 设置队列Name
	SetName(ctx context.Context, name string) error
	// 推送
	Push(ctx context.Context, job *Job) error
	// 处理
	Pop(ctx context.Context, handler func(job *Job) error) error
}

// 预备消费队列
type PrepareQueue interface {
	QueueIFace
}

// job发布队列
type TaskQueue interface {
	QueueIFace
}
