package delay_msg

import (
	"context"
	"testing"
	"time"

	"github.com/canbefree/delay-msg/spi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type RedisTestSuite struct {
	suite.Suite
	redisZSet     *RedisZSet
	redisDataList *RedisDataList
	ctx           context.Context
}

func (suite *RedisTestSuite) SetupTest() {
	redisPool := spi.NewRedisPool()
	suite.redisZSet = &RedisZSet{
		RedisPool: redisPool,
		Key:       "test_key",
	}
	suite.redisDataList = &RedisDataList{
		RedisPool: redisPool,
	}
	suite.ctx = context.TODO()
	suite.redisZSet.Add(suite.ctx, 1, 123)
}

func (suite *RedisTestSuite) TearDownTest() {
	suite.redisZSet.Empty(suite.ctx)
}

func (suite *RedisTestSuite) TestRedisZSet_Range() {
	jobIDs, err := suite.redisZSet.Range(suite.ctx, 0, 10)
	assert.Equal(suite.T(), nil, err)
	assert.Equal(suite.T(), []int64{1}, jobIDs)
}

func (suite *RedisTestSuite) TestRedisZSet_Remove() {
	err := suite.redisZSet.Remove(suite.ctx, 1)
	assert.Equal(suite.T(), nil, err)
}

func (suite *RedisTestSuite) TestRedisDataList_SGet() {
	job := &Job{
		ID:          1,
		InvokerTime: time.Now().Unix() + 600,
	}
	err := suite.redisDataList.Set(suite.ctx, job)
	assert.Equal(suite.T(), nil, err, "set err")
	rjob, err := suite.redisDataList.Get(suite.ctx, job.ID)
	assert.Equal(suite.T(), nil, err, "get err")
	assert.Equal(suite.T(), job, rjob)
}

func TestRedisTestSuite(t *testing.T) {
	suite.Run(t, new(RedisTestSuite))
}
