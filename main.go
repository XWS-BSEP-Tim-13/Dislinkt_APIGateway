package main

import (
	"github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/startup"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/startup/config"
)

func main() {
	config := config.NewConfig()
	server := startup.NewServer(config)
	server.Start()

}
