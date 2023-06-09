package utils

import (
	logger "github.com/sirupsen/logrus"
	"sync"
)

type Locker struct {
	data  map[string]string
	mutex *sync.RWMutex
}

var (
	lockerOnce sync.Once
	locker     Locker
)

func InitLocker() {
	lockerOnce.Do(func() {
		locker = Locker{
			data:  make(map[string]string),
			mutex: new(sync.RWMutex),
		}
	})
}

func LockRead(key string) (string, bool) {
	// Create a read lock which will allow concurrent reads but prevent any writes
	locker.mutex.RLock()
	// Unlock when exiting the function
	defer locker.mutex.RUnlock()
	val, exists := locker.data[key]
	logger.Debugf("LockRead: %s exists: %t", key, exists)
	return val, exists
}

func LockWrite(key, value string) bool {
	created := false
	locker.mutex.RLock()
	defer locker.mutex.RUnlock()

	_, exists := locker.data[key]
	if exists {
		logger.Debugf("LockWrite: %s exists: %t", key, exists)
		return created
	} else {
		logger.Debugf("Locking %s with %s", key, value)
		locker.data[key] = value
		created = true
		return created
	}
}

func LockDelete(key string) {
	locker.mutex.RLock()
	defer locker.mutex.RUnlock()
	logger.Debugf("Deleting lock: %s", key)
	delete(locker.data, key)
}
