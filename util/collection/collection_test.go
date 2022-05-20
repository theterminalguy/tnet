package collection

import "testing"

// smoke test to validate docker setup
func TestAddition(t *testing.T) {
	got := 2 + 2
	expected := 4
	if got != expected {
		t.Errorf("Got: '%v', wanted: '%v'", got, expected)
	}
}

func TestContains_ItemExists(t *testing.T) {
	slice := []string{"a", "b", "c"}
	if !Contains(slice, "b") {
		t.Errorf("Expected 'b' to be in slice")
	}
}

func TestContains_ItemDoesNotExist(t *testing.T) {
	slice := []string{"a", "b", "c"}
	if Contains(slice, "d") {
		t.Errorf("Expected 'd' to not be in slice")
	}
}
