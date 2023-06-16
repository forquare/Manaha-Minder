package actions

import (
	"fmt"
	"github.com/gookit/event"
	logger "github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
	"github.com/forquare/manaha-minder/config"
	"github.com/forquare/manaha-minder/utils"
	"regexp"
	"time"
)

func OperatorMonitor() {
	logger.Debug("Starting Operator Monitor")
	event.On("msgevt", event.ListenerFunc(fnMonitorHandler), event.AboveNormal)
}

var fnMonitorHandler = func(e event.Event) error {
	log := e.Get("logline").(utils.LogLine)
	re := regexp.MustCompile("^op me$")

	if re.MatchString(log.Message) {
		logger.Debug("Operator Monitor event")
		config := config.GetConfig()
		player := log.Player
		if slices.Contains(config.Operator.Operators, player) {
			d := config.Operator.Duration
			go makeOperator(player, d)
		} else {
			logger.Debugf("Player %s is not an authorised operator", player)
			utils.RunBroadcastCommand(fmt.Sprintf("Sorry %s, you are not an authorized operator.", player))
		}
	}
	return nil
}

func makeOperator(player string, duration int) {
	if utils.LockWrite(fmt.Sprintf("operator_%s", player), player) {
		logger.Debugf("Locking operator for %s", player)
		defer utils.LockDelete(fmt.Sprintf("operator_%s", player))
	} else {
		logger.Debugf("Player %s is already being made an operator", player)
		return
	}

	logger.Debugf("Making %s an operator for %d minutes", player, (duration / 60))
	utils.RunServerCommand("op", "add", player)
	utils.RunBroadcastCommand(fmt.Sprintf("%s has been made an operator for %d minutes.", player, (duration / 60)))

	time.Sleep(time.Duration(duration-10) * time.Second)

	utils.DmPlayer(fmt.Sprintf("10 SECONDS LEFT OF OP POWERS!"), player)
	utils.DmPlayer(fmt.Sprintf("WILL REMOVE CREATIVE MODE TOO!!"), player)

	time.Sleep(time.Duration(10) * time.Second)

	logger.Debugf("Removing operator status from %s", player)
	utils.RunServerCommand("op", "remove", player)
	utils.RunMinecraftCommand(fmt.Sprintf("/gamemode survival %s", player))
	utils.RunBroadcastCommand(fmt.Sprintf("%s is no longer an operator.", player))

	logger.Debugf("Operator Monitor finished for player %s", player)
}
