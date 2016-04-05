package registries

import "errors"

var registeredFactories = make(map[string]AdapterFactory)

// AdapterFactory specifies a constructor for RegistryAdapter factories.
type AdapterFactory interface {
	// New builds a RegistryAdapter, which should be a client of a registry listening on the given
	// address.
	New(address string) (RegistryAdapter, error)
}

// RegistryAdapter specifies the contract a registry adapter (consul, etcd) should follow.
type RegistryAdapter interface {
	Register(service *Service) error
	Deregister(service *Service) error
}

// Service represents a service.
type Service struct {
	ID   string
	Name string
	IP   string
	Port uint
}

// Register registers an AdapterFactory for use.
func Register(rf AdapterFactory, name string) error {
	if _, ok := registeredFactories[name]; ok {
		// Should be unique (either "consul", "etcd", etc.)
		return errors.New("A registry with the name \"" + name + "\" was already registered.")
	}
	registeredFactories[name] = rf
	return nil
}

// Deregister deregisters an existent factory. (Mostly here for testing.)
func Deregister(name string) bool {
	_, ok := registeredFactories[name]
	delete(registeredFactories, name)
	return ok
}

// LookUp returns an AdapterFactory registered with a given name.
func LookUp(name string) (AdapterFactory, bool) {
	registry, ok := registeredFactories[name]
	return registry, ok
}
