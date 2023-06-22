package main

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {

	// TestMain entry point of testing in this package
	os.Exit(m.Run())
}
