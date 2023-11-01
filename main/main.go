package main

import (
	"context"
	"time"

	"github.com/maltsev25/shutdown"
)

func main() {
	gs := shutdown.New()

	// ctx := context.Background()

	gs.MustAdd("node1", func(ctx context.Context) {
		time.Sleep(time.Millisecond * 100)
		println("node1")
	})

	gs.MustAdd("node10", func(ctx context.Context) {
		time.Sleep(time.Millisecond * 100)
		println("node10")
	})

	gs.MustAdd("node2", func(ctx context.Context) {
		time.Sleep(time.Millisecond * 100)
		println("node2")
	}, "node1")

	gs.MustAdd("node3", func(ctx context.Context) {
		time.Sleep(time.Millisecond * 100)
		println("node3")
	}, "node1", "node2")

	gs.MustAdd("node4", func(ctx context.Context) {
		time.Sleep(time.Millisecond * 100)
		println("node4")
	}, "node2", "node3")

	gs.MustAdd("node5", func(ctx context.Context) {
		time.Sleep(time.Millisecond * 100)
		println("node5")
	}, "node1")

	gs.MustAdd("node7", func(ctx context.Context) {
		time.Sleep(time.Millisecond * 100)
		println("node7")
	}, "node1", "node10")

	gs.MustAdd("node6", func(ctx context.Context) {
		time.Sleep(time.Millisecond * 100)
		println("node6")
	}, "node4", "node5", "node10")

	if err := gs.Wait(); err != nil {
		panic(err)
	}
}
