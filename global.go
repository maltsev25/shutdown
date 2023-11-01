package shutdown

import (
	"time"
)

var globalShutdown *Shutdown

func InitGlobal() {
	globalShutdown = New()
}

// RegisterTimeout for a different shutdown timeout
func RegisterTimeout(duration time.Duration) {
	timeout = duration
}

// Timeout for shutdown
func Timeout() time.Duration {
	return timeout
}

// MustAdd shutdown callback to a global GracefulShutdown
func MustAdd(name string, callbackFunc CallbackFunc, parents ...string) {
	globalShutdown.MustAdd(name, callbackFunc, parents...)
}

// Add shutdown callback to a global Shutdown
func Add(name string, callbackFunc CallbackFunc, parents ...string) error {
	return globalShutdown.Add(name, callbackFunc, parents...)
}

// Wait for a global GracefulShutdown, check GracefulShutdown.Wait
func Wait() error {
	return globalShutdown.Wait()
}

// ForceShutdown for a global Shutdown, check Shutdown.Shutdown
func ForceShutdown() {
	globalShutdown.Shutdown()
}
