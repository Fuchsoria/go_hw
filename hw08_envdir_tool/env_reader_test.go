package main

import (
	"testing"
)

func TestReadDir(t *testing.T) {
	// Place your code here
	t.Run("Test envirovments", func(t *testing.T) {
		ReadDir("testdata/env")
	})
}
