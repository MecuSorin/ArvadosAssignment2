package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("The bytes generator", func() {

	Specify("Should accept only a valid size", func() {
		err := validateBlobSize(0)
		Expect(err).Should(HaveOccurred())
		err = validateBlobSize(1 + 64*MIB)
		Expect(err).Should(HaveOccurred())
		err = validateBlobSize(1)
		Expect(err).Should(Succeed())
	})
})
