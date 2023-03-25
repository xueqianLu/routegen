package utils

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"github.com/xueqianLu/routegen/config"
	"github.com/xueqianLu/routegen/log"
	"os"
)

func Exists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

// initConfig reads in config file and ENV variables if set.
func InitConfig() {
	log.InitLog()

	viper.SetConfigName("config") // name of config file (without extension)
	viper.AddConfigPath(".")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Info("Using config file:", viper.ConfigFileUsed())
	} else {
		log.Error("Read config failed", "error", err)
		return
	}

	_, err := config.ParseConfig(viper.ConfigFileUsed())
	if err != nil {
		log.WithField("error", err).Fatal("parse config failed")
	}
	//log.Info("config is", config.GetConfig())

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Warn("Config file changed:", e.Name)
	})
}
