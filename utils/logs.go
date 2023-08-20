package utils

import (
	"github.com/gookit/event"
	"github.com/nxadm/tail"
	logger "github.com/sirupsen/logrus"
	"io"
	"regexp"
	"strings"
)

type LogLine struct {
	UnsplitLogLine string
	Time           string
	Thread         string
	Message        string
	Player         string
}

var ignoreLogLines = []string{
	"Can't keep up! Is the server overloaded? Running",
	"logged in with entity id",
	"lost connection: Disconnected",
	"You whisper to",
}

func LogScraper(logFile string) {
	t, err := tail.TailFile(logFile,
		tail.Config{Location: &tail.SeekInfo{Offset: 0, Whence: io.SeekEnd}, Follow: true, ReOpen: true, MustExist: true, Logger: tail.DiscardingLogger})
	if err != nil {
		logger.Fatal(err)
	}

	defer t.Cleanup()

	for line := range t.Lines {
		if contains(ignoreLogLines, line.Text) {
			continue
		}
		l, skip, is_logonoff := SplitLog(line.Text)
		if !skip {
			if is_logonoff {
				event.MustFire("logonoffevt", event.M{"logline": l})
			}

			event.MustFire("msgevt", event.M{"logline": l})
		}
	}
}

func splitPlayerMessage(message string) (string, string, bool) {
	// Returns playerName and message

	// [19:09:44] [Server thread/INFO]: forquare joined the game
	// [19:11:39] [Server thread/INFO]: forquare left the game
	// [19:11:12] [Server thread/INFO]: <forquare> op me

	logOnOff := regexp.MustCompile(`^(\w+) ((joined|left) the game)$`)
	chatter := regexp.MustCompile(`^<(\w+)> (.*)$`)

	if logOnOff.MatchString(message) {
		logger.Trace("Matched logOnOff")
		match := logOnOff.FindStringSubmatch(message)
		return match[1], match[2], true
	} else if chatter.MatchString(message) {
		logger.Trace("Matched chatter")
		match := chatter.FindStringSubmatch(message)
		return match[1], match[2], false
	} else {
		logger.Trace("Matched nothing - no player")
		return "", message, false
	}
}

func SplitLog(line string) (LogLine, bool, bool) {
	skip := true
	re := regexp.MustCompile(`^\[(\d{2}:\d{2}:\d{2})\] \[(.*)\]: (.*)$`)
	if re.MatchString(line) {
		match := re.FindStringSubmatch(line)

		player, message, is_logonoff := splitPlayerMessage(match[3])

		ll := LogLine{UnsplitLogLine: line, Time: match[1], Thread: match[2], Message: message, Player: player}
		skip = false
		return ll, skip, is_logonoff
	}

	return LogLine{}, skip, false
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if strings.Contains(str, v) {
			return true
		}
	}
	return false
}
