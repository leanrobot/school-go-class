/*
counter package provides a thread-safe singleton counter implementation.
The data is shared across the entire application, providing a simple counter
that can be used across a distributed application to keep track of statistics.
*/
package counter

var (
	counts    map[string]int
	increment chan string
	reset     chan string
	clear     chan bool
	export    chan chan map[string]int
)

const DEF_CHAN_CAP = 10

func init() {
	counts = make(map[string]int)
	increment = make(chan string, DEF_CHAN_CAP)
	reset = make(chan string, DEF_CHAN_CAP)
	clear = make(chan bool, DEF_CHAN_CAP)
	export = make(chan chan map[string]int, DEF_CHAN_CAP)

	go semaphore()
}

// Increment adds one to the value of the specified key.
func Increment(key string) {
	increment <- key
}

// Reset changes the value of the specified key to zero.
func Reset(key string) {
	reset <- key
}

// Clear sets all data for every key to zero.
func Clear() {
	clear <- true
}

// Export returns a copy of the counter data map.
func Export() map[string]int {
	listener := make(chan map[string]int)
	export <- listener
	return <-listener
}

/*
semphore is a goroutine who implements the access to the data store.
A channel exists for every action that the counter supports. when data is
received on the channel the action is performed on the datastore.
*/
func semaphore() {
	for {
		select {
		// increment a key
		case key := <-increment:
			counts[key]++
		// reset a key
		case key := <-reset:
			counts[key] = 0
		// clear all counter data
		case <-clear:
			clearData()
		// more complicated, a channel is sent through export. copy retuns
		// a copy of the data map, which is then returned through exportChan
		// to the caller.
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
	dataCopy := make(map[string]int)
	for key, value := range counts {
		dataCopy[key] = value
	}
	return dataCopy
}
