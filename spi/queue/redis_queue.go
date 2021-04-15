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
type RedisSubPub struct {
	topic     string
	RedisPool *redis.Pool
}

// 设置队列tag
func (q *RedisSubPub) SetTag(ctx context.Context, topic string) error {
	q.topic = topic
	return nil
}

// 发布
func (q *RedisSubPub) Publish(ctx context.Context, job *delay_msg.Job) error {
	conn := q.RedisPool.Get()
	jobS, err := json.Marshal(job)
	if err != nil {
		return err
	}
	return conn.Send("PUBLISH", q.getTopic(), jobS)
}

// 订阅
func (q *RedisSubPub) Subscribe(ctx context.Context, handle func()) (chan *delay_msg.Job, chan error) {
	jobCh := make(chan *delay_msg.Job)
	errCh := make(chan error)
	go func() {
		for {
			conn := q.RedisPool.Get()
			psc := redis.PubSubConn{Conn: conn}
			psc.Subscribe(q.getTopic())
			switch v := psc.Receive().(type) {
			case redis.Message:
				var job *delay_msg.Job
				err := json.Unmarshal(v.Data, &job)
				if err != nil {
					common.Log.Errorf(ctx, "RedisSubPub,Subscribe,err:%v", err)
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

func (q *RedisSubPub) getTopic() string {
	return fmt.Sprintf("topic_%v", q.topic)
}
