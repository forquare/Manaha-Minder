package serverAdmin

import (
	"github.com/gookit/event"
	logger "github.com/sirupsen/logrus"
	"manaha_minder/utils"
	"net"
	"strings"
	"time"
)

func Watchdog() {
	logger.Debug("Setting server watchdog")
	event.On("msgevt", event.ListenerFunc(fnWatchdogHandler), event.Max)
}

var fnWatchdogHandler = func(e event.Event) error {
	log := e.Get("logline").(utils.LogLine)
	str := log.UnsplitLogLine
	if strings.Contains(str, "Server thread/ERROR") || strings.Contains(str, "Server Watchdog/ERROR") {
		logger.Debugf("Server Watchdog triggered: %s", str)

		// Try stopping the server, give some grace time for it to stop
		utils.RunServerCommand("stop", "now")

		// Wait for the server to stop
		for {
			conn, err := net.DialTimeout("tcp", net.JoinHostPort("localhost", "25565"), time.Duration(10)*time.Second)
			if err == nil {
				if conn != nil {
					logger.Debug("Waiting for server to stop...")
					time.Sleep(time.Duration(10) * time.Second)
					continue
				}
			}
			logger.Debug("Server stopped.  Waiting...")
			time.Sleep(time.Duration(30) * time.Second)
			logger.Debug("Server stopped.  Breaking...")
			break
		}

		// Restart the server
		logger.Debug("Restarting server")
		utils.RunServerCommand("start")
	}
	return nil
}
