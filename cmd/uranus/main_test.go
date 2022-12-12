package main

import "testing"

func TestSimple(t *testing.T) {
	if true != true {
		t.Errorf("Something's wrong")
	}
}
