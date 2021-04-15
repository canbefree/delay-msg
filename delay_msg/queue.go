package delay_msg

import "context"

type QueueIFace interface {
	// 设置队列tag
	SetTag(ctx context.Context, tag string) error
	// 发布
	Push(ctx context.Context, job *Job) error
	// 订阅
	Pop(ctx context.Context) chan *Job
}
