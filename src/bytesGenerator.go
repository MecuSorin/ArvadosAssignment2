package main

import (
	"encoding/base64"
	"errors"
	"io"
	"math/rand"
)

// MIB is a standard measure
const MIB = 1048576

func getBlob(size int) (io.Reader, error) {
	if nil != validateBlobSize(size) {
		return nil, errors.New("Invalid size for the blob")
	}

	pipeReader, pipeWriter := io.Pipe()
	wb64 := base64.NewEncoder(base64.StdEncoding, pipeWriter)
	go func() {
		defer pipeWriter.Close()
		defer wb64.Close()

		for i := 0; i < size; i++ {
			wb64.Write([]byte{byte(rand.Intn(256))})
		}
	}()
	return pipeReader, nil
}

func validateBlobSize(size int) error {
	if size > 64*MIB || size < 1 {
		return errors.New("Invalid size for the blob")
	}
	return nil
}
