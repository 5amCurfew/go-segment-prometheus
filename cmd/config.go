package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

//InitConfig init config
func InitConfig() {
	viper.SetDefault("port", "8000")
	viper.AutomaticEnv()
	for key, value := range viper.AllSettings() {
		log.Infof("%s:%s", key, value)
	}
}
