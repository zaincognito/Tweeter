package main

import(
	"sync"
)


type Lock struct{
	mut sync.RWMutex
}

//gets a lock, if read->shared, write->exclusive
func acquireLock(lockType string, lock *Lock) {
	
	//if read lock requested, do a share lock
	if lockType == "read" {

		lock.mut.RLock()

	//if write lock requested, do an exclusive lock
	} else if lockType == "write" {

		lock.mut.Lock()

	} else {
		panic("Locking: Not read nor write")
	}
}

//takes lock, if read->unlock the shared lock (might have others still using)
//or releases totally if write lock
func releaseLock(lockType string, lock *Lock) {
	
	if lockType == "read" {
		lock.mut.RUnlock()
	} else if lockType == "write" {
		lock.mut.Unlock()
	} else {
		panic("Unlocking: Not read nor write")
	}
}