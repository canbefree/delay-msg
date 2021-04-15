package delay_msg

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDelayServer_push(t *testing.T) {
	// 测试delay_server
	s := NewDelayServer(context.TODO())
	job1 := &Job{ID: 1}
	job2 := &Job{ID: 2}
	s.push(job1)
	s.push(job2)
	rjob1, _ := s.pop()
	rjob2, _ := s.pop()
	assert.Equal(t, job1, rjob1)
	assert.Equal(t, job2, rjob2)
}
