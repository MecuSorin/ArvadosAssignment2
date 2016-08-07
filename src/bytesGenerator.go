package main

import (
	"encoding/base64"
	"errors"
	"io"
	"math/rand"
)

// MIB is a standard measure
const MIB = 1048576


func validateBlobSize(size int) error {
	if size > 64*MIB || size < 1 {
		return errors.New("Invalid size for the blob")
	}
	return nil
}
