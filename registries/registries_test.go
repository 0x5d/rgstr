package registries_test

import (
	"testing"

	"github.com/castillobg/rgstr/registries"
	"github.com/stretchr/testify/assert"
)

func TestRegisterUnique(t *testing.T) {
	var factory registries.AdapterFactory
	name := "factory1"

	// When a factory's name is unique, there shouldn't be a problem.
	err := registries.Register(factory, name)
	// To avoid conflicts with other tests, deregister it.
	defer registries.Deregister(name)
	assert.NoError(t, err, "err should be nil if the factory name is unique. Instead got: %v", err)
}

func TestRegisterDup(t *testing.T) {
	var factory registries.AdapterFactory
	name := "factory1"

	// When a factory's name is unique, there shouldn't be a problem.
	err := registries.Register(factory, name)
	// To avoid conflicts with other tests, deregister it.
	defer registries.Deregister(name)
	assert.NoError(t, err, "Unexpected error: %v", err)
	// When a factory's name is a duplicate, client should get an error.
	err = registries.Register(factory, name)
	assert.NotNil(t, err, "err shouldn't be nil when registering a duplicate factory.")
}

func TestDeregisterExistent(t *testing.T) {
	var factory registries.AdapterFactory
	name := "factory1"

	// The factory's name is unique, there shouldn't be a problem.
	err := registries.Register(factory, name)
	assert.NoError(t, err, "Unexpected error: %v", err)
	ok := registries.Deregister(name)
	assert.True(t, ok, "ok shouldn't be false if the factory existed.")
}

func TestDeregisterInexistent(t *testing.T) {
	name := "factory1"
	assert.False(t, registries.Deregister(name), "Deregister shouldn't return false if the factory existed.")
}

func TestLookUpExistent(t *testing.T) {
	var factory registries.AdapterFactory
	name := "factory1"

	// The factory's name is unique, there shouldn't be a problem.
	err := registries.Register(factory, name)
	defer registries.Deregister(name)
	assert.NoError(t, err, "Unexpected error: %v", err)
	_, ok := registries.LookUp(name)
	assert.True(t, ok, "ok should be true when a factory exists.")
}

func TestLookUpInexistent(t *testing.T) {
	_, ok := registries.LookUp("inexistent")
	assert.False(t, ok, "ok should be false when a factory doesn't exist.")
}
