package main

import (
	"go-segment-prometheus/cmd"
)

func main() {
	cmd.InitConfig()
	cmd.InitHTTPServer()
}
