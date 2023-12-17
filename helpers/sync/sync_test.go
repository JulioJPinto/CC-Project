package sync

import (
	"testing"
	"reflect"
)

func TestRemoveIf(t *testing.T) {
	// Create a new set
	set := NewSyncSet[int]()
	set.Add(1)
	set.Add(2)
	set.Add(3)

	// Remove elements greater than 1
	set.RemoveIf(func(element int) bool {
		return element > 1
	})

	// Verify that the set now contains only elements not removed
	expected := []int{1}
	actual := set.List()

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Test failed: Expected %v, got %v", expected, actual)
	}
}

func TestAnyMatch(t *testing.T) {
	// Define a simple Set with an underlying Map
	set := Set[int]{} // replace with your actual Set initialization function

	// Add some elements to the set
	set.Add(1)
	set.Add(2)
	set.Add(3)

	// Test case 1: Predicate that matches an element in the set
	predicate1 := func(element int) bool {
		return element == 2
	}
	result1 := set.AnyMatch(predicate1)
	if !result1 {
		t.Errorf("Test case 1 failed: Expected true, got false")
	}

	// Test case 2: Predicate that does not match any element in the set
	predicate2 := func(element int) bool {
		return element == 4
	}
	result2 := set.AnyMatch(predicate2)
	if result2 {
		t.Errorf("Test case 2 failed: Expected false, got true")
	}
}
