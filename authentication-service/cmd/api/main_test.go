package main

import (
	"os"
	"testing"

	"github.com/MikoBerries/go-micro-services/authentication-service/data"
)

var testApp Config

func TestMain(m *testing.M) {

	repo := data.NewPostgresTestRepository(nil)
	testApp.Repo = repo

	// TestMain entry point of testing in this package
	os.Exit(m.Run())
}
