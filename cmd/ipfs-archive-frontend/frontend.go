package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/rakyll/statik/fs"
	"gopkg.in/urfave/cli.v1"

	"context"
	"time"

	"fmt"

	_ "github.com/jirwin/ipfs-archive/frontend/statik"
	"github.com/jirwin/ipfs-archive/version"
)

func main() {
	app := cli.NewApp()
	app.Name = "ipfs-archive-frontend"
	app.Usage = "ipfs-archive frontend service"
	app.Version = version.Version
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "address,a",
			Usage: "The address to listen on",
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

	statikFS, err := fs.New()
	if err != nil {
		logger.Error(err.Error())
		return cli.NewExitError(err.Error(), -1)
	}

	stop := make(chan os.Signal, 2)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(statikFS))

	addr := cliCtx.String("address")

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	logger.Info(fmt.Sprintf("Listening on %s", addr))

	go func() {
		if err := server.ListenAndServe(); err != nil {
			logger.Fatal(err.Error(), zap.Error(err))
		}
	}()

	<-stop

	logger.Info("Shutting down server.")
	shutdownCtx, _ := context.WithTimeout(ctx, time.Second*10)
	server.Shutdown(shutdownCtx)
	logger.Info("Server shutdown.")

	return nil
}
