package main

import (
	"context"
	"os"
	"os/signal"

	loads "github.com/go-openapi/loads"
	"go.uber.org/zap"
	"gopkg.in/urfave/cli.v1"

	"github.com/jirwin/ipfs-archive/api"
	"github.com/jirwin/ipfs-archive/api/swagger/restapi"
	"github.com/jirwin/ipfs-archive/api/swagger/restapi/operations"
	"github.com/jirwin/ipfs-archive/version"
)

func main() {
	app := cli.NewApp()
	app.Name = "ipfs-archive-api"
	app.Usage = "IPFS api service"
	app.Version = version.Version
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "port,p",
			Usage: "The port to listen on",
			Value: 7002,
		},
	}
	app.Action = run

	app.Run(os.Args)
}

func run(cliCtx *cli.Context) error {
	ctx := context.Background()
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	// load embedded swagger file
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		logger.Error("Error loading swaggerApi", zap.Error(err))
		return cli.NewExitError("", -1)
	}

	// create new service API
	swaggerApi := operations.NewAPI(swaggerSpec)
	server := restapi.NewServer(swaggerApi)

	handler := api.NewServer(ctx, logger, []string{"localhost:7001"})
	handler.Handle(swaggerApi)

	server.SetAPI(swaggerApi)

	server.Port = cliCtx.Int("port")
	go func() {
		if err := server.Serve(); err != nil {
			logger.Fatal(err.Error(), zap.Error(err))
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop

	logger.Info("Shutting down the server...")
	server.Shutdown()
	logger.Info("Server shutdown")

	return nil
}
