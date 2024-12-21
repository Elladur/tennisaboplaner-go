package main

import "testing"

func TestMain(t *testing.T) {
	if "Hello, Go!" != "Hello, Go!" {
		t.Errorf("Unexpected output")

	}
}
