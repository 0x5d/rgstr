package bridge_test

import (
	. "github.com/castillobg/rgstr/bridge"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("bridge", func() {
	Describe(".Register()", func() {
		Context("When the name is unique", func() {
			It("Registers the new AdapterFactory", func() {

				var factory AdapterFactory
				Register(factory, "factory1")
				_, ok := LookUp("factory1")
				Expect(ok).To(BeTrue())
			})
		})
	})
})
