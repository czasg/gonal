package gonal

import (
	"context"
	"errors"
	"github.com/czasg/go-queue"
	"runtime"
	"testing"
	"time"
)

func TestNotify(t *testing.T) {
	var label Label
	Bind(func(ctx context.Context, payload Payload) {
		label = payload.Label
	}, Label{
		"type": "test",
	})
	label2 := Label{
		"type":     "test",
		"property": "test",
	}
	err := Notify(Payload{
		Label: label2,
	})
	time.Sleep(time.Millisecond)
	if err != nil {
		t.Error(err)
	}
	for k, v := range label2 {
		vv, ok := label[k]
		if !ok {
			t.Error(k)
		}
		if v != vv {
			t.Error(v, vv)
		}
	}
}

func Test_notifyQueuePopErr(t *testing.T) {
	notifyQueuePopErr(errors.New("test"))
}

func Test_notifyJsonErr(t *testing.T) {
	notifyJsonErr(errors.New("test"))
}

func Test_notifyHandlerPanic(t *testing.T) {
	notifyHandlerPanic("test")
}

func TestSetMaxConcurrent(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	SetMaxConcurrent(ctx, runtime.NumCPU()*4, queue.NewFifoMemoryQueue(1024))
	time.Sleep(time.Millisecond)
	cancel()
}
