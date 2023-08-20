package actions

import (
	"github.com/forquare/manaha-minder/config"
	"github.com/forquare/manaha-minder/utils"
	"github.com/gookit/event"
	logger "github.com/sirupsen/logrus"
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
	configuration := config.GetConfig()

	if len(configuration.CustomActions) > 0 {
		for _, action := range configuration.CustomActions {
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
