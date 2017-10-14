package main

import (
	"fmt"
	"os"

	"context"

	"os/exec"
	"path"

	"net/url"
	"strings"

	"time"

	"github.com/briandowns/spinner"
	"github.com/jirwin/ipfs-archive/scraper"
	"github.com/jirwin/ipfs-archive/version"
	"github.com/pborman/uuid"
	"go.uber.org/zap"
	cli "gopkg.in/urfave/cli.v1"
)

func main() {
	app := cli.NewApp()
	app.Name = "ipfs-archive"
	app.Usage = "Use to archive url to ipfs"
	app.Version = version.Version
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "url",
			Usage: "Archive `URL` to ipfs",
		},
	}
	app.Action = run

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
		cli.ShowAppHelpAndExit(cliCtx, -1)
	}

	seedUrl := cliCtx.String("url")

	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.FinalMSG = fmt.Sprintf("Fetched %s\n", seedUrl)
	s.Prefix = fmt.Sprintf("Fetching %s ", seedUrl)
	s.Writer = os.Stderr
	s.Start()

	scraper := scraper.NewScraper(ctx, logger, uuid.New(), seedUrl)
	err = scraper.Scrape()
	if err != nil {
		return cli.NewExitError(err.Error(), -1)
	}
	s.Stop()

	s = spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.FinalMSG = fmt.Sprintf("Archived %s\n", seedUrl)
	s.Prefix = fmt.Sprintf("Archiving %s ", seedUrl)
	s.Start()
	cmd := exec.Command("ipfs", "add", "-r", path.Join(scraper.SnapshotDir, scraper.Id))
	output, err := cmd.Output()
	if err != nil {
		logger.Error("Error running ipfs command", zap.Error(err))
		return cli.NewExitError("", -1)
	}

	outputLines := strings.Split(string(output), "\n")
	parts := strings.Split(outputLines[len(outputLines)-2], " ")

	if len(parts) < 2 {
		return cli.NewExitError("Unexpected ipfs output.", -1)
	}

	archiveUrl := &url.URL{
		Scheme: "https",
		Host:   "ipfs.io",
		Path:   path.Join("ipfs", parts[1]),
	}

	s.Stop()

	fmt.Printf("%s is archived at %s\n", seedUrl, archiveUrl.String())

	err = scraper.Cleanup()
	if err != nil {
		return cli.NewExitError(err.Error(), -1)
	}

	return nil
}
