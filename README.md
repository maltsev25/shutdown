shutdown
===========

The `shutdown` package is package that provides a simple way to shutdown your app.

Example Use
===========
```go
shutdown := shutdown.New()

shutdown.MustAdd("node1", func(ctx context.Context) {
    t.Log("node1 shutdown success")
})

shutdown.MustAdd("node2", func(ctx context.Context) {
    t.Log("node2 shutdown success")
}, "node1")

shutdown.Wait()
```