package player

import (
	"github.com/forquare/manaha-minder/config"
	"github.com/forquare/manaha-minder/utils"
	"github.com/gookit/event"
	logger "github.com/sirupsen/logrus"
	"regexp"
	"time"
)

type Activity struct {
	ID         uint `gorm:"primaryKey"`
	Timestamp  int64
	PlayerUuid string
	Action     string
}

func Accounting() {
	logger.Debug("Starting activity")
	err := utils.GetDatabase().AutoMigrate(&Activity{})
	if err != nil {
		logger.Panicf("Error initializing activity: %v\n", err)
	}
	event.On("logonoffevt", event.ListenerFunc(fnAccountingHandler), event.Low)
}

var fnAccountingHandler = func(e event.Event) error {
	log := e.Get("logline").(utils.LogLine)
	logOn := regexp.MustCompile(`joined the game$`)
	logOff := regexp.MustCompile(`left the game$`)
	config := config.GetConfig()
	action := ""
	playerUuid := WhitelistPlayerName2UUID(log.Player)
	insertTime := time.Now().Unix()

	logger.Debug("Accounting event found")

	if logOn.MatchString(log.Message) {
		action = "joined"
	} else if logOff.MatchString(log.Message) {
		action = "left"
	}

	activity := Activity{PlayerUuid: playerUuid, Timestamp: insertTime, Action: action}
	if config.Activity.LogActivity {
		utils.GetDatabase().Create(&activity)
	}

	if logOn.MatchString(log.Message) {
		go SyncPlayers(activity)
	} else if logOff.MatchString(log.Message) {
		go updateTimePlayed(playerUuid, insertTime)
	}

	return nil
}
