/*
Concurrent map contains an implementation of a key-value store which is
thread-safe. The underlying implementation is a map[string] string which is
locked with a sync.Mutex.
*/
package concurrentmap
