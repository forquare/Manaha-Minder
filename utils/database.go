package utils

import (
	"github.com/forquare/manaha-minder/config"
	logger "github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"sync"
)

var (
	databaseOnce sync.Once
	database     *gorm.DB
)

func GetDatabase() *gorm.DB {
	logger.Debugf("Loading database")
	databaseOnce.Do(func() {
		config := config.GetConfig()
		logger.Tracef("Database file: %s", config.ManahaMinder.DatabaseFile)
		db, err := gorm.Open(sqlite.Open(config.ManahaMinder.DatabaseFile), &gorm.Config{})
		if err != nil {
			logger.Panicln("failed to connect database")
		}
		database = db
	})
	return database
}
