//go:generate swagger generate server -t swagger -f ./swagger.yml --exclude-main -A api
//go:generate swagger-codegen generate -i ./swagger.yml -l javascript -o ../frontend/client

package api

import (
	"context"

	"math/rand"
	"sync"
	"time"

	"errors"

	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/jirwin/ipfs-archive/api/swagger/models"
	"github.com/jirwin/ipfs-archive/api/swagger/restapi/operations"
	ipfsOps "github.com/jirwin/ipfs-archive/api/swagger/restapi/operations/ipfs"
	"github.com/jirwin/ipfs-archive/backend"
	"github.com/pborman/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Server struct {
	ctx    context.Context
	logger *zap.Logger

	clients   map[string]backend.BackendServiceClient
	clientMtx sync.RWMutex // manage clients map
	backends  []string
}

func (s *Server) getClient(backend string) (string, backend.BackendServiceClient, error) {
	s.clientMtx.RLock()
	defer s.clientMtx.RUnlock()

	if backend == "" {
		randomBackend := s.backends[rand.Int()%len(s.backends)]
		return randomBackend, s.clients[randomBackend], nil

	}

	if b, ok := s.clients[backend]; ok {
		return backend, b, nil
	}

	return "", nil, errors.New("Unknown backend.")

}

func (s *Server) archiveUrl(params ipfsOps.ArchiveURLParams) middleware.Responder {
	url := swag.StringValue(params.Body.URL)

	primeBackend, client, err := s.getClient("")
	if err != nil {
		return ipfsOps.NewArchiveURLInternalServerError().WithPayload(&models.Error{
			Message: fmt.Sprintf("Error while archiving url. (%s)", err.Error()),
		})
	}
	id := uuid.New()

	scrapeResp, err := client.Scrape(s.ctx, &backend.ScrapeReq{
		Id:  id,
		Url: url,
	})
	if err != nil {
		return ipfsOps.NewArchiveURLInternalServerError().WithPayload(&models.Error{
			Message: fmt.Sprintf("Error while archiving url. (%s)", err.Error()),
		})
	}

	for _, b := range s.backends {
		if b == primeBackend {
			continue
		}
		_, c, err := s.getClient(b)
		if err != nil {
			s.logger.Error("Unable to get client",
				zap.Error(err),
				zap.String("backend", b),
			)
			continue
		}
		go func(b string, c backend.BackendServiceClient) {
			_, err := c.Pin(s.ctx, &backend.PinReq{
				Id:   id,
				Hash: scrapeResp.GetHash(),
			})
			if err != nil {
				s.logger.Error("Unable to pin hash",
					zap.Error(err),
					zap.String("backend", b),
					zap.String("hash", scrapeResp.GetHash()),
				)
				return
			}
		}(b, c)
	}

	return ipfsOps.NewArchiveURLCreated().WithPayload(&models.ArchiveResponse{
		ArchivedURL: scrapeResp.GetArchiveUrl(),
		ID:          scrapeResp.GetId(),
	})
}

func (s *Server) Handle(api *operations.API) {
	api.IPFSArchiveURLHandler = ipfsOps.ArchiveURLHandlerFunc(s.archiveUrl)
}

func NewServer(ctx context.Context, logger *zap.Logger, backends []string) *Server {
	rand.Seed(time.Now().UnixNano())

	server := &Server{
		ctx:      ctx,
		logger:   logger.Named("ipfs-archive-api"),
		clients:  make(map[string]backend.BackendServiceClient),
		backends: backends,
	}

	server.clientMtx.Lock()
	for _, b := range backends {
		conn, err := grpc.Dial(b, grpc.WithInsecure())
		if err != nil {
			panic(err)
		}

		// FIXME: Figure out tls here
		server.clients[b] = backend.NewBackendServiceClient(conn)
	}
	server.clientMtx.Unlock()

	return server
}
