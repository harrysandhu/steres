package main

import (
	"sync"
	"time"

	"github.com/syndtr/goleveldb/leveldb"
)

/**
App struct maintains a central state for
the server
**/
type App struct {
	tokenNodes    map[string][]Node
	db            *leveldb.DB
	mutexLock     sync.Mutex
	lock          map[string]struct{}
	volumes       []string
	fallback      string
	replicas      int
	subvolumes    int
	protect       bool
	md5sum        bool
	volumeTimeout time.Duration
	tokenSize     int
}

/*
* Unlock key
 */

func (a *App) UnlockKey(key []byte) {
	a.mutexLock.Lock()
	delete(a.lock, string(key))
	a.mutexLock.Unlock()
}

/*
* Lock key
* Put a mutex lock on an address.
**/
func (a *App) LockKey(key []byte) bool {
	a.mutexLock.Lock()
	defer a.mutexLock.Unlock()
	if _, ok := a.lock[string(key)]; ok {
		return false
	}
	a.lock[string(key)] = struct{}{}
	return true
}

// lock the sequence
func (a *App) LockSequence(key []byte) bool {
	// Split the sequence bytearray at " ", of size = size
	tks := NSplit(key, a.tokenSize)
	a.mutexLock.Lock()
	defer a.mutexLock.Unlock()
	for _, value := range tks {
		if _, ok := a.lock[value]; ok {
			return false
		}
		a.lock[value] = struct{}{}
	}
	return true
}

// Unlock the sequence
func (a *App) UnlockSequence(key []byte) {
	a.mutexLock.Lock()
	tks := NSplit(key, a.tokenSize)
	for _, value := range tks {
		delete(a.lock, value)
	}
	a.mutexLock.Unlock()
}
