package spi_queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/canbefree/delay-msg/common"
	"github.com/canbefree/delay-msg/delay_msg"
	"github.com/gomodule/redigo/redis"
)

//
type RedisQueue struct {
	Name      string
	RedisPool *redis.Pool
}

// 设置队列tag
func (q *RedisQueue) SetName(ctx context.Context, name string) error {
	q.Name = name
	return nil
}

// 发布
func (q *RedisQueue) Publish(ctx context.Context, job *delay_msg.Job) error {
	conn := q.RedisPool.Get()
	jobS, err := json.Marshal(job)
	if err != nil {
		return err
	}
	return conn.Send("PUBLISH", q.getTopic(), jobS)
}

// 订阅
func (q *RedisQueue) Subscribe(ctx context.Context, handle func()) (chan *delay_msg.Job, chan error) {
	jobCh := make(chan *delay_msg.Job)
	errCh := make(chan error)

	conn := q.RedisPool.Get()
	psc := redis.PubSubConn{Conn: conn}
	psc.Subscribe(q.getTopic())

	go func() {
		for {
			switch v := psc.Receive().(type) {
			case redis.Message:
				var job *delay_msg.Job
				err := json.Unmarshal(v.Data, &job)
				if err != nil {
					common.Log.Errorf(ctx, "RedisQueue,Subscribe,err:%v", err)
					errCh <- err
				}
				jobCh <- job
			case redis.Subscription:
			case redis.Error:
				errCh <- v
			}
		}
	}()
	return jobCh, errCh
}

func (q *RedisQueue) getName() string {
	return fmt.Sprintf("topic_%v", q.Name)
}
