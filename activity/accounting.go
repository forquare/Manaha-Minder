package activity

import (
	"fmt"
	"github.com/gookit/event"
	logger "github.com/sirupsen/logrus"
	"manaha_minder/config"
	"manaha_minder/utils"
	"os"
	"regexp"
	"time"
)

func Accounting() {
	logger.Debug("Starting activity")
	event.On("logonoffevt", event.ListenerFunc(fnAccountingHandler), event.Low)
}

var fnAccountingHandler = func(e event.Event) error {
	log := e.Get("logline").(utils.LogLine)
	logOn := regexp.MustCompile(`joined the game$`)
	logOff := regexp.MustCompile(`left the game$`)

	logger.Debug("Accounting event found")
	player := log.Player
	activityMessage := ""
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	if logOn.MatchString(log.Message) {
		activityMessage = fmt.Sprintf("%s %s joined\n", currentTime, player)
	} else if logOff.MatchString(log.Message) {
		activityMessage = fmt.Sprintf("%s %s left\n", currentTime, player)
		go ActivityCalulator()
	}

	logger.Debugf("Activity message: %s", activityMessage)

	f, err := os.OpenFile(config.GetConfig().Activity.Log, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logger.Error("Could not open activity log for ", player)
	}
	if _, err := f.Write([]byte(activityMessage)); err != nil {
		f.Close() // ignore error; Write error takes precedence
		logger.Error("Could not log activity for", player)
	}
	if err := f.Close(); err != nil {
		logger.Panic("Could not close activity log")
	}

	return nil
}
