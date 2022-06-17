package gonal

import (
    "context"
    "encoding/json"
    "fmt"
    "github.com/czasg/go-queue"
    "runtime"
    "sync"
    "time"
)

var gonal = NewGonal(context.Background())

func NewGonal(ctx context.Context, q ...queue.Queue) *Gonal {
    if len(q) < 1 {
        q = append(q, queue.NewFifoMemoryQueue())
    }
    gonal := Gonal{
        queue:    q[0],
        handlers: map[string][]HandlerWrap{},
    }
    gonal.ctx, gonal.cancel = context.WithCancel(ctx)
    gonal.SetConcurrent(runtime.NumCPU() * 2)
    return &gonal
}

func Notify(ctx context.Context, labels Labels, data []byte) error {
    return gonal.Notify(ctx, labels, data)
}

func BindHandler(labels Labels, handlers ...Handler) {
    gonal.BindHandler(labels, handlers...)
}

func FetchHandler(labels Labels) []Handler {
    return gonal.FetchHandler(labels)
}

func SetConcurrent(concurrent int) {
    gonal.SetConcurrent(concurrent)
}

func UnsafeSetGonal(g *Gonal) {
    gonal = g
}

func Close() error {
    return gonal.Close()
}

type Placeholder struct{}
type Labels map[string]string
type Handler func(ctx context.Context, labels Labels, data []byte)

type Gonal struct {
    ctx      context.Context
    cancel   context.CancelFunc
    reset    context.CancelFunc
    handlers map[string][]HandlerWrap
    queue    queue.Queue
    lock     sync.Mutex
    index    int
}

type Payload struct {
    Labels Labels
    Data   []byte
}

type HandlerWrap struct {
    Index   int
    Handler Handler
}

func (g *Gonal) Notify(ctx context.Context, labels Labels, data []byte) error {
    select {
    case <-g.ctx.Done():
        return g.ctx.Err()
    default:
    }
    body, err := json.Marshal(Payload{Labels: labels, Data: data})
    if err != nil {
        return err
    }
    return g.queue.Put(ctx, body)
}

func (g *Gonal) BindHandler(labels Labels, handlers ...Handler) {
    g.lock.Lock()
    defer g.lock.Unlock()
    handlerWraps := []HandlerWrap{}
    for _, handler := range handlers {
        handlerWraps = append(handlerWraps, HandlerWrap{
            Index:   g.index,
            Handler: handler,
        })
        g.index++
    }
    for key, value := range labels {
        hk := fmt.Sprintf("%s=%s", key, value)
        g.handlers[hk] = append(g.handlers[hk], handlerWraps...)
    }
}

func (g *Gonal) FetchHandler(labels Labels) []Handler {
    g.lock.Lock()
    defer g.lock.Unlock()
    results := []Handler{}
    handlerSet := map[int]Placeholder{}
    for key, value := range labels {
        hk := fmt.Sprintf("%s=%s", key, value)
        handlers, ok := g.handlers[hk]
        if !ok {
            continue
        }
        for _, handler := range handlers {
            _, ok := handlerSet[handler.Index]
            if ok {
                continue
            }
            handlerSet[handler.Index] = Placeholder{}
            results = append(results, handler.Handler)
        }
    }
    return results
}

func (g *Gonal) Close() error {
    if g.cancel != nil {
        g.cancel()
    }
    return g.queue.Close()
}

func (g *Gonal) SetConcurrent(concurrent int) {
    g.lock.Lock()
    defer g.lock.Unlock()
    if g.reset != nil {
        g.reset()
        time.Sleep(time.Millisecond * 100) // 预留线程退出时间
    }
    ctx, cancel := context.WithCancel(g.ctx)
    limit := make(chan Placeholder, concurrent)
    g.reset = cancel
    go g.loop(ctx, limit)
}

func (g *Gonal) loop(ctx context.Context, limit chan Placeholder) {
    defer close(limit)
    for {
        select {
        case <-ctx.Done():
            return
        default:
        }
        data, err := g.queue.Get(ctx)
        if err != nil {
            continue
        }
        var payload Payload
        err = json.Unmarshal(data, &payload)
        if err != nil {
            continue
        }
        for _, handler := range g.FetchHandler(payload.Labels) {
            select {
            case <-ctx.Done():
                return
            case limit <- Placeholder{}:
            }
            go func(handler Handler) {
                defer func() {
                    if err := recover(); err != nil {
                    }
                    <-limit
                }()
                handler(ctx, payload.Labels, payload.Data)
            }(handler)
        }
    }
}
