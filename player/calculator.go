package player

import (
	"fmt"
	"github.com/forquare/manaha-minder/config"
	"github.com/forquare/manaha-minder/utils"
	logger "github.com/sirupsen/logrus"
	"os"
	"strings"
	"time"
)

func RecalculateTimePlayed() {
	initPlayers()
	config := config.GetConfig()
	logger.Tracef("Activity calculator: %v", config)

	for {
		if utils.LockWrite("activityCalc", "activityCalc") {
			logger.Tracef("Got lock: activityCalc")
			break
		}
		time.Sleep(1 * time.Second)
	}
	defer utils.LockDelete("activityCalc")

	players := GetPlayersByUuid()

	for uuid := range players {
		var activities []Activity
		utils.GetDatabase().Where("player_uuid = ?", uuid).Find(&activities)
		var timePlayed int64 = 0

		for _, activity := range activities {
			if strings.Contains(activity.Action, "joined") {
				logger.Tracef("Activity calculator joined: Player: %v ", activity)
				players[uuid].RecalculateLastJoined = activity.Timestamp
				players[uuid].RecalculateLastJoinedFlag = true
			} else if strings.Contains(activity.Action, "left") {
				if players[uuid].RecalculateLastJoinedFlag {
					timePlayed += activity.Timestamp - players[uuid].RecalculateLastJoined
					players[uuid].RecalculateLastJoinedFlag = false
				}
				logger.Tracef("Activity calculator left: Player %s | TimePlayed %d | LastJoined %d | LastLeft %d", players[uuid].Name, timePlayed, players[uuid].RecalculateLastJoined, activity.Timestamp)
			}
		}
		SetTimePlayed(uuid, timePlayed)
	}

	writeTimePlayedTable()
}

func updateTimePlayed(playerUuid string, logOfftime int64) {
	// Called from accounting.go when a player logs off

	for {
		if utils.LockWrite("activityCalc", "activityCalc") {
			logger.Tracef("Got lock: activityCalc")
			break
		}
		time.Sleep(1 * time.Second)
	}
	defer utils.LockDelete("activityCalc")

	config := config.GetConfig()

	if !config.Activity.LogActivity {
		// Don't run if activity logging is disabled
		return
	}

	player, err := GetPlayerByUuid(playerUuid)
	if err != nil {
		logger.Errorf("Error getting player: %v", err)
	}

	player.TimePlayed += logOfftime - player.LastJoined
	utils.GetDatabase().Save(&player)

	writeTimePlayedTable()
}

func writeTimePlayedTable() {
	config := config.GetConfig()

	if !config.Activity.LogActivity {
		// Don't run if activity logging is disabled
		return
	}
	if !config.Activity.GenerateTimePlayedOutput {
		// Don't run if output generation is disabled
		return
	}

	for {
		if utils.LockWrite("activityWrite", "activityWrite") {
			logger.Tracef("Got lock: activityWrite")
			break
		}
		time.Sleep(1 * time.Second)
	}
	defer utils.LockDelete("activityWrite")

	players := GetPlayersSortedByTimePlayed()

	activityOut, err := os.OpenFile(config.Activity.TimePlayedFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	defer activityOut.Close()
	if err != nil {
		logger.Errorf("Error opening activity log: %v\n", err)
		return
	}

	for _, player := range players {
		timePlayed := time.Duration(player.TimePlayed * int64(time.Second))
		if err != nil {
			logger.Errorf("Error parsing time: %v\n", err)
		}
		_, err = activityOut.WriteString(fmt.Sprintf("%s: %s\n", player.Name, timePlayed))
		if err != nil {
			logger.Errorf("Error writing to activity log: %v\n", err)
		}
	}
}
