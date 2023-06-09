package activity

import (
	"bufio"
	"fmt"
	"github.com/gookit/event"
	logger "github.com/sirupsen/logrus"
	"manaha_minder/config"
	"manaha_minder/utils"
	"os"
	"regexp"
	"strings"
	"time"
)

func StatusSetter() {
	logger.Debug("Starting status setter")
	event.On("logonoffevt", event.ListenerFunc(fnStatusHandler), event.Low)
}

var fnStatusHandler = func(e event.Event) error {
	log := e.Get("logline").(utils.LogLine)
	logOn := regexp.MustCompile(`joined the game$`)
	logOff := regexp.MustCompile(`left the game$`)

	logger.Debug("Status setter event found")
	player := log.Player
	date := time.Now().Format("02 Jan 2006")
	time := time.Now().Format("15:04")
	in := fmt.Sprintf("<tr bgcolor='LimeGreen'><td>%s</td><td><b>ACTIVE NOW!</b></td></tr>", player)
	out := fmt.Sprintf("<tr bgcolor='red'><td>%s</td><td>%s&nbsp;&nbsp;&nbsp;&nbsp;%s</td></tr>", player, time, date)

	output := make([]string, 0)

	if logOn.MatchString(log.Message) {
		output = append(output, in)
		logger.Debugf("Status setter: %s joined", player)
		logger.Debug(output)
	} else if logOff.MatchString(log.Message) {
		output = append(output, out)
		logger.Debugf("Status setter: %s left", player)
		logger.Debug(output)
	}

	f, err := os.OpenFile(config.GetConfig().Activity.Status, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		logger.Error("Could not open activity output file", player, err)
	}

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		tmp := scanner.Text()
		if !strings.Contains(tmp, player) {
			output = append(output, tmp)
		}
	}

	if err := scanner.Err(); err != nil {
		logger.Error("Could not read activity output file", player, err)
	}
	if err := f.Close(); err != nil {
		logger.Panic("Could not close activity output file", err)
	}

	logger.Debugf("Status setter: %s", output)

	logger.Debug("Status setter: writing to activity output file")

	f, err = os.OpenFile(config.GetConfig().Activity.Status, os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		logger.Error("Could not open activity output file", player, err)
	}

	for _, line := range output {
		_, err := f.WriteString(strings.TrimRight(line, "\n") + "\n")
		if err != nil {
			logger.Panic("Could not write to activity output file", err)
		}
	}

	if err := f.Close(); err != nil {
		logger.Panic("Could not close activity output file", err)
	}

	return nil
}
