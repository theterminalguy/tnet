package collection_test

import "testing"

// smoke test to validate docker setup
func TestAddition(t *testing.T) {
	got := 2 + 2
	expected := 4
	if got != expected {
		t.Errorf("Got: '%v', wanted: '%v'", got, expected)
	}
}
