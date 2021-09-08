package cmd

// go get ./...
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Event struct {
	AnonymousID string `json:"anonymousId"`
	Context     struct {
		IP      string `json:"ip"`
		Library struct {
			Name    string `json:"name"`
			Version string `json:"version"`
		} `json:"library"`
		Locale string `json:"locale"`
		Page   struct {
			Path     string `json:"path"`
			Referrer string `json:"referrer"`
			Search   string `json:"search"`
			Title    string `json:"title"`
			URL      string `json:"url"`
		} `json:"page"`
		SignalDeviceID string `json:"signalDeviceId"`
		UserAgent      string `json:"userAgent"`
	} `json:"context"`
	Integrations struct {
	} `json:"integrations"`
	MessageID         string    `json:"messageId"`
	OriginalTimestamp time.Time `json:"originalTimestamp"`
	Properties        struct {
		App struct {
			Version string `json:"version"`
		} `json:"app"`
		Environment string `json:"environment"`
		Path        string `json:"path"`
		Referrer    string `json:"referrer"`
		Search      string `json:"search"`
		Title       string `json:"title"`
		URL         string `json:"url"`
		UserAgent   string `json:"userAgent"`
		UserID      string `json:"userId"`
		WorkspaceID string `json:"workspaceId"`
	} `json:"properties"`
	ReceivedAt time.Time `json:"receivedAt"`
	SentAt     time.Time `json:"sentAt"`
	Timestamp  time.Time `json:"timestamp"`
	Type       string    `json:"type"`
	UserID     string    `json:"userId"`
}

var (
	totalPageCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "total_page_views",
		Help: "The total number of page views recorded by segment.io",
	})

	searchPageCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "total_search_page_views",
		Help: "The total number of search page views recorded by segment.io",
	})

	insightsPageCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "total_insights_page_views",
		Help: "The total number of insight page views recorded by segment.io",
	})

	dashboardPageCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "total_dashboard_page_views",
		Help: "The total number of dashboard page views recorded by segment.io",
	})
)

func segmentHandler(c *gin.Context) {
	payload, err := ioutil.ReadAll(c.Request.Body)
	if handleError(c, "unable to read request", err) {
		return // exit
	}

	var e Event
	err = json.Unmarshal(payload, &e)
	if err != nil {
		panic(err)
	}

	switch e.Type {
	case "page":
		totalPageCounter.Inc()
		if strings.Contains(e.Properties.Path, "search") {
			searchPageCounter.Inc()
			if strings.Contains(e.Properties.Path, "insights") {
				insightsPageCounter.Inc()
			}
		}
		if strings.Contains(e.Properties.Path, "dashboards") {
			dashboardPageCounter.Inc()
		}
		c.JSON(200, "success")
	default:
		message := fmt.Sprintf("Event %s ignored\n", e.Type)
		c.JSON(200, message)
		return
	}
}
