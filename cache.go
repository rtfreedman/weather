package weather

import (
	"sync"
	"time"
)

type cache struct {
	weathers map[int]Weather
	sync.RWMutex
}

var cachedWeather = cache{
	weathers: map[int]Weather{},
}

// retrieve will retrieve weather from the cache if present
func (c *cache) retrieve(zip int) (w Weather, present, expired bool) {
	// rlock to protect the map from simultaneous read/writes with the add function
	c.RLock()
	defer c.RUnlock()

	// check if it exists in the map
	w, present = c.weathers[zip]
	if !present {
		return
	}
	// determine whether we're expired
	expired = w.Expiry.Unix() < time.Now().Unix()
	return
}

func (c *cache) add(w Weather) {
	c.Lock()
	defer c.Unlock()
	c.weathers[w.ZipCode] = w
}
