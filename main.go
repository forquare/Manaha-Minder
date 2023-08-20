// Command fsnotify provides example usage of the fsnotify library.
package main

import (
	"github.com/forquare/manaha-minder/actions"
	"github.com/forquare/manaha-minder/config"
	"github.com/forquare/manaha-minder/player"
	"github.com/forquare/manaha-minder/serverAdmin"
	"github.com/forquare/manaha-minder/utils"
	logger "github.com/sirupsen/logrus"
)

var Version string

func main() {
	done := make(chan bool)

	// Print version
	logger.Info("Manaha Minder " + Version)

	// Set debug until we parse config
	logger.SetLevel(logger.DebugLevel)

	// Load config
	config := config.GetConfig()

	// Set log level
	l, _ := logger.ParseLevel(config.ManahaMinder.LogLevel)
	logger.SetLevel(l)

	// Load database
	utils.GetDatabase()

	// Init locker
	utils.InitLocker()

	// Start log scraper
	go utils.LogScraper(config.MinecraftServer.LatestLog)

	// Start accounting
	go player.Accounting()

	if config.Activity.LogActivity {
		// Start status setter
		if config.Activity.GenerateStatusTable {
			go player.StatusSetter()
		}

		// Run calculator
		if config.Activity.RecalculateActivityOnStartup {
			go player.RecalculateTimePlayed()
		}
	}

	// Start actions
	go actions.StartActions()

	// Set server restart timer
	if config.MinecraftServer.Restart.Enabled {
		go serverAdmin.SetServerRestart(config.MinecraftServer.Restart.Cron)
	}

	// Start watchdog
	if config.MinecraftServer.Watchdog {
		go serverAdmin.Watchdog()
	}

	// Log decompression
	if config.MinecraftServer.LogDecompress {
		go serverAdmin.LogDecompressor()
	}

	// Block forever
	<-done
}

/*
log.Trace("Something very low level.")
log.Debug("Useful debugging information.")
log.Info("Something noteworthy happened!")
log.Warn("You should probably take a look at this.")
log.Error("Something failed but I'm not quitting.")
// Calls os.Exit(1) after logging
log.Fatal("Bye.")
// Calls panic() after logging
log.Panic("I'm bailing.")
*/
