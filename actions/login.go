package actions

import (
	"fmt"
	"github.com/gookit/event"
	logger "github.com/sirupsen/logrus"
	"golang.org/x/exp/rand"
	"manaha_minder/config"
	"manaha_minder/utils"
	"regexp"
	"strings"
	"time"
)

func LoginMonitor() {
	logger.Debug("starting login monitor")
	event.On("logonoffevt", event.ListenerFunc(fnLoginHandler), event.AboveNormal)
}

var fnLoginHandler = func(e event.Event) error {
	log := e.Get("logline").(utils.LogLine)
	re := regexp.MustCompile(`joined the game`)

	if re.MatchString(log.Message) {
		logger.Debug("login event found")
		config := config.GetConfig()
		player := log.Player
		go greet(player, config.Login.WelcomeMessage)
	}
	return nil
}

func greet(player string, greeting string) {
	logger.Debugf("greeting %s", player)
	time.Sleep(time.Duration(5) * time.Second)
	greeting = strings.ReplaceAll(greeting, "<PLAYER>", player)
	for _, line := range strings.Split(strings.TrimRight(greeting, "\n"), "\n") {
		utils.DmPlayer(line, player)
	}

	if config.GetConfig().Login.GiveRandomGift {
		time.Sleep(time.Duration(2) * time.Second)
		go randomGift(player)
	}
}

type gift struct {
	name        string
	description string
}

func randomGift(player string) {
	gifts := map[int]gift{
		47: {
			name:        "bookshelf",
			description: "a sturdy bookshelf to hold your favorite books.",
		},
		263: {
			name:        "coal",
			description: "a lump of coal that can be used as fuel or crafting material.",
		},
		264: {
			name:        "diamond",
			description: "a precious diamond gemstone for crafting or trade.",
		},
		265: {
			name:        "iron_ingot",
			description: "a solid iron ingot that can be used for crafting various items.",
		},
		266: {
			name:        "gold_ingot",
			description: "a shiny gold ingot that can be used for crafting valuable items.",
		},
		276: {
			name:        "diamond_sword",
			description: "a powerful sword made of diamond, perfect for slaying enemies.",
		},
		297: {
			name:        "bread",
			description: "freshly baked bread, a tasty and filling food source.",
		},
		344: {
			name:        "egg",
			description: "a fragile egg that can be used for various recipes or breeding animals.",
		},
		329: {
			name:        "saddle",
			description: "a saddle that allows you to ride and control certain animals.",
		},
		331: {
			name:        "redstone",
			description: "a valuable mineral used for creating complex circuits and mechanisms.",
		},
		340: {
			name:        "book",
			description: "a book filled with knowledge and stories waiting to be explored.",
		},
		341: {
			name:        "slime",
			description: "a squishy slimeball that can be used in crafting and various slime-related activities.",
		},
		420: {
			name:        "lead",
			description: "a sturdy lead that can be used to guide and tether animals.",
		},
		421: {
			name:        "name_tag",
			description: "a special tag that allows you to name and personalize your pets.",
		},
	}

	randomItem := rand.Intn(432)

	logger.Debug("random item: ", randomItem)

	if gift, ok := gifts[randomItem]; ok {
		utils.DmPlayer(fmt.Sprintf("Congratulations! You have just won %s", gift.description), player)
		utils.RunMinecraftCommand(fmt.Sprintf("give %s %s", player, gift.name))
	}
}
