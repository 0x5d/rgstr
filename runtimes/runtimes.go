package runtimes

import (
	"errors"

	"github.com/castillobg/rgstr/registries"
)

var registeredFactories = make(map[string]AdapterFactory)

// AdapterFactory specifies a constructor for RuntimeAdapter factories.
type AdapterFactory interface {
	New(address string, registry *registries.RegistryAdapter) (RuntimeAdapter, error)
}

// RuntimeAdapter specifies the contract a container runtime adapter should follow.
type RuntimeAdapter interface {
	Listen(errs chan error)
}

// Register registers an AdapterFactory for use.
func Register(rf AdapterFactory, name string) error {
	if _, ok := registeredFactories[name]; ok {
		// Should be unique
		return errors.New("A runtime with the name \"" + name + "\" was already registered.")
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
