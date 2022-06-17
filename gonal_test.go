package gonal

import (
    "context"
    "testing"
)

func TestGonal(t *testing.T) {
    {
        UnsafeSetGonal(NewGonal(context.Background()))
    }
    {
        check := make(chan Placeholder)
        test1 := func(ctx context.Context, labels Labels, data []byte) {
            check <- Placeholder{}
        }
        BindHandler(map[string]string{"test": "1"}, test1)
        FetchHandler(map[string]string{"test": "1"})
        SetConcurrent(1)
        _ = Notify(nil, map[string]string{"test": "1"}, nil)
        <-check
        Close()
    }
    {
        check := make(chan Placeholder)
        test1 := func(ctx context.Context, labels Labels, data []byte) {
            check <- Placeholder{}
        }
        g := NewGonal(context.Background())
        g.SetConcurrent(1)
        g.BindHandler(map[string]string{"test": "1"}, test1, test1)
        _ = g.Notify(nil, map[string]string{"test": "2"}, nil)
        _ = g.Notify(nil, map[string]string{"test": "1"}, nil)
        <-check
        g.Close()
    }
    {
        check := make(chan Placeholder)
        test1 := func(ctx context.Context, labels Labels, data []byte) {
            check <- Placeholder{}
        }
        test2 := func(ctx context.Context, labels Labels, data []byte) {
            check <- Placeholder{}
        }
        g := NewGonal(context.Background())
        g.SetConcurrent(1)
        g.BindHandler(map[string]string{"test1": "1"}, test1)
        g.BindHandler(map[string]string{"test2": "2"}, test2)
        _ = g.Notify(nil, map[string]string{"test1": "1", "test2": "2"}, nil)
        <-check
        <-check
        g.Close()
    }
    {
        ctx, cancel := context.WithCancel(context.Background())
        cancel()
        g := NewGonal(ctx)
        if err := g.Notify(nil, nil, nil); err != ctx.Err() {
            t.Error("上下文cancel测试异常")
        }
    }
}
