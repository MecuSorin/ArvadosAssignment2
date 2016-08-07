package main

import (
	"encoding/base64"
	"io/ioutil"

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

	Specify("Should generate exactly the requested size", func() {
		r, err := getBlob(10)
		Expect(err).Should(Succeed())
		rb64 := base64.NewDecoder(base64.StdEncoding, r)
		data, err := ioutil.ReadAll(rb64)

		Expect(err).Should(Succeed())

		Expect(len(data)).Should(Equal(10))
	})
})
