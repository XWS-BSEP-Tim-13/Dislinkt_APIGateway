package main

import (
	"context"
	logger "github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/logging"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/startup"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/startup/config"
)

func main() {
	config := config.NewConfig()
	logger := logger.InitLogger("api-gateway", context.TODO())
	server := startup.NewServer(config, logger)
	server.Start(logger)

	defer server.CloseTracer()
}
