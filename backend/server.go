//go:generate protoc --gogoslick_out=plugins=grpc:. ../pb/backend.proto -I ../pb

package backend

import (
	"context"
	"net/url"
	"os/exec"
	"path"
	"strings"

	"time"

	"github.com/jirwin/ipfs-archive/scraper"
	"go.uber.org/zap"
	netctx "golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BackendServer struct {
	ctx    context.Context
	logger *zap.Logger
}

func (b *BackendServer) runIpfsCmd(cmdArgs []string) ([]byte, error) {
	cmd := exec.Command("ipfs", cmdArgs...)
	output, err := cmd.Output()
	if err != nil {
		b.logger.Error("Error running ipfs command", zap.Error(err))
		return nil, err
	}

	return output, nil
}

func (b *BackendServer) Scrape(ctx netctx.Context, req *ScrapeReq) (*ScrapeResp, error) {
	scrapeCtx, _ := context.WithTimeout(ctx, time.Second*30)
	sc := scraper.NewScraper(scrapeCtx, b.logger, req.Id, req.Url)

	err := sc.Scrape()
	if err != nil {
		b.logger.Error("Error scraping.",
			zap.Error(err),
			zap.String("url", req.Url),
		)
		return nil, status.Errorf(codes.Internal, "Unable to scrape %s: %s", req.Url, err.Error())
	}

	output, err := b.runIpfsCmd([]string{"add", "-r", path.Join(sc.SnapshotDir, sc.Id)})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error scraping url: %ss", req.GetUrl())
	}

	outputLines := strings.Split(string(output), "\n")
	parts := strings.Split(outputLines[len(outputLines)-2], " ")
	if len(parts) < 2 {
		return nil, status.Error(codes.Internal, "Invalid IPFS output")
	}

	hash := parts[1]

	archiveUrl := &url.URL{
		Scheme: "https",
		Host:   "gateway.archive.network",
		Path:   path.Join("ipfs", hash),
	}

	resp := &ScrapeResp{
		Id:         req.Id,
		ArchiveUrl: archiveUrl.String(),
		Hash:       hash,
	}

	return resp, nil
}

func (b *BackendServer) Pin(ctx netctx.Context, req *PinReq) (*PinResp, error) {
	hash := req.GetHash()

	if hash == "" {
		return nil, status.Error(codes.InvalidArgument, "You must provide a hash to pin.")
	}

	_, err := b.runIpfsCmd([]string{"pin", "add", hash})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error pinning hash: %s", hash)
	}

	return &PinResp{
		Id: req.Id,
	}, nil
}

func NewServer(ctx context.Context, logger *zap.Logger) BackendServiceServer {
	return &BackendServer{
		ctx:    ctx,
		logger: logger,
	}
}
