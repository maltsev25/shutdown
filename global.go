package shutdown

var globalShutdown *Shutdown

func InitGlobal(opts ...Option) {
	globalShutdown = New(opts...)
}

// GetNodesNames returns all nodes names
func GetNodesNames() []string {
	return globalShutdown.GetNodesNames()
}

// MustAdd shutdown callback to a global Shutdown, will panic if error
func MustAdd(name string, callbackFunc CallbackFunc, parentNames ...string) {
	globalShutdown.MustAdd(name, callbackFunc, parentNames...)
}

// Add shutdown callback to a global Shutdown, can return ErrorNodeNotFound, ErrorNodeExists
func Add(name string, callbackFunc CallbackFunc, parentNames ...string) error {
	return globalShutdown.Add(name, callbackFunc, parentNames...)
}

// Wait for a global Shutdown, check Shutdown.Wait
func Wait() error {
	return globalShutdown.Wait()
}

// ForceShutdown for a global Shutdown, check Shutdown.Shutdown
func ForceShutdown() {
	globalShutdown.Shutdown()
}
