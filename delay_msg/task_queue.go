package delay_msg

// job发布队列
type TaskQueue interface {
	QueueIFace
}

var MyTaskQueue TaskQueue
