package shutdown

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	signals = []os.Signal{syscall.SIGINT, syscall.SIGTERM}

	timeout = time.Second * 5

	ErrorNodeNotFound = errors.New("parent node not found")
	ErrorNodeExists   = errors.New("node already exists")
	ErrorForceStop    = errors.New("shutdown force stopped")
	ErrorTimeout      = errors.New("shutdown timeout stopped")
)

type (
	// CallbackFunc to Add
	CallbackFunc func(ctx context.Context)

	// Shutdown keeps track of all of your application Close/Shutdown dependencies
	Shutdown struct {
		mx sync.Mutex

		once    sync.Once
		started chan struct{}
		done    chan struct{}

		dependencyMap map[string]*Node
	}
)

// New GracefulShutdown constructor
func New() *Shutdown {
	osCTX, cancel := signal.NotifyContext(context.Background(), signals...)

	shutdown := &Shutdown{
		mx:            sync.Mutex{},
		once:          sync.Once{},
		started:       make(chan struct{}),
		done:          make(chan struct{}),
		dependencyMap: make(map[string]*Node),
	}

	go func() {
		defer cancel()

		<-osCTX.Done()

		shutdown.Shutdown()
	}()

	return shutdown
}

// GetNodesNames returns all nodes names
func (s *Shutdown) GetNodesNames() []string {
	s.mx.Lock()
	defer s.mx.Unlock()

	result := make([]string, 0, len(s.dependencyMap))
	for k := range s.dependencyMap {
		result = append(result, k)
	}

	return result
}

// MustAdd adds a callback to a Shutdown instance
func (s *Shutdown) MustAdd(name string, callbackFunc CallbackFunc, parents ...string) {
	if err := s.Add(name, callbackFunc, parents...); err != nil {
		panic(err)
	}
}

// Add adds a callback to a GracefulShutdown instance, can return ErrorNodeNotFound, ErrorNodeExists
func (s *Shutdown) Add(name string, callbackFunc CallbackFunc, parents ...string) error {
	s.mx.Lock()
	defer s.mx.Unlock()

	if _, ok := s.dependencyMap[name]; ok {
		return ErrorNodeExists
	}

	node := &Node{
		name:         name,
		parents:      []*Node{},
		children:     []*Node{},
		wg:           sync.WaitGroup{},
		callbackFunc: callbackFunc,
	}

	for _, parentName := range parents {
		parent, ok := s.dependencyMap[parentName]
		if !ok {
			return ErrorNodeNotFound
		}

		parent.wg.Add(1)

		node.parents = append(node.parents, parent)
		parent.children = append(parent.children, node)
	}

	s.dependencyMap[name] = node

	return nil
}

// Shutdown processes all shutdown callbacks concurrently in a limited time frame (Timeout)
func (s *Shutdown) Shutdown() {
	s.once.Do(func() {
		defer close(s.done)

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		close(s.started)

		s.shutdown(ctx)
	})
}

// Wait for shutdown initiated by OS signal, can be forced, cancelled by timeout, finished correctly.
func (s *Shutdown) Wait() error {
	<-s.started

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, signals...)

	select {
	case <-s.done:
		return nil
	case <-time.After(timeout):
		return ErrorTimeout
	case <-stop:
		return ErrorForceStop
	}
}

// shutdown running shutdown callbacks from parents to children concurrently
func (s *Shutdown) shutdown(ctx context.Context) {
	s.mx.Lock()
	defer s.mx.Unlock()

	waitGroup := &sync.WaitGroup{}

	for _, node := range s.dependencyMap {
		if len(node.parents) == 0 {
			waitGroup.Add(1)

			go func(node *Node) {
				defer waitGroup.Done()

				node.shutdown(ctx)
			}(node)
		}
	}

	waitGroup.Wait()

	s.dependencyMap = make(map[string]*Node)
}
