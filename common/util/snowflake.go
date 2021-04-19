package util

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/bwmarrin/snowflake"
)

var once sync.Once
var lock sync.Mutex

func GenerateSnowflake() int64 {
	var nodeIndex int64
	once.Do(func() {
		lock.Lock()
		nodeIndex = getRandNode() % (-1 ^ (-1 << snowflake.NodeBits))
		lock.Unlock()
	})
	lock.Lock()

	node, err := snowflake.NewNode(nodeIndex)
	if err != nil {
		panic(fmt.Sprintf("snowflake err:%v", err))
	}
	return node.Generate().Int64()
}

func getRandNode() int64 {
	rand.Seed(time.Now().UnixNano())
	r := rand.Int63()
	if r == 0 {
		return getRandNode()
	}
	return r
}
