package main

import (
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/canbefree/delay-msg/delay_msg"
)

// sample for delay_msg

func main() {
	// mock
	job := &delay_msg.Job{}
	node, _ := snowflake.NewNode(1)
	job.ID = node.Generate().Int64()
	job.InvokerTime = time.Now().Unix() + 600
	job.CreatedAt = time.Now().Unix()

	// add job to db
	service := &delay_msg.JobService{}
	service

	// push into redis
	// update job

	//  redis -> queue -> update job

}
