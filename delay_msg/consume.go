package delay_msg

import (
	"container/list"
	"sync"
)

// type ConsumeServiceIFace interface {
// 	// 获取需要消费的job
// 	GetHeadJob() *Job
// 	// 消费出栈,获取下一个
// 	Next() error
// 	// 执行
// 	Run() error
// }

type ConsumeQueue struct {
	lock sync.Cond
	List list.List
}

func (q *ConsumeQueue) Insert(jobId int64, execTime int64) {
}

type ConsumeService struct {
	Queue ConsumeQueue
}

type ConsumeQueueIFace interface {
	Insert(*Job) error
	Pop() *Job
}

func NewConsumeClin() *ConsumeService {
	return &ConsumeService{}
}

func (c *ConsumeService) Run() {
}
