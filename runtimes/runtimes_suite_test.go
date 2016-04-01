package runtimes_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestRuntimes(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Runtimes Suite")
}
