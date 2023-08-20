package actions

import (
	"fmt"
	"github.com/forquare/manaha-minder/config"
	"github.com/forquare/manaha-minder/player"
	"github.com/forquare/manaha-minder/utils"
	"github.com/gookit/event"
	logger "github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
	"regexp"
	"time"
)

func OperatorMonitor() {
	logger.Debug("Starting Operator Monitor")

	for _, p := range player.GetPlayersByUuid() {
		_, lockExists := utils.LockReadDB(fmt.Sprintf("operator_%s", p.Uuid))
		if lockExists {
			logger.Debugf("Removing lock for %s", p.Name)
			deOpNow(*p)
		}
	}

	event.On("msgevt", event.ListenerFunc(fnMonitorHandler), event.AboveNormal)
}

var fnMonitorHandler = func(e event.Event) error {
	log := e.Get("logline").(utils.LogLine)
	re := regexp.MustCompile("^op me$")

	if re.MatchString(log.Message) {
		logger.Debug("Operator Monitor event")
		config := config.GetConfig()
		player, err := player.GetPlayerByName(log.Player)
		if err != nil {
			logger.Errorf("Could not get player %s", log.Player)
		}
		if slices.Contains(config.Operator.Operators, player.Name) {
			d := config.Operator.Duration
			go makeOperator(*player, d)
		} else {
			logger.Debugf("Player %s is not an authorised operator", player.Name)
			utils.RunBroadcastCommand(fmt.Sprintf("Sorry %s, you are not an authorized operator.", player.Name))
		}
	}
	return nil
}

func makeOperator(player player.Player, duration int) {
	if utils.LockWriteDB(fmt.Sprintf("operator_%s", player.Uuid), player.Uuid) {
		logger.Debugf("Locking operator for %s", player.Name)
		defer utils.LockDeleteDB(fmt.Sprintf("operator_%s", player.Uuid))
	} else {
		logger.Debugf("Player %s is already being made an operator", player.Name)
		return
	}

	logger.Debugf("Making %s an operator for %d minutes", player.Name, (duration / 60))
	utils.RunServerCommand("op", "add", player.Name)
	utils.RunBroadcastCommand(fmt.Sprintf("%s has been made an operator for %d minutes.", player.Name, (duration / 60)))

	time.Sleep(time.Duration(duration-10) * time.Second)

	deOp(player)
}

func deOp(player player.Player) {
	utils.DmPlayer("10 SECONDS LEFT OF OP POWERS!", player.Name)
	utils.DmPlayer("WILL REMOVE CREATIVE MODE TOO!!", player.Name)

	time.Sleep(time.Duration(10) * time.Second)

	deOpNow(player)
}

func deOpNow(player player.Player) {
	logger.Debugf("Removing operator status from %s", player.Name)
	utils.RunServerCommand("op", "remove", player.Name)
	utils.RunMinecraftCommand(fmt.Sprintf("/gamemode survival %s", player.Name))
	utils.RunBroadcastCommand(fmt.Sprintf("%s is no longer an operator.", player.Name))

	_, lockExists := utils.LockReadDB(fmt.Sprintf("operator_%s", player.Uuid))
	if lockExists {
		utils.LockDeleteDB(fmt.Sprintf("operator_%s", player.Uuid))
	}

	logger.Debugf("Operator Monitor finished for player %s", player.Name)
}
