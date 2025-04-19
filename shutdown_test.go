package shutdown

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

func TestShutdown(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Parallel()

		t.Run("simple registration", func(t *testing.T) {
			shutdown := New()

			shutdown.MustAdd("node1", func(ctx context.Context) {
				t.Log("node1 shutdown success")
			})

			shutdown.MustAdd("node2", func(ctx context.Context) {
				t.Log("node2 shutdown success")
			}, "node1")

			shutdown.MustAdd("node3", func(ctx context.Context) {
				t.Log("node3 shutdown success")
			}, "node1", "node2")

			shutdown.MustAdd("node4", func(ctx context.Context) {
				t.Log("node4 shutdown success")
			}, "node2", "node3")

			shutdown.MustAdd("node5", func(ctx context.Context) {
				t.Log("node5 shutdown success")
			})

			t.Run("get nodes names", func(t *testing.T) {
				nodes := shutdown.GetNodesNames()
				if len(nodes) != 5 {
					t.Errorf("got %d, want 5", len(nodes))
				}
			})

			shutdown.Shutdown()
		})

		t.Run("concurrent registration", func(t *testing.T) {
			shutdown := New()

			shutdown.MustAdd("node11", func(ctx context.Context) {
				t.Log("node11 shutdown success")
			})

			wg := sync.WaitGroup{}
			wg.Add(3)

			go func() {
				shutdown.MustAdd("node12", func(ctx context.Context) {
					t.Log("node12 shutdown success")
				}, "node11")
				wg.Done()
			}()

			go func() {
				shutdown.MustAdd("node13", func(ctx context.Context) {
					t.Log("node13 shutdown success")
				}, "node11")
				wg.Done()
			}()

			go func() {
				shutdown.MustAdd("node14", func(ctx context.Context) {
					t.Log("node14 shutdown success")
				}, "node11")
				wg.Done()
			}()

			wg.Wait()

			shutdown.Shutdown()
		})
	})

	t.Run("error node already exists", func(t *testing.T) {
		t.Parallel()

		shutdown := New()

		shutdown.MustAdd("node1", func(ctx context.Context) {
			t.Log("node1 shutdown success")
		})

		err := shutdown.Add("node1", func(ctx context.Context) {
			t.Log("node1 shutdown success")
		})

		if !errors.Is(err, ErrorNodeExists) {
			t.Errorf("got %v, want %v", err, ErrorNodeExists)
		}
	})

	t.Run("error node not found", func(t *testing.T) {
		t.Parallel()

		shutdown := New()

		err := shutdown.Add("node2", func(ctx context.Context) {
			t.Log("node2 shutdown success")
		}, "node1")

		if !errors.Is(err, ErrorNodeNotFound) {
			t.Errorf("got %v, want %v", err, ErrorNodeNotFound)
		}
	})

	t.Run("shutdown wait", func(t *testing.T) {
		t.Parallel()

		shutdown := New()

		shutdown.MustAdd("node1", func(ctx context.Context) {
			time.Sleep(time.Millisecond * 100)
			t.Log("node1 shutdown success")
		})

		shutdown.MustAdd("node2", func(ctx context.Context) {
			time.Sleep(time.Millisecond * 100)
			t.Log("node2 shutdown success")
		}, "node1")

		go func() {
			// syscall.Kill(syscall.Getpid(), syscall.SIGINT)
			shutdown.Shutdown()
		}()

		err := shutdown.Wait()
		if err != nil {
			t.Errorf("got %v, want nil", err)
		}
	})

	t.Run("shutdown twice", func(t *testing.T) {
		t.Parallel()

		shutdown := New()

		shutdown.MustAdd("node1", func(ctx context.Context) {
			time.Sleep(time.Millisecond * 100)
			t.Log("node1 shutdown success")
		})

		shutdown.MustAdd("node2", func(ctx context.Context) {
			time.Sleep(time.Millisecond * 100)
			t.Log("node2 shutdown success")
		}, "node1")

		go func() {
			shutdown.Shutdown()
			shutdown.Shutdown()
		}()

		err := shutdown.Wait()
		if err != nil {
			t.Errorf("got %v, want nil", err)
		}
	})

	t.Run("shutdown with timeout", func(t *testing.T) {
		t.Parallel()

		shutdown := New(WithTimeout(time.Millisecond * 100))

		shutdown.MustAdd("node1", func(ctx context.Context) {
			time.Sleep(time.Millisecond * 100)
			t.Log("node1 shutdown success")
		})

		shutdown.MustAdd("node2", func(ctx context.Context) {
			time.Sleep(time.Millisecond * 100)
			t.Log("node2 shutdown success")
		}, "node1")

		go func() {
			shutdown.Shutdown()
		}()

		err := shutdown.Wait()
		if !errors.Is(err, ErrorTimeout) {
			t.Errorf("expected ErrorTimeout, got %v", err)
		}
	})
}
