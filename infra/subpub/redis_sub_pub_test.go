package infra_subpub

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/canbefree/delay-msg/delay_msg"
	"github.com/canbefree/delay-msg/infra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TestSuiteSubPub struct {
	suite.Suite
	ctx    context.Context
	subPub *RedisSubPub
}

func (s *TestSuiteSubPub) SetupTest() {
	pool := infra.NewRedisPool()
	s.subPub = &RedisSubPub{
		Topic:     "test_topic",
		RedisPool: pool,
	}
	s.ctx = context.TODO()
}

func (s *TestSuiteSubPub) TestRedisSubPub() {
	var err error
	var fn = func(job *delay_msg.Job) error {
		fmt.Println(job.ID)
		return nil
	}
	//
	s.subPub.Subscribe(s.ctx, fn)

	for i := 0; i < 10; i++ {
		time.Sleep(100 * time.Millisecond)
		err = s.subPub.Publish(s.ctx, &delay_msg.Job{
			ID: time.Now().Unix(),
		})
		assert.Equal(s.T(), nil, err)
	}

	assert.Equal(s.T(), 1, 1)

}

func TestRedisTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuiteSubPub))
}
