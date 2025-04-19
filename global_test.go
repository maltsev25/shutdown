package shutdown

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestGlobal(t *testing.T) {
	t.Run("init", func(t *testing.T) {
		t.Parallel()

		InitGlobal(WithTimeout(time.Second))

		t.Run("add", func(t *testing.T) {
			t.Run("success", func(t *testing.T) {
				MustAdd("node1", func(ctx context.Context) {
					t.Log("node1 shutdown success")
				})

				err := Add("node2", func(ctx context.Context) {
					t.Log("node2 shutdown success")
				})
				if err != nil {
					t.Error(err)
				}

				t.Run("get nodes names", func(t *testing.T) {
					nodes := GetNodesNames()
					if len(nodes) != 2 {
						t.Errorf("got %d, want 2", len(nodes))
					}
				})
			})
			t.Run("error", func(t *testing.T) {
				MustAdd("node3", func(ctx context.Context) {
					t.Log("node3 shutdown success")
				})

				err := Add("node3", func(ctx context.Context) {
					t.Log("node3 shutdown success")
				})
				if !errors.Is(err, ErrorNodeExists) {
					t.Errorf("got %v, want %v", err, ErrorNodeExists)
				}
			})
			t.Run("parent not found", func(t *testing.T) {
				err := Add("node4", func(ctx context.Context) {
					t.Log("node4 shutdown success")
				}, "node5")
				if !errors.Is(err, ErrorNodeNotFound) {
					t.Errorf("got %v, want %v", err, ErrorNodeNotFound)
				}
			})
			t.Run("shutdown twice", func(t *testing.T) {
				ForceShutdown()
				ForceShutdown()
				err := Wait()
				if err != nil {
					t.Error(err)
				}
			})
		})
	})
}
