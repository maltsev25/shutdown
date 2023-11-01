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

		InitGlobal()

		t.Run("timeout", func(t *testing.T) {
			t.Run("register", func(t *testing.T) {

				RegisterTimeout(time.Second)

				t.Run("get", func(t *testing.T) {

					if Timeout() != time.Second {
						t.Error("timeout not set")
					}
				})
			})
		})

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
		})
	})

}
