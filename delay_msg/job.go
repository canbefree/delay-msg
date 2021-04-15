package delay_msg

import "context"

const (
	JOB_STATUS_READY int32 = iota + 1
	JOB_STATUS_DELAY
	JOB_STATUS_DELETED
	JOB_STATUS_PREPARED
	JOB_STATUS_INVOKED
)

type Job struct {
	// 任务ID
	ID int64
	// 任务主题
	Topic string
	// 唤醒时间
	InvokerTime int64
	// 超时时间 为0，尝试重试
	Timeout int64

	CreatedAt int64
	UpdateAt  int64
}

type JobService struct {
}

// func (js *JobService) AddJob(context.Context){

// }

type JobRepo interface {
	AddJob(ctx context.Context, job *Job) error
	UpdateJob(ctx context.Context, job *Job) error
}
