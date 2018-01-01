package book

import (
	"sync"
)

const cacheSize = 15

var cache = map[string]Book{}
var cacheSavings = []string{}
var cacheMutex = &sync.Mutex{}

func cacheUpdate(path string, b Book, ok bool, err error) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	if _, ok := cache[path]; ok {
		return
	}

	if ok && err == nil {
		cache[path] = b
		cacheSavings = append(cacheSavings, path)
		if len(cacheSavings) > cacheSize {
			cacheSavings = cacheSavings[len(cacheSavings)-cacheSize:]
		}
	}
}
