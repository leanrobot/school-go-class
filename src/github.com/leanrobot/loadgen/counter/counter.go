package counter

import ()

var (
	counts    map[string]int
	increment chan string
	reset     chan int
	clear     chan bool
	export    chan map[string]int
)

const DEF_CHAN_CAP = 10

func init() {
	counts = make(map[string]int)
	increment = make(chan string, DEF_CHAN_CAP)
	reset = make(chan int, DEF_CHAN_CAP)
	clear = make(chan bool, DEF_CHAN_CAP)
	export = make(chan chan map[string]int, DEF_CHAN_CAP)
}

func Increment(key string) {
	increment <- key
}

func Reset() {
	reset <- key
}

func Clear() {
	clear <- true
}

func Export() {
	listener = make(chan map[string]int)
	export <- listener
	return <-listener
}

func semaphore() {
	for {
		select {
		case key := <-increment:
			counts[key]++
		case key := <-reset:
			counts[key] = 0
		case key := <-clear:
			clearData()
		case exportChan := <-export:
			exportChan <- copy()
		}
	}
}

func clearData() {
	for key, _ := range counts {
		delete(counts, key)
	}
}

func copy() map[string]int {
	copy := new(map[string]int)
	for key, value := range counts {
		copy[key] = value
	}
	return copy
}
