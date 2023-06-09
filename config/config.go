package config

import (
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"runtime"
	"sync"
)

var (
	once   sync.Once
	config Config
)

type MinecraftServerConfig struct {
	ServerName    string `mapstructure:"server_name"`
	LatestLog     string `mapstructure:"latest_log_file"`
	LogDir        string `mapstructure:"log_dir"`
	WorldDir      string `mapstructure:"world_dir"`
	WhitelistFile string `mapstructure:"whitelist_file"`
	MsmBinary     string `mapstructure:"msm_binary"`
	Watchdog      bool   `mapstructure:"watchdog"`
	LogDecompress bool   `mapstructure:"log_decompress"`
	Restart       struct {
		Enabled bool   `mapstructure:"enabled"`
		Cron    string `mapstructure:"cron"`
	}
}

type ManahaMinder struct {
	LogLevel string `mapstructure:"log_level"`
}

type LoginConfig struct {
	WelcomeMessage string `mapstructure:"welcome_message"`
	GiveRandomGift bool   `mapstructure:"give_random_gift"`
}

type ActivityConfig struct {
	Log            string `mapstructure:"log"`
	LogActivity    bool   `mapstructure:"log_activity"`
	Output         string `mapstructure:"output"`
	GenerateOutput bool   `mapstructure:"generate_output"`
	Status         string `mapstructure:"status"`
	GenerateStatus bool   `mapstructure:"generate_status"`
}

type OperatorConfig struct {
	Duration  int      `mapstructure:"duration"`
	Operators []string `mapstructure:"players"`
}

type CustomActionConfig struct {
	Name        string   `mapstructure:"name"`
	Description string   `mapstructure:"description"`
	Pattern     string   `mapstructure:"pattern"`
	Commands    []string `mapstructure:"commands"`
}

type Config struct {
	MinecraftServer MinecraftServerConfig `mapstructure:"minecraft_server"`
	ManahaMinder    ManahaMinder          `mapstructure:"manaha_minder"`
	Login           LoginConfig           `mapstructure:"login"`
	Activity        ActivityConfig        `mapstructure:"activity"`
	Operator        OperatorConfig        `mapstructure:"operator"`
	CustomActions   []CustomActionConfig  `mapstructure:"custom_actions"`
}

func GetConfig() Config {
	once.Do(func() {
		logger.Debug("Reading config file")
		viper.SetConfigName("mminder")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("$HOME/config")
		viper.AddConfigPath("/etc/")
		viper.AddConfigPath("/usr/local/etc/")
		err := viper.ReadInConfig()
		if err != nil {
			logger.Fatal("Cannot read config file. File may not exist, or be in the wrong format.")
		}
		err = viper.Unmarshal(&config)
		if err != nil {
			logger.Fatal("Cannot read config file. File may be in the wrong format.")
		}
	})

	if logger.GetLevel() == logger.TraceLevel {
		pc, _, _, ok := runtime.Caller(1)
		details := runtime.FuncForPC(pc)
		if ok && details != nil {
			logger.Trace("Returning config to %s", details.Name())
		}
	}

	return config
}
