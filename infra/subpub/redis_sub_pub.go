package infra_subpub

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
	Topic     string
	RedisPool *redis.Pool
}

// 设置队列tag
func (q *RedisSubPub) SetName(ctx context.Context, topic string) error {
	q.Topic = topic
	return nil
}

// 发布
func (q *RedisSubPub) Publish(ctx context.Context, job *delay_msg.Job) error {
	conn := q.RedisPool.Get()
	defer conn.Close()
	jobS, err := json.Marshal(job)
	if err != nil {
		return err
	}
	return conn.Send("PUBLISH", q.getTopic(), jobS)
}

// 订阅
func (q *RedisSubPub) Subscribe(ctx context.Context, handle func(*delay_msg.Job) error) chan error {
	errCh := make(chan error)

	conn := q.RedisPool.Get()
	psc := redis.PubSubConn{Conn: conn}
	psc.Subscribe(q.getTopic())

	go func() {
		defer conn.Close()
		for {
			switch v := psc.Receive().(type) {
			case redis.Message:
				var job *delay_msg.Job
				common.Log.Infof(ctx, "message")
				err := json.Unmarshal(v.Data, &job)
				if err != nil {
					common.Log.Errorf(ctx, "RedisSubPub,Subscribe,err:%v", err)
					errCh <- err
				}
				go func() {
					if err := handle(job); err != nil {
						errCh <- err
					}

				}()
			case redis.Subscription:
				common.Log.Infof(ctx, "channel:%v", v.Channel)
			case redis.Error:
				errCh <- v
			}
		}
	}()
	return errCh
}

func (q *RedisSubPub) getTopic() string {
	return fmt.Sprintf("topic_%v", q.Topic)
}
