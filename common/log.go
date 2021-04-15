package common

import (
	"context"
	"fmt"
)

type LogIFace interface {
	Errorf(ctx context.Context, format string, params ...interface{})
}

type DefaultLog struct {
}

func (log *DefaultLog) Errorf(ctx context.Context, format string, params ...interface{}) {
	fmt.Errorf(format, params...)
}

var Log = new(DefaultLog)
