package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestRgstr(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "rgstr suite")
}
