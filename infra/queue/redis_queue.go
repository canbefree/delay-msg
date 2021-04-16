package infra_queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/canbefree/delay-msg/delay_msg"
	"github.com/gomodule/redigo/redis"
)

//
type RedisQueue struct {
	Name      string
	RedisPool *redis.Pool
}

func NewRedisQueue(name string, pool *redis.Pool) *RedisQueue {
	return &RedisQueue{
		Name:      name,
		RedisPool: pool,
	}
}

// 设置队列tag
func (q *RedisQueue) SetName(ctx context.Context, name string) error {
	q.Name = name
	return nil
}

// 发布
func (q *RedisQueue) Push(ctx context.Context, job *delay_msg.Job) error {
	conn := q.RedisPool.Get()
	defer conn.Close()
	jobS, err := json.Marshal(job)
	if err != nil {
		return err
	}
	conn.Do("LPUSH", q.getName(), jobS)
	if err != nil {
		return err
	}
	return nil
}

// 订阅
func (q *RedisQueue) PopWithCallback(ctx context.Context, handle func(*delay_msg.Job) error) error {
	conn := q.RedisPool.Get()
	defer conn.Close()
	bts, err := redis.Bytes(conn.Do("RPOP", q.getName()))
	if err != nil {
		return err
	}
	var job *delay_msg.Job
	if err := json.Unmarshal(bts, &job); err != nil {
		return err
	}
	return handle(job)
}

func (q *RedisQueue) getName() string {
	return fmt.Sprintf("queue_%v", q.Name)
}
