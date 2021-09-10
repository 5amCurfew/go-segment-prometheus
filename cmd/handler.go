package cmd

import (
	"fmt"

	"github.com/Depado/ginprom"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func handleError(c *gin.Context, message string, err error) bool {
	if err != nil {
		logMessage := fmt.Sprintf("reason: %s %v", message, err.Error())
		log.Error(logMessage)
		c.Error(err)
		c.JSON(500, gin.H{"reason": message, "error": err.Error()})
		return true // signal that there was an error and the caller should return
	}
	return false // no error, can continue
}

//InitHTTPServer init server
func InitHTTPServer() {
	router := gin.New()
	p := ginprom.New(
		ginprom.Engine(router),
		ginprom.Subsystem("gin"),
		ginprom.Path("/metrics"),
	)
	router.Use(p.Instrument())
	router.Use(gin.Recovery())
	router.GET("/ping", func(c *gin.Context) { c.JSON(200, gin.H{"message": "pong!"}) })
	router.POST("/monitor_ui", func(c *gin.Context) { segmentHandler(c) })
	router.Run(fmt.Sprintf(":%v", viper.GetInt("port")))
}
