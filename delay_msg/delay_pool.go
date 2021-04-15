package delay_msg

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/gomodule/redigo/redis"
)

const EXTRA_EXPIRED_TIME = 120

var (
	ErrInovkerTimeInvaild = errors.New("invoke time out of range")
)

type DelayPool interface {
	// 每隔多少分钟一次 从pool获取数据
	InitPrepareJob(context.Context) chan error
}

type RedisDataList struct {
	RedisPool *redis.Pool
}

func NewRedisDataList(redisPool *redis.Pool) *RedisDataList {
	return &RedisDataList{
		RedisPool: redisPool,
	}
}

func (l *RedisDataList) Get(ctx context.Context, jobID int64) (job *Job, err error) {
	conn := l.RedisPool.Get()
	defer conn.Close()
	bts, err := redis.Bytes(conn.Do("GET", jobID))
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bts, &job)
	if err != nil {
		return nil, err
	}
	return
}

func (l *RedisDataList) Set(ctx context.Context, job *Job) error {
	conn := l.RedisPool.Get()
	defer conn.Close()
	bts, err := json.Marshal(job)
	if err != nil {
		return err
	}
	if job.InvokerTime < time.Now().Unix() {
		return ErrInovkerTimeInvaild
	}
	expired := job.InvokerTime + EXTRA_EXPIRED_TIME - time.Now().Unix()
	_, err = conn.Do("SETEX", job.ID, expired, bts)
	if err != nil {
		return err
	}
	return nil
}

type RedisZSet struct {
	// 小于当前时间 + ctrlSecond直接入列
	CtrlSecond int64
	RedisPool  *redis.Pool
	Key        string
}

func NewRedisZSet(redisPool *redis.Pool) *RedisZSet {
	return &RedisZSet{
		RedisPool: redisPool,
		Key:       "delay_msg_zset",
	}
}

func (z *RedisZSet) Range(ctx context.Context, start, end int32) ([]int64, error) {
	conn, err := z.RedisPool.GetContext(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	rely, err := redis.Int64s(conn.Do("ZRANGE", z.Key, start, end))
	if err != nil {
		return nil, err
	}
	return rely, nil
}

func (z *RedisZSet) Add(ctx context.Context, jobID int64, metric int64) error {
	conn, err := z.RedisPool.GetContext(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.Do("ZADD", z.Key, metric, jobID)
	return err
}

func (z *RedisZSet) Remove(ctx context.Context, jobId int64) error {
	conn, err := z.RedisPool.GetContext(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.Do("ZREM", z.Key, jobId)
	if err != nil {
		return err
	}
	return nil
}

func (z *RedisZSet) Empty(ctx context.Context) error {
	conn, err := z.RedisPool.GetContext(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.Do("DEL", z.Key)
	return err
}

type RedisPool struct {
	*RedisDataList
	*RedisZSet
	PrepareQueue
}

func NewRedisPool(redisPool *redis.Pool) *RedisPool {
	redisDataList := NewRedisDataList(redisPool)
	redisZSet := NewRedisZSet(redisPool)
	return &RedisPool{
		RedisDataList: redisDataList,
		RedisZSet:     redisZSet,
	}
}

func (pool *RedisPool) InitPrepareJob(ctx context.Context) chan error {
	ch := make(chan error)
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		for range ticker.C {
			jobs, err := pool.getRecentJobs(ctx)
			if err != nil {
				ch <- err
			}
			for _, job := range jobs {
				// publish
				err := pool.PrepareQueue.Push(ctx, job)
				if err == nil {
					pool.RedisZSet.Remove(ctx, job.ID)
				}
			}
		}
		ch <- nil
	}()
	return ch
}

// range中获取前100的数据
// 遍历job中的 time 当小于当前配置时间，直接返回【】jObs
func (pool *RedisPool) getRecentJobs(ctx context.Context) ([]*Job, error) {
	var jobs = make([]*Job, 0)
	jobIds, err := pool.RedisZSet.Range(ctx, 0, 100)
	if err != nil {
		return nil, err
	}
	nowFunc := func() int64 {
		return time.Now().Unix()
	}
	for _, v := range jobIds {
		job, err := pool.RedisDataList.Get(ctx, v)
		if err != nil {
			return nil, err
		}
		if job != nil {
			jobs = append(jobs, job)
		} else {
			// TODO redis数据丢失(根据状态码重新从DB获取数据)

		}
		if nowFunc()+pool.CtrlSecond > job.InvokerTime {
			// prepare ctrlsecond jobs
			break
		}
	}
	return jobs, nil
}
