package counter

import (
	tst "testing"
)

var expectedCounts map[string]int

func init() {
	expectedCounts = map[string]int{
		"key1":    1734,
		"bar":     118,
		"ff8338d": 2948,
	}
}

func TestSimpleCount(t *tst.T) {
	Clear()
	for key, actualCount := range expectedCounts {
		for i := 0; i < actualCount; i++ {
			Increment(key)
		}
		count := Get(key)
		if count != actualCount {
			t.Errorf("for counts[%s], got %d, expected %d",
				key, count, actualCount)
		}
	}
}

func TestConcurrentCount(t *tst.T) {
	Clear()
	controller := make(chan bool)
	incrementer := func(key string, times int) {
		for i := 0; i < times; i++ {
			i = i
			Increment(key)
		}
		controller <- true
	}

	// create an incrementer
	for key, actualCount := range expectedCounts {
		go incrementer(key, actualCount)
	}

	numWorkers := len(expectedCounts)
	// wait for workers
	for i := 0; i < numWorkers; i++ {
		<-controller
	}

	// check counts
	data := Export()
	for key, actualCount := range expectedCounts {
		if data[key] != actualCount {
			t.Errorf("for count %s, got %d, expected %d",
				key, data[key], actualCount)
		}
	}
}

// should be 5.
func TestReset(t *tst.T) {
	key := "foobar"
	const count = 5
	for i := 0; i < count; i++ {
		Increment(key)
	}
	Reset(key)
	for i := 0; i < count; i++ {
		Increment(key)
	}
	if Get(key) != count {
		t.Errorf("Reset incorrect. got %d, expected %d", count, Get(key))
	}
}
