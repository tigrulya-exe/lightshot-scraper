package util

import (
	"os"
	"sync"
)

func CreateDirectoryIfNotExists(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		_ = os.MkdirAll(path, os.ModePerm)
	}
}

type Counter struct {
	value int
	lock  sync.Mutex
}

func NewCounter(value int) *Counter {
	return &Counter{
		value: value,
	}
}

func (c *Counter) Increment() {
	c.lock.Lock()
	c.value++
	c.lock.Unlock()
}

func (c *Counter) Get() int {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.value
}
