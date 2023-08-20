package utils

import (
	logger "github.com/sirupsen/logrus"
	"sync"
)

type Locker struct {
	data  map[string]string
	mutex *sync.RWMutex
}

type DBLock struct {
	ID     uint   `gorm:"primaryKey"`
	Key    string `gorm:"uniqueIndex"`
	Value  string
	Locked bool
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
		err := GetDatabase().AutoMigrate(&DBLock{})
		if err != nil {
			logger.Panicf("Error initializing locker: %v\n", err)
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

func LockReadDB(key string) (string, bool) {
	locker.mutex.RLock()
	defer locker.mutex.RUnlock()
	var dbLock DBLock

	GetDatabase().Where("key = ?", key).First(&dbLock)
	logger.Debugf("LockReadDB: %s exists: %t", key, dbLock.Locked)
	return dbLock.Value, dbLock.Locked
}

func LockWrite(key, value string) bool {
	locker.mutex.RLock()
	defer locker.mutex.RUnlock()

	_, exists := locker.data[key]
	if exists {
		logger.Debugf("LockWrite: %s exists: %t", key, exists)
		return false
	} else {
		logger.Debugf("Locking %s with %s", key, value)
		locker.data[key] = value
		return true
	}
}

func LockWriteDB(key, value string) bool {
	locker.mutex.RLock()
	defer locker.mutex.RUnlock()
	var dbLock DBLock

	GetDatabase().Where("key = ?", key).First(&dbLock)
	if dbLock.Locked {
		logger.Debugf("LockWriteDB: %s exists: %t", key, dbLock.Locked)
		return false
	} else {
		logger.Debugf("Locking %s with %s", key, value)
		dbLock.Key = key
		dbLock.Value = value
		dbLock.Locked = true
		GetDatabase().Create(&dbLock)
		return true
	}
}

func LockDelete(key string) {
	locker.mutex.RLock()
	defer locker.mutex.RUnlock()
	logger.Debugf("Deleting lock: %s", key)
	delete(locker.data, key)
}

func LockDeleteDB(key string) {
	locker.mutex.RLock()
	defer locker.mutex.RUnlock()
	logger.Debugf("Deleting lock: %s", key)
	GetDatabase().Where("key = ?", key).Delete(&DBLock{})
}
