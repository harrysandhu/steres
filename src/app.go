package main

import (
	"sync"
	"time"

	"github.com/google/uuid"
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
	tokenSize     int
	threshold     float64
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

/**

	task - insert the sequence in the right place

				  < this computation occurs on RAM >
	get sequence : [build tokenNodes : resulting id] : put sequence


**/
func (a *App) GetSequence(key []byte) (map[string]Node, map[string]int) {
	tokenNodes := make(map[string]Node)
	countTable := make(map[string]int)
	tokens := NSplit(key, a.tokenSize)
	t, _ := uuid.NewUUID()
	id := t.String()
	for index, token := range tokens {
		data, err := a.db.Get([]byte(token), nil)
		// can we hash this?
		//tmp node
		n := Node{nvolumes: []string{}, deleted: HARD, current: token, id: id, next: getNext(&tokens, index), prev: getPrev(&tokens, index)}
		if err != leveldb.ErrNotFound {
			// byte arr -> Nodes wrapper
			var nodes Nodes = toNodes(data)

			// which one of them have next as diff(getNext)
			for _, nodeObj := range nodes.L {
				tmpnode := nodeFromMap(nodeObj)

				if nodePassesThreshold(n, tmpnode, a.threshold) {
					n = tmpnode
					id = tmpnode.id
					break
				}
			}
		}
		tokenNodes[token] = n
		// tokenNodes[dbToken] = append(tokenNodes[dbToken], n)
		countTable[id] += 1
	}
	// winId := id
	// maxCount := 0
	// for id, count := range countTable {
	// 	if count > maxCount {
	// 		maxCount = count
	// 		winId = id
	// 	}
	// }

	return tokenNodes, countTable

}
