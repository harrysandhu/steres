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

// get record

// put record
