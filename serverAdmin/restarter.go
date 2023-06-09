package serverAdmin

import (
	"github.com/go-co-op/gocron"
	logger "github.com/sirupsen/logrus"
	"manaha_minder/utils"
	"time"
)

func SetServerRestart(t string) {
	logger.Debug("Setting server restart")
	s := gocron.NewScheduler(time.UTC)
	s.Cron(t).Do(func() {
		logger.Debug("Restarting server")
		utils.RunServerCommand("say", "The server is going down for a restart. It will be back up after a few minutes.  10 second countdown")
		time.Sleep(5 * time.Second)
		utils.RunServerCommand("say", "The server is going down for a restart in 5 seconds.")
		time.Sleep(5 * time.Second)
		utils.RunServerCommand("say", "The server is going down for a restart NOW. Back up soon.")
		utils.RunServerCommand("say", "The server is going down for a restart NOW. Back up soon.")

		utils.RunServerCommand("restart", "now")
	})
}
