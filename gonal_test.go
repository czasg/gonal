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

func TestHub_notify(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	q := queue.NewLifoMemoryQueue(0)
	_ = q.Push([]byte{1})
	h := Hub{
		Ctx:       ctx,
		CtxCancel: cancel,
	}
	_ = h.notify(Payload{})

	h = Hub{
		Ctx: context.Background(),
		Q:   q,
	}
	_ = h.notify(Payload{})
}

func TestHub_fetch(t *testing.T) {
	Bind(func(ctx context.Context, payload Payload) {
	}, Label{
		"v1": "v1",
		"v2": "v1",
	})
	hub.fetch(Label{
		"v1": "v1",
		"v2": "v1",
	})
}

func TestHub_loop(t *testing.T) {
	SetMaxConcurrent(context.Background(), runtime.NumCPU()*4, queue.NewFifoMemoryQueue(1024))
	Bind(func(ctx context.Context, payload Payload) {
		panic("test")
	}, Label{"type": "test"})
	_ = hub.Q.Push([]byte{1})
	_ = Notify(Payload{
		Label: Label{"type": "test"},
	})
	time.Sleep(time.Second)
}

func TestHub_loop2(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	SetMaxConcurrent(ctx, 0, queue.NewFifoMemoryQueue(1024))
	Bind(func(ctx context.Context, payload Payload) {
		cancel()
		time.Sleep(time.Second)
	}, Label{"type": "test"})
	_ = hub.Q.Push([]byte{1})
	_ = Notify(Payload{
		Label: Label{"type": "test"},
	})
	time.Sleep(time.Second)
}
