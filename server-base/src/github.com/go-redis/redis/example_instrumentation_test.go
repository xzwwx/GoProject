package redis_test

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type redisHook struct{}

var _ Hook = redisHook{}

func (redisHook) BeforeProcess(ctx context.Context, cmd Cmder) (context.Context, error) {
	fmt.Printf("starting processing: <%s>\n", cmd)
	return ctx, nil
}

func (redisHook) AfterProcess(ctx context.Context, cmd Cmder) error {
	fmt.Printf("finished processing: <%s>\n", cmd)
	return nil
}

func (redisHook) BeforeProcessPipeline(ctx context.Context, cmds []Cmder) (context.Context, error) {
	fmt.Printf("pipeline starting processing: %v\n", cmds)
	return ctx, nil
}

func (redisHook) AfterProcessPipeline(ctx context.Context, cmds []Cmder) error {
	fmt.Printf("pipeline finished processing: %v\n", cmds)
	return nil
}

func Example_instrumentation() {
	rdb := NewClient(&Options{
		Addr: ":6379",
	})
	rdb.AddHook(redisHook{})

	rdb.Ping(ctx)
	// Output: starting processing: <ping: >
	// finished processing: <ping: PONG>
}

func ExamplePipeline_instrumentation() {
	rdb := NewClient(&Options{
		Addr: ":6379",
	})
	rdb.AddHook(redisHook{})

	rdb.Pipelined(ctx, func(pipe Pipeliner) error {
		pipe.Ping(ctx)
		pipe.Ping(ctx)
		return nil
	})
	// Output: pipeline starting processing: [ping:  ping: ]
	// pipeline finished processing: [ping: PONG ping: PONG]
}

func ExampleClient_Watch_instrumentation() {
	rdb := NewClient(&Options{
		Addr: ":6379",
	})
	rdb.AddHook(redisHook{})

	rdb.Watch(ctx, func(tx *Tx) error {
		tx.Ping(ctx)
		tx.Ping(ctx)
		return nil
	}, "foo")
	// Output:
	// starting processing: <watch foo: >
	// finished processing: <watch foo: OK>
	// starting processing: <ping: >
	// finished processing: <ping: PONG>
	// starting processing: <ping: >
	// finished processing: <ping: PONG>
	// starting processing: <unwatch: >
	// finished processing: <unwatch: OK>
}
