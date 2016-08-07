package main

import "sync"

// utility class used to count results routine safe
type opsCounter struct {
	sync.RWMutex
	Sends       int
	TotalOps    int
	SendsToStop int
}

func (c *opsCounter) Done() bool {
	c.RLock()
	defer c.RUnlock()
	return c.Sends >= c.SendsToStop
}

func (c *opsCounter) Process(err error) {
	c.Lock()
	defer c.Unlock()
	c.TotalOps++
	if nil == err {
		c.Sends++
	}
}

func (c *opsCounter) GetSends() (int, int) {
	c.RLock()
	defer c.RUnlock()
	return c.Sends, c.TotalOps
}
