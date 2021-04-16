package infra_queue

import (
	"context"
	"fmt"
	"testing"

	"github.com/canbefree/delay-msg/delay_msg"
	"github.com/canbefree/delay-msg/infra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TestRedisQueue struct {
	suite.Suite
	ctx        context.Context
	redisQueue *RedisQueue
}

func (s *TestRedisQueue) SetupTest() {
	pool := infra.NewRedisPool()
	s.redisQueue = &RedisQueue{
		Name:      "test_queue",
		RedisPool: pool,
	}
	s.ctx = context.TODO()
}

func (s *TestRedisQueue) TestRedisQueue() {
	var err error
	var fn = func(job *delay_msg.Job) error {
		fmt.Println(job.ID)
		return nil
	}
	err = s.redisQueue.Push(s.ctx, &delay_msg.Job{ID: 1})
	assert.Equal(s.T(), nil, err)
	s.redisQueue.Pop(s.ctx, fn)
}

func TestRedisTestSuite(t *testing.T) {
	suite.Run(t, new(TestRedisQueue))
}
