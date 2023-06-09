package activity

import (
	"bufio"
	"fmt"
	logger "github.com/sirupsen/logrus"
	"manaha_minder/config"
	"manaha_minder/utils"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"
)

type activityLogLine struct {
	date   string
	time   string
	player string
	action string
}

func ActivityCalulator() {
	// Gets called from accounting.go when a player logs off

	config := config.GetConfig()

	if !config.Activity.LogActivity {
		// Don't run if activity logging is disabled
		return
	}
	if !config.Activity.GenerateOutput {
		// Don't run if output generation is disabled
		return
	}

	if utils.LockWrite("activity", "activity") {
		logger.Debug("Starting activity calculator")
		defer utils.LockDelete("activity")
	} else {
		logger.Debug("Activity calculator already running")
		return
	}

	players := utils.GetPlayers()
	activityLog, err := os.Open(config.Activity.Log)

	if err != nil {
		logger.Errorf("Error reading activity log: %v\n", err)
	}

	defer activityLog.Close()

	scanner := bufio.NewScanner(activityLog)

	for scanner.Scan() {
		line := scanner.Text()
		splitLine := splitActivityLogLine(line)
		if strings.Contains(splitLine.action, "joined") {
			logger.Tracef("Activity calculator: Player: %v ", splitLine)
			players[splitLine.player].LastJoined, _ = time.Parse(time.DateTime, fmt.Sprintf("%s %s", splitLine.date, splitLine.time))
			players[splitLine.player].LastJoinedFlag = true
		} else if strings.Contains(splitLine.action, "left") {
			logger.Tracef("Activity calculator: Player: %v ", splitLine)
			// TODO this is failing
			if players[splitLine.player].LastJoinedFlag {
				leftTime, _ := time.Parse(time.DateTime, fmt.Sprintf("%s %s", splitLine.date, splitLine.time))
				duration := leftTime.Sub(players[splitLine.player].LastJoined)
				tmp := 0.0

				logger.Tracef("Activity calculator PRE: Player %s, leftTime %s, joinTime %s, duration %s, tmp %f, activity %f", splitLine.player, leftTime, players[splitLine.player].LastJoined, duration, tmp, players[splitLine.player].Activity.Seconds())

				players[splitLine.player].Activity += duration

				players[splitLine.player].LastJoinedFlag = false

				logger.Tracef("Activity calculator POST: Player %s, leftTime %s, joinTime %s, duration %s, tmp %f, activity %f", splitLine.player, leftTime, players[splitLine.player].LastJoined, duration, tmp, players[splitLine.player].Activity.Seconds())
			}
		}
	}

	var playerNames []string
	for name, _ := range players {
		playerNames = append(playerNames, name)
	}

	sort.Slice(playerNames, func(i, j int) bool {
		return players[playerNames[i]].Activity.Seconds() > players[playerNames[j]].Activity.Seconds()
	})

	activityOut, err2 := os.OpenFile(config.Activity.Output, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err2 != nil {
		logger.Errorf("Error opening activity log: %v\n", err)
	}

	defer activityOut.Close()

	for _, name := range playerNames {
		// Make the time output prettier
		_, err = activityOut.WriteString(fmt.Sprintf("%s: %s\n", name, players[name].Activity))
		if err != nil {
			logger.Errorf("Error writing to activity log: %v\n", err)
		}
	}
}

func splitActivityLogLine(line string) activityLogLine {
	re := regexp.MustCompile(`^(\d{4}-\d{2}-\d{2}) (\d{2}:\d{2}:\d{2}) (.+) (.*)$`)

	if !re.MatchString(line) {
		logger.Panicf("Could not parse activity log line: %s", line)
	}

	match := re.FindStringSubmatch(line)
	var activityLL activityLogLine

	activityLL.date = match[1]
	activityLL.time = match[2]
	activityLL.player = match[3]
	activityLL.action = match[4]
	if strings.Contains(match[4], "joined the game") {
		activityLL.action = "joined"
	} else if strings.Contains(match[4], "left the game") {
		activityLL.action = "left"
	}

	return activityLL
}
