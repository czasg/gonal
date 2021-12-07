package gonal

import (
    "context"
    "errors"
    "github.com/czasg/go-queue"
    "reflect"
    "testing"
    "time"
)

func assertErr(t *testing.T, err1, err2 error) {
    if err1 != err2 {
        t.Error("failure", err1, err2)
    }
}

func assertInterface(t *testing.T, data1, data2 interface{}) {
    if !reflect.DeepEqual(data1, data2) {
        t.Error("failure", data1, data2)
    }
}

func sleep(ctx context.Context, payload Payload) {
    time.Sleep(time.Second)
}

func Test_Gonal(t *testing.T) {
    {
        g := &Gonal{
            LabelsMatcher: map[string][]Handler{},
            C:             make(chan struct{}, 1),
        }
        _ = g.SetContext(context.Background())
        g.Close()
        assertErr(t, g.Notify(Payload{}), context.Canceled)
    }
    {
        g := &Gonal{
            LabelsMatcher: map[string][]Handler{},
            C:             make(chan struct{}, 1),
        }
        assertErr(t, g.SetContext(context.Background()), nil)
        assertErr(t, g.SetConcurrent(0), nil)
        assertErr(t, g.SetQueue(queue.NewFifoMemoryQueue(2)), nil)
        g.Bind(Label{"test": "test"}, sleep, sleep, sleep)
        g.Bind(Label{"test": "test"})
        assertInterface(t, len(g.Fetch(Label{"test": "test"})), 1)
        assertInterface(t, len(g.Fetch(Label{"test": "none"})), 0)
        assertErr(t, g.Q.Push([]byte{1}), nil)
        assertErr(t, g.Notify(Payload{Label: Label{"test": "test"}}), nil)
        _ = g.Notify(Payload{Label: Label{"test": "test"}})
        _ = g.Notify(Payload{Label: Label{"test": "test"}})
        _ = g.Notify(Payload{Label: Label{"test": "test"}})
        _ = g.Notify(Payload{Label: Label{"test": "test"}})
        assertErr(t, g.SetContext(context.Background()), ErrRunning)
        assertErr(t, g.SetConcurrent(0), ErrRunning)
        assertErr(t, g.SetQueue(queue.NewFifoMemoryQueue(1)), ErrRunning)
        time.Sleep(time.Millisecond * 40)
        assertErr(t, g.Notify(Payload{Label: Label{"test": "test"}}), nil)
        g.Cancel()
    }
    {
        notifyQueuePopErr(errors.New("test"))
        notifyQueuePopErr(errors.New("test"))
        notifyQueuePopErr(errors.New("test"))
    }
    {
        notifyJsonErr(errors.New("test"))
        notifyJsonErr(errors.New("test"))
        notifyJsonErr(errors.New("test"))
    }
    {
        notifyHandlerPanic("test")
        notifyHandlerPanic("test")
        notifyHandlerPanic("test")
    }
    {
        Fetch(Label{"test": "test"})
        Bind(Label{"test": "test"}, func(ctx context.Context, payload Payload) {
            time.Sleep(time.Second)
        })
        Bind(Label{"test": "panic"}, func(ctx context.Context, payload Payload) {
            panic("test")
        })
        assertErr(t, Notify(Payload{Label: Label{"test": "test"}}), nil)
        assertErr(t, Notify(Payload{Label: Label{"test": "test"}}), nil)
        assertErr(t, Notify(Payload{Label: Label{"test": "test"}}), nil)
        assertErr(t, Notify(Payload{Label: Label{"test": "test"}}), nil)
        assertErr(t, Notify(Payload{Label: Label{"test": "panic"}}), nil)
        assertErr(t, Close(), nil)
    }
}
