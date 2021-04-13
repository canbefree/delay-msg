package delay_msg

// job发布队列
type TaskQueue interface {
	Pushlish(job *Job) error
	Subscribe() <-chan *Job
}

var MyTaskQueue TaskQueue
