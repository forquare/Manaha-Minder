package player

import (
	"encoding/json"
	"github.com/forquare/manaha-minder/config"
	"github.com/forquare/manaha-minder/utils"
	logger "github.com/sirupsen/logrus"
	"os"
	"sort"
	"sync"
)

type Player struct {
	Name                      string `json:"name"`
	Uuid                      string `json:"uuid" gorm:"primaryKey"`
	TimePlayed                int64
	LastJoined                int64
	RecalculateLastJoinedFlag bool  // Used for recalculations, true if we've read a LastJoined
	RecalculateLastJoined     int64 // Used for recalculations
	Active                    bool  // False if player is in DB but not in whitelist
}

var (
	playersOnce sync.Once
	syncLock    sync.Mutex
)

func initPlayers() {
	playersOnce.Do(func() {
		logger.Debugf("Initializing players")
		utils.GetDatabase().AutoMigrate(&Player{})
	})
}

func readWhitelist() []Player {
	players := make([]Player, 0)

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

	return players
}

func getPlayers() []Player {
	var players []Player
	if result := utils.GetDatabase().Find(&players); result.Error != nil {
		logger.Panicf("Error reading players from DB: %v\n", result.Error)
	}
	return players
}

func SyncPlayers(activities ...Activity) {
	initPlayers()

	// Some calls will be from the activity calculator, which will pass a UUID
	// if we've seen this player before we don't need to sync.
	if len(activities) > 0 {
		for _, activity := range activities {
			dbPlayer, err := GetPlayerByUuid(activity.PlayerUuid)
			if err != nil {
				logger.Tracef("Player %s not in DB", activity.PlayerUuid)
				break
			}
			if activity.PlayerUuid == dbPlayer.Uuid {
				// Set LastJoined
				dbPlayer.LastJoined = activity.Timestamp
				utils.GetDatabase().Save(&dbPlayer)

				// Check player name is correct
				if WhitelistPlayerUUID2Name(activity.PlayerUuid) != dbPlayer.Name {
					logger.Tracef("Player %s name incorrect, updating", activity.PlayerUuid)
					dbPlayer.Name = WhitelistPlayerUUID2Name(activity.PlayerUuid)
					utils.GetDatabase().Save(&dbPlayer)
				}
				return
			}
		}
	}

	logger.Tracef("Locking syncLock")
	syncLock.Lock()
	defer syncLock.Unlock()

	whitelistPlayers := readWhitelist()

	dbPlayers := getPlayers()

	// Add players from whitelist to DB
	for _, player := range whitelistPlayers {
		logger.Tracef("Player: %v", player)
		found := false
		for _, dbPlayer := range dbPlayers {
			logger.Tracef("DBPlayer: %v", dbPlayer)
			if player.Uuid == dbPlayer.Uuid {
				logger.Tracef("Player %s already in DB", player.Name)
				found = true
			}
		}
		if !found {
			logger.Tracef("Player %s not in DB", player.Name)
			utils.GetDatabase().Create(&Player{Name: player.Name, Uuid: player.Uuid, Active: true})
		}
	}

	for _, dbPlayer := range dbPlayers {
		found := false
		for _, player := range whitelistPlayers {
			if dbPlayer.Uuid == player.Uuid {
				found = true
			}
		}
		if !found {
			logger.Tracef("Player %s not in whitelist", dbPlayer.Name)
			dbPlayer.Active = false
			utils.GetDatabase().Save(&dbPlayer)
		}
	}
}

func GetPlayersByName() map[string]*Player {
	players := getPlayers()
	playerMap := make(map[string]*Player)

	for _, player := range players {
		p := player
		playerMap[player.Name] = &p
		logger.Tracef("Player: %v", playerMap[player.Name])
	}

	logger.Tracef("PlayerMap: %v", playerMap)

	return playerMap
}

func GetPlayersByUuid() map[string]*Player {
	players := getPlayers()
	playerMap := make(map[string]*Player)

	for _, player := range players {
		p := player
		playerMap[player.Uuid] = &p
		logger.Tracef("Player: %v", playerMap[player.Uuid])
	}

	logger.Tracef("PlayerMap: %v", playerMap)

	return playerMap
}

func GetPlayersSortedByTimePlayed() []*Player {
	players := GetPlayersByUuid()
	var uuids []string
	var sortedPlayers []*Player

	for _, player := range players {
		uuids = append(uuids, player.Uuid)
	}

	sort.Slice(uuids, func(i, j int) bool {
		return players[uuids[i]].TimePlayed > players[uuids[j]].TimePlayed
	})

	for _, uuid := range uuids {
		sortedPlayers = append(sortedPlayers, players[uuid])
	}

	return sortedPlayers
}

func GetPlayerByName(playerName string) (*Player, error) {
	var player Player
	if result := utils.GetDatabase().Where("name = ?", playerName).First(&player); result.Error != nil {
		logger.Errorf("Error reading player from DB: %v\n", result.Error)
		return nil, result.Error
	} else {
		return &player, nil
	}
}

func GetPlayerByUuid(playerUUID string) (*Player, error) {
	var player Player
	if result := utils.GetDatabase().Where("uuid = ?", playerUUID).First(&player); result.Error != nil {
		logger.Errorf("Error reading player from DB: %v\n", result.Error)
		return nil, result.Error
	} else {
		return &player, nil
	}
}

func PlayerName2UUID(playerName string) string {
	players := GetPlayersByName()
	return players[playerName].Uuid
}

func WhitelistPlayerName2UUID(playerName string) string {
	players := readWhitelist()

	for _, player := range players {
		if player.Name == playerName {
			return player.Uuid
		}
	}

	return ""
}

func PlayerUUID2Name(playerUUID string) string {
	players := GetPlayersByUuid()
	return players[playerUUID].Name
}
func WhitelistPlayerUUID2Name(playerUUID string) string {
	players := readWhitelist()

	for _, player := range players {
		if player.Uuid == playerUUID {
			return player.Name
		}
	}

	return ""
}

func GetLastJoinedActivity(playerUUID string) int64 {
	player, err := GetPlayerByUuid(playerUUID)
	if err != nil {
		logger.Errorf("Error getting player: %v", err)
	}
	return player.LastJoined
}

func SetTimePlayed(playerUUID string, timePlayed int64) {
	player, err := GetPlayerByUuid(playerUUID)
	if err != nil {
		logger.Errorf("Error getting player: %v", err)
	}
	player.TimePlayed = timePlayed
	utils.GetDatabase().Save(&player)
}
