package main

import (
	"context"
	"net"
	"os"

	"fmt"

	"github.com/jirwin/ipfs-archive/backend"
	"github.com/jirwin/ipfs-archive/version"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	cli "gopkg.in/urfave/cli.v1"
)

func main() {
	app := cli.NewApp()
	app.Name = "ipfs-archive-backend"
	app.Usage = "IPFS backend service"
	app.Version = version.Version
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "address,a",
			Usage: "The address to listen on",
			Value: "0.0.0.0:7001",
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

	backendServer := backend.NewServer(ctx, logger.Named("ipfs-backend-server"))
	listenAddr := cliCtx.String("address")

	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		logger.Fatal("failed to listen.", zap.Error(err))
		return cli.NewExitError("Failed to listen.", -1)
	}
	logger.Info(fmt.Sprintf("ipfs-scraper listening on %s", listenAddr))

	grpcServer := grpc.NewServer()
	backend.RegisterBackendServiceServer(grpcServer, backendServer)
	grpcServer.Serve(lis)

	return nil
}
