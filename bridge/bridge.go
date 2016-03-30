package bridge

import "errors"

var registeredFactories = make(map[string]RegistryFactory)

// RegistryFactory specifies a constructor for factories.
type RegistryFactory interface {
	New() (RegistryAdapter, error)
}

// RegistryAdapter specifies the contract a container runtime adapter (docker, rkt) should follow.
type RegistryAdapter interface {
	Register(service *Service) error
	Deregister(service *Service) error
}

// Service represents a service.
type Service struct {
	ID   string
	Name string
	IP   string
	Port int
}

// Register registers a new RegistryFactory for use.
func Register(rf RegistryFactory, name string) error {
	if _, ok := registeredFactories[name]; ok {
		// Should be unique (either "consul", "etcd", etc.)
		return errors.New("A registry with the name \"" + name + "\" was already registered.")
	}
	registeredFactories[name] = rf
	return nil
}

// LookUp returns a RegistryFactory registered with a given name.
func LookUp(name string) (RegistryFactory, bool) {
	runtime, ok := registeredFactories[name]
	return runtime, ok
}
