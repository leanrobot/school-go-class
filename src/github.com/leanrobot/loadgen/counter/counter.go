/*
counter package provides a thread-safe singleton counter implementation.
The data is shared across the entire application, providing a simple counter
that can be used across a distributed application to keep track of statistics.
*/
package counter

import (
	"fmt"
)

var (
	commands chan *Action
	counts   map[string]int
)

type Action struct {
	key           string
	action        Command
	valueReceiver chan int
	copyReceiver  chan map[string]int
}

type Command int

const (
	getCmd = iota
	incrementCmd
	resetCmd
	exportCmd
	clearCmd
)

const DEF_CHAN_CAP = 1000

func init() {
	commands = make(chan *Action, DEF_CHAN_CAP)
	counts = make(map[string]int)
	go semaphore()

	fmt.Println(getCmd)
	fmt.Println(incrementCmd)
	fmt.Println(resetCmd)
	fmt.Println(exportCmd)
	fmt.Println(clearCmd)

}

func newAction(key string, cmd Command) *Action {
	return &Action{
		key:           key,
		action:        cmd,
		valueReceiver: make(chan int),
		copyReceiver:  make(chan map[string]int),
	}
}

// Increment adds one to the value of the specified key.
func Increment(key string) {
	incr := newAction(key, incrementCmd)
	commands <- incr
}

func Get(key string) int {
	get := newAction(key, getCmd)
	commands <- get
	return <-get.valueReceiver
}

// Reset changes the value of the specified key to zero.
func Reset(key string) {
	reset := newAction(key, resetCmd)
	commands <- reset
}

// Clear sets all data for every key to zero.
func Clear() {
	clear := newAction("", clearCmd)
	commands <- clear
}

// Export returns a copy of the counter data map.
func Export() map[string]int {
	export := newAction("", exportCmd)
	commands <- export
	return <-export.copyReceiver
}

/*
semphore is a goroutine who implements the access to the data store.
*/
func semaphore() {
	for {
		cmd := <-commands
		fmt.Printf("%v\n", cmd)

		switch cmd.action {
		case getCmd:
			cmd.valueReceiver <- counts[cmd.key]
		case incrementCmd:
			counts[cmd.key]++
		case resetCmd:
			counts[cmd.key] = 0
		case clearCmd:
			clearData()
		case exportCmd:
			cmd.copyReceiver <- copy()
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
