package delay_msg

import (
	"context"
	"time"
)

type DelayPool interface {
	// 每隔多少分钟一次 从pool获取数据
	InitPrepareJob(context.Context) error
}

type RedisDataList struct{}

func NewRedisDataList() *RedisDataList {
	return &RedisDataList{}
}

type RedisZSet struct {
	// 小于当前时间 + ctrlSecond直接入列
	CtrlSecond int32
}

func NewRedisZSet() *RedisZSet {
	return &RedisZSet{}
}

type RedisPool struct {
	RedisDataList
	RedisZSet
	PrepareQueue
}

func (pool *RedisPool) InitPrepareJob(ctx context.Context) error {
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		for range ticker.C {
			// 	// range中获取前100的数据
			// 	// 遍历job中的 time 当小于当前配置时间，直接塞入ch
			// 	ch <- &Job{}
			pool.PrepareQueue.Pushlish(ctx, &Job{})
		}
	}()
	return nil
}
