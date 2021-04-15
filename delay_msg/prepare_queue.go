package delay_msg

// 预备消费队列
type PrepareQueue interface {
	QueueIFace
}

var MyPrepareQueue PrepareQueue
