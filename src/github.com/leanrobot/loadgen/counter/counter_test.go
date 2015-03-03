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
}
