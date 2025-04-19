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

	defaultTimeout = time.Second * 5

	ErrorNodeNotFound = errors.New("parent node not found")
	ErrorNodeExists   = errors.New("node already exists")
	ErrorForceStop    = errors.New("shutdown force stopped")
	ErrorTimeout      = errors.New("shutdown timeout stopped")
)

type (
	// CallbackFunc to Add
	CallbackFunc func(ctx context.Context)

	Option func(*Shutdown)

	// Shutdown keeps track of all of your application Close/Shutdown dependencies
	Shutdown struct {
		mu sync.Mutex

		once    sync.Once
		started chan struct{}
		done    chan struct{}

		timeout time.Duration

		dependencyMap map[string]*Node
	}
)

// WithTimeout set shutdown timeout
func WithTimeout(duration time.Duration) Option {
	return func(s *Shutdown) {
		s.timeout = duration
	}
}

// New Shutdown constructor, default timeout 5s
func New(opts ...Option) *Shutdown {
	signalsCtx, cancel := signal.NotifyContext(context.Background(), signals...)

	shutdown := &Shutdown{
		mu:            sync.Mutex{},
		once:          sync.Once{},
		started:       make(chan struct{}),
		done:          make(chan struct{}),
		timeout:       defaultTimeout,
		dependencyMap: make(map[string]*Node),
	}

	for _, opt := range opts {
		opt(shutdown)
	}

	go func() {
		defer cancel()

		<-signalsCtx.Done()

		shutdown.Shutdown()
	}()

	return shutdown
}

// GetNodesNames returns all nodes names
func (s *Shutdown) GetNodesNames() []string {
	s.mu.Lock()
	defer s.mu.Unlock()

	result := make([]string, 0, len(s.dependencyMap))
	for k := range s.dependencyMap {
		result = append(result, k)
	}

	return result
}

// MustAdd adds a callback to a Shutdown, will panic if error
func (s *Shutdown) MustAdd(name string, callbackFunc CallbackFunc, parentNames ...string) {
	if err := s.Add(name, callbackFunc, parentNames...); err != nil {
		panic(err)
	}
}

// Add adds a callback to a Shutdown, can return ErrorNodeNotFound, ErrorNodeExists
func (s *Shutdown) Add(name string, callbackFunc CallbackFunc, parentNames ...string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.dependencyMap[name]; exists {
		return ErrorNodeExists
	}

	node := &Node{
		name:         name,
		parents:      make([]*Node, 0, len(parentNames)),
		children:     make([]*Node, 0),
		wg:           sync.WaitGroup{},
		callbackFunc: callbackFunc,
	}

	for _, parentName := range parentNames {
		parent, exists := s.dependencyMap[parentName]
		if !exists {
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

		ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
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
	case <-time.After(s.timeout):
		return ErrorTimeout
	case <-stop:
		return ErrorForceStop
	}
}

// shutdown running shutdown callbacks from parents to children concurrently
func (s *Shutdown) shutdown(ctx context.Context) {
	s.mu.Lock()
	defer s.mu.Unlock()

	wg := &sync.WaitGroup{}

	for _, node := range s.dependencyMap {
		if len(node.parents) == 0 {
			wg.Add(1)

			go func(node *Node) {
				defer wg.Done()

				node.shutdown(ctx)
			}(node)
		}
	}

	wg.Wait()

	s.dependencyMap = nil
}
