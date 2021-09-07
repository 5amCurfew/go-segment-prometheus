package cmd

// go get ./...
import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Event struct {
	Version    float64                `json:"version"`
	Event      string                 `json:"type"`
	UserId     string                 `json:"UserId"`
	Properties map[string]interface{} `json:"properties"`
	Timestamp  string                 `json:"timestamp"`
}

var (
	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "myapp_processed_ops_total",
		Help: "The total number of processed events",
	})
)

func segmentHandler(c *gin.Context) {
	payload, err := ioutil.ReadAll(c.Request.Body)
	if handleError(c, "unable to read request", err) {
		return // exit
	}
	log.Println(string(payload))

	var e Event
	err = json.Unmarshal(payload, &e)
	if err != nil {
		panic(err)
	}
	log.Printf("%+v", e)

	opsProcessed.Inc()

	c.JSON(200, "success")
}
