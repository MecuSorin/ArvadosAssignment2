package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestArvados(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Arvados Suite")
}
