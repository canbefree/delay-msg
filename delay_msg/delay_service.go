package delay_msg

import (
	"container/list"
	"context"
	"sync"
	"time"

	"github.com/canbefree/delay-msg/infra"
)

type DelayServer struct {
	lock sync.Mutex
	List *list.List
	Pool DelayPool
	PrepareQueue
	TaskQueue
}

func NewDelayServer(ctx context.Context) *DelayServer {
	return &DelayServer{
		lock: sync.Mutex{},
		List: list.New(),
		Pool: NewRedisPool(infra.NewRedisPool()),
	}
}

func (d *DelayServer) Init() {
}

func (d *DelayServer) Run(ctx context.Context) error {
	//  获取job
	errCh := d.Pool.InitPrepareJob(ctx)
	go d.WatchPrepareJob(ctx)
	go d.Handle(ctx)
	err := <-errCh
	return err
}

func (d *DelayServer) WatchPrepareJob(ctx context.Context) error {
	jobCh := d.PrepareQueue.Pop(ctx)
	for job := range jobCh {
		d.push(job)
	}
	return nil
}

func (d *DelayServer) Handle(ctx context.Context) error {
	c := time.NewTicker(500 * time.Millisecond)
	for range c.C {
		job, err := d.pop()
		if err != nil {
			return err
		}
		d.TaskQueue.Push(ctx, job)
	}
	return nil
}

func (d *DelayServer) push(job *Job) {
	d.lock.Lock()
	d.List.PushBack(job)
	d.lock.Unlock()
}

func (d *DelayServer) pop() (*Job, error) {
	d.lock.Lock()
	elem := d.List.Front()
	if elem == nil {
		return nil, nil
	}
	jobI := d.List.Remove(elem)
	d.lock.Unlock()
	return jobI.(*Job), nil
}
