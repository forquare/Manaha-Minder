package utils

import (
	logger "github.com/sirupsen/logrus"
	"manaha_minder/config"
	"os/exec"
)

func RunMinecraftCommand(args ...string) string {
	logger.Debug("Running Minecraft command")
	config := config.GetConfig()
	s := append([]string{config.MinecraftServer.ServerName}, "cmd")
	s = append(s, args...)
	return RunMsmCommand(s)
}

func RunBroadcastCommand(args ...string) string {
	logger.Debug("Running broadcast command")
	config := config.GetConfig()
	s := append([]string{config.MinecraftServer.ServerName}, "say")
	s = append(s, args...)
	return RunMsmCommand(s)
}

func DmPlayer(message string, player string) string {
	logger.Debug("Running DM command")
	return RunMinecraftCommand("/tell", player, message)
}

func RunServerCommand(args ...string) string {
	logger.Debug("Running server command")
	config := config.GetConfig()
	s := append([]string{config.MinecraftServer.ServerName}, args...)
	return RunMsmCommand(s)
}

func RunMsmCommand(s []string) string {
	config := config.GetConfig()
	cmd := exec.Command(config.MinecraftServer.MsmBinary, s...)
	logger.Debugf("Running MSM command: %s", cmd.String())
	out, err := cmd.Output()
	if err != nil {
		logger.Error(err)
	}
	logger.Debugf("MSM command output: %s", out)
	return string(out)
}
