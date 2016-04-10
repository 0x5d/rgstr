package registries_test

import (
	"testing"

	"github.com/castillobg/rgstr/registries"
)

func TestRegisterUnique(t *testing.T) {
	var factory registries.AdapterFactory
	name := "factory1"

	// When a factory's name is unique, there shouldn't be a problem.
	err := registries.Register(factory, name)
	// To avoid conflicts with other tests, deregister it.
	defer registries.Deregister(name)
	if err != nil {
		t.Errorf("err should be nil if the factory name is unique. Instead got: %v", err)
	}
}

func TestRegisterDup(t *testing.T) {
	var factory registries.AdapterFactory
	name := "factory1"

	// When a factory's name is unique, there shouldn't be a problem.
	err := registries.Register(factory, name)
	// To avoid conflicts with other tests, deregister it.
	defer registries.Deregister(name)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	// When a factory's name is a duplicate, client should get an error.
	err = registries.Register(factory, name)
	if err == nil {
		t.Error("err shouldn't be nil when registering a duplicate factory.")
	}
}

func TestDeregisterExistent(t *testing.T) {
	var factory registries.AdapterFactory
	name := "factory1"

	// The factory's name is unique, there shouldn't be a problem.
	err := registries.Register(factory, name)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	ok := registries.Deregister(name)
	if !ok {
		t.Error("ok shouldn't be false if the factory existed.")
	}
}

func TestDeregisterInexistent(t *testing.T) {
	name := "factory1"
	ok := registries.Deregister(name)
	if ok {
		t.Error("ok should be false if the factory deidn't exist.")
	}
}

func TestLookUpExistent(t *testing.T) {
	var factory registries.AdapterFactory
	name := "factory1"

	// The factory's name is unique, there shouldn't be a problem.
	err := registries.Register(factory, name)
	defer registries.Deregister(name)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	_, ok := registries.LookUp(name)
	if !ok {
		t.Error("ok should be true when a factory exists.")
	}
}

func TestLookUpInexistent(t *testing.T) {
	_, ok := registries.LookUp("inexistent")
	if ok {
		t.Error("ok should be false when a factory doesn't exist.")
	}
}
