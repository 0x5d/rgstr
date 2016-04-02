package runtimes_test

import (
	. "github.com/castillobg/rgstr/runtimes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("runtimes", func() {
	Describe(".Register()", func() {
		Context("When the name is unique", func() {
			It("Registers the new AdapterFactory", func() {
				var factory AdapterFactory
				name := "factory1"
				Expect(Register(factory, name)).To(Succeed())
				// Deregister it to avoid conflicts with other tests.
				Deregister(name)
			})
		})
		Context("When the name is a dup", func() {
			It("Returns an error", func() {
				var factory AdapterFactory
				name := "factory1"
				Expect(Register(factory, name)).To(Succeed())
				// Trying to register the same runtime twice results in an error.
				Expect(Register(factory, name)).To(HaveOccurred())
				Deregister(name)
			})
		})
	})

	Describe(".Deregister()", func() {
		Context("When a given AdapterFactory exists", func() {
			It("Returns true.", func() {
				var factory AdapterFactory
				name := "factory1"
				Expect(Register(factory, name)).To(Succeed())
				Expect(Deregister(name)).To(BeTrue())
			})
		})
		Context("When a given AdapterFactory doesn't exist", func() {
			It("Returns false.", func() {
				Expect(Deregister("inexistent")).To(BeFalse())
			})
		})
	})

	Describe(".LookUp()", func() {
		Context("When a given service exists", func() {
			It("Returns the service and true", func() {
				var factory AdapterFactory
				name := "factory1"
				Expect(Register(factory, name)).To(Succeed())
				_, ok := LookUp(name)
				Expect(ok).To(BeTrue())
				// Deregister it to avoid conflicts with other tests.
				Deregister(name)
			})
		})
		Context("When a given service doesn't exist", func() {
			It("Returns the nil and false", func() {
				_, ok := LookUp("inexistent")
				Expect(ok).To(BeFalse())
			})
		})
	})
})
