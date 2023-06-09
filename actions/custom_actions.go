package actions

import (
	"github.com/gookit/event"
	logger "github.com/sirupsen/logrus"
	"manaha_minder/config"
	"manaha_minder/utils"
	"regexp"
	"strings"
	"time"
)

func CustomActions() {
	logger.Debug("Starting custom actions")
	event.On("msgevt", event.ListenerFunc(fnCustomActionsHandler), event.Low)
}

var fnCustomActionsHandler = func(e event.Event) error {
	log := e.Get("logline").(utils.LogLine)
	go customLogAction(log)
	return nil
}

func customLogAction(log utils.LogLine) {
	config := config.GetConfig()

	if len(config.CustomActions) > 0 {
		for _, action := range config.CustomActions {
			re := regexp.MustCompile(action.Pattern)
			if re.MatchString(log.Message) {
				logger.Debugf("Found custom action: %s", action.Name)
				for _, command := range action.Commands {
					time.Sleep(time.Duration(1) * time.Second)
					cmd := strings.ReplaceAll(command, "<PLAYER>", log.Player)
					utils.RunMinecraftCommand(cmd)
				}
			}
		}
	}
}
