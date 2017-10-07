package main

import (
	"fmt"
	"os"

	"context"

	"github.com/jirwin/ipfs-archive/scraper"
	"go.uber.org/zap"
	cli "gopkg.in/urfave/cli.v1"
)

func main() {
	app := cli.NewApp()
	app.Name = "ipfs-archive"
	app.Usage = "Use to archive url to ipfs"
	app.Commands = []cli.Command{
		cli.Command{
			Name:   "archive",
			Usage:  "Archive a url to ipfs",
			Action: run,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "url",
					Usage: "Archive `URL` to ipfs",
				},
			},
		},
	}

	app.Run(os.Args)
}

func run(cliCtx *cli.Context) error {
	ctx := context.Background()

	logger, err := zap.NewProduction()
	if err != nil {
		return cli.NewExitError(err.Error(), -1)
	}
	ctx = context.WithValue(ctx, "logger", logger)

	if !cliCtx.IsSet("url") {
		cli.ShowCommandHelpAndExit(cliCtx, "archive", -1)
	}

	url := cliCtx.String("url")

	scraper := scraper.NewScraper(ctx, url)
	err = scraper.Scrape()
	if err != nil {
		return cli.NewExitError(err.Error(), -1)
	}

	fmt.Printf("Finished scraping %s to %s\n", scraper.Id, url)

	return nil
}
