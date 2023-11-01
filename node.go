package shutdown

import (
	"context"
	"sync"
)

type (
	Node struct {
		name     string
		parents  []*Node
		children []*Node

		wg           sync.WaitGroup
		callbackFunc CallbackFunc
		execution    sync.Once
	}
)

func (n *Node) shutdown(ctx context.Context) {
	for _, child := range n.children {
		go func(child *Node) {
			child.shutdown(ctx)
		}(child)
	}

	n.wg.Wait()

	n.execution.Do(func() {
		n.callbackFunc(ctx)

		for _, parent := range n.parents {
			parent.wg.Done()
		}
	})
}
