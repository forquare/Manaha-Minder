// Command fsnotify provides example usage of the fsnotify library.
package main

import (
	logger "github.com/sirupsen/logrus"
	"manaha_minder/actions"
	"manaha_minder/activity"
	"manaha_minder/config"
	"manaha_minder/serverAdmin"
	"manaha_minder/utils"
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

	// Init locker
	utils.InitLocker()

	// Start log scraper
	go utils.LogScraper(config.MinecraftServer.LatestLog)

	// Start accounting
	if config.Activity.LogActivity {
		go activity.Accounting()

		if config.Activity.GenerateStatus {
			// Start status setter
			go activity.StatusSetter()
		}

		if config.Activity.GenerateOutput {
			// Run calculator
			go activity.ActivityCalulator()
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
