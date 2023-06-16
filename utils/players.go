package utils

import (
	"encoding/json"
	logger "github.com/sirupsen/logrus"
	"github.com/forquare/manaha-minder/config"
	"os"
	"time"
)

type Player struct {
	Name           string `json:"name"`
	Uuid           string `json:"uuid"`
	Activity       time.Duration
	LastJoined     time.Time
	LastJoinedFlag bool // True if we've read a LastJoined
}

func GetPlayers() map[string]*Player {
	players := make([]Player, 0)
	playerMap := make(map[string]*Player)

	jsonData, err := os.Open(config.GetConfig().MinecraftServer.WhitelistFile)
	if err != nil {
		logger.Errorf("Error reading JSON file: %v\n", err)
	}
	defer jsonData.Close()

	decoder := json.NewDecoder(jsonData)

	err = decoder.Decode(&players)

	logger.Tracef("Players: %v", players)

	if err != nil {
		logger.Panicf("Error decoding JSON: %v\n", err)
	}

	for _, player := range players {
		p := player
		playerMap[player.Name] = &p
		playerMap[player.Name].LastJoinedFlag = false
		logger.Tracef("Player: %v", playerMap[player.Name])
	}

	logger.Tracef("PlayerMap: %v", playerMap)

	return playerMap
}
