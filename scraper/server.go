//go:generate protoc --go_out=plugins=grpc:. ../pb/scraper.proto -I ../pb

package scraper

import (
	"context"
	"errors"
	"net/url"
	"os/exec"
	"path"
	"strings"

	"time"

	"go.uber.org/zap"
	netctx "golang.org/x/net/context"
)

type Server struct {
	logger *zap.Logger
}

func (s *Server) Scrape(ctx netctx.Context, req *ScrapeReq) (*ScrapeResp, error) {
	scrapeCtx, _ := context.WithTimeout(ctx, time.Second*30)
	scraper := NewScraper(scrapeCtx, req.Id, req.Url)

	err := scraper.Scrape()
	if err != nil {
		s.logger.Error("Error scraping.",
			zap.Error(err),
			zap.String("url", req.Url),
		)
		return nil, err
	}

	cmd := exec.Command("ipfs", "add", "-r", path.Join(scraper.SnapshotDir, scraper.Id))
	output, err := cmd.Output()
	if err != nil {
		s.logger.Error("Error running ipfs command", zap.Error(err))
		return nil, err
	}

	outputLines := strings.Split(string(output), "\n")
	parts := strings.Split(outputLines[len(outputLines)-2], " ")
	if len(parts) < 2 {
		return nil, errors.New("IPFS output invalid.")
	}

	archiveUrl := &url.URL{
		Scheme: "https",
		Host:   "ipfs.io",
		Path:   path.Join("ipfs", parts[1]),
	}

	resp := &ScrapeResp{
		Id:         req.Id,
		ArchiveUrl: archiveUrl.String(),
	}
	return resp, nil
}

func NewServer(ctx context.Context, listenAddr string) ScraperServiceServer {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	return &Server{
		logger: logger.Named("scraper_server"),
	}
}
