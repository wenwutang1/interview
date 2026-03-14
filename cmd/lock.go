package main

import "sync"

type MyLock struct {
	// ttl  int
	key  string
	lock sync.Mutex
}
