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
		cli.StringFlag{
			Name:  "address,a",
			Usage: "The address to listen on",
		},

		cli.StringSliceFlag{
			Name:  "backend,b",
			Usage: "The ipfs-archive backends to connect to.",
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

	if !cliCtx.IsSet("address") {
		cli.ShowAppHelpAndExit(cliCtx, -1)
	}

	if !cliCtx.IsSet("backend") {
		cli.ShowAppHelpAndExit(cliCtx, -1)
	}

	// load embedded swagger file
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		logger.Error("Error loading swaggerApi", zap.Error(err))
		return cli.NewExitError("Unable to start api.", -1)
	}

	// create new service API
	swaggerApi := operations.NewAPI(swaggerSpec)
	server := restapi.NewServer(swaggerApi)

	handler := api.NewServer(ctx, logger, cliCtx.StringSlice("backend"))
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
