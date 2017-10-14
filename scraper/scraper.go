package scraper

import (
	"context"
	"io"
	"net/http"

	"os"

	"net/url"
	"path"

	"sync"

	"compress/gzip"
	"time"

	"compress/zlib"

	"crypto/tls"

	"io/ioutil"

	"fmt"

	"github.com/PuerkitoBio/goquery"
	"github.com/jirwin/ipfs-archive/version"
	"go.uber.org/zap"
)

type Scraper struct {
	ctx           context.Context
	cancel        context.CancelFunc
	wg            sync.WaitGroup
	seed          string
	baseUrl       *url.URL
	resourceQueue chan Resource
	SnapshotDir   string
	Log           *zap.Logger
	Id            string
	client        *http.Client
}

func (s *Scraper) request(url string) (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept-Encoding", "gzip")
	req.Header.Set("User-Agent", fmt.Sprintf("ipfs-archive/%s", version.Version))

	resp, err := s.client.Do(req)
	if err != nil {
		s.Log.Error("Error grabbing url",
			zap.Error(err),
			zap.String("url", url),
		)
		return nil, err
	}

	var respReader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		respReader, err = gzip.NewReader(resp.Body)
		if err != nil {
			s.Log.Error("Error gunzipping response",
				zap.Error(err),
				zap.String("url", url),
			)
			return nil, err
		}
	case "deflate":
		respReader, err = zlib.NewReader(resp.Body)
		if err != nil {
			s.Log.Error("Error deflating response",
				zap.Error(err),
				zap.String("url", url),
			)
			return nil, err
		}
	default:
		respReader = resp.Body
	}

	return respReader, nil
}

func (s *Scraper) Scrape() error {
	reader, err := s.request(s.seed)
	defer reader.Close()
	if err != nil {
		s.Log.Error("Unable to fetch seed",
			zap.Error(err),
			zap.String("url", s.seed),
		)
		return err
	}

	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		s.Log.Error("Unable to scrape seed",
			zap.Error(err),
			zap.String("seed", s.seed),
		)
		return err
	}

	go s.process()

	doc.Find("img").Each(func(i int, sel *goquery.Selection) {
		img, exists := sel.Attr("src")
		if !exists {
			return
		}

		sel.SetAttr("src", s.rewriteAssetUrl(img))

		s.queueResource(NewImage(img))
	})

	doc.Find("link").Each(func(i int, sel *goquery.Selection) {
		link, exists := sel.Attr("href")
		if !exists {
			return
		}
		rel, exists := sel.Attr("rel")
		if !exists || rel != "stylesheet" {
			return
		}

		sel.SetAttr("href", s.rewriteAssetUrl(link))
		sel.RemoveAttr("crossorigin")

		s.queueResource(NewStylesheet(link))
	})

	doc.Find("script").Each(func(i int, sel *goquery.Selection) {
		scriptSrc, exists := sel.Attr("src")
		if !exists {
			return
		}

		sel.SetAttr("src", s.rewriteAssetUrl(scriptSrc))
		sel.RemoveAttr("crossorigin")

		s.queueResource(NewScript(scriptSrc))
	})

	s.wg.Wait()

	_, indexFilename, err := s.ensureFilename(&index{
		url: s.seed,
	}, nil)

	out, err := os.Create(indexFilename)
	defer out.Close()
	if err != nil {
		fmt.Println(s.seed)
		s.Log.Error("Error creating index file",
			zap.Error(err),
			zap.String("filename", indexFilename),
			zap.String("url", s.seed),
		)
		return err
	}

	outIndex, err := doc.Html()
	if err != nil {
		return err
	}

	out.WriteString(outIndex)
	out.Sync()

	return nil
}

func (s *Scraper) rewriteAssetUrl(rawurl string) string {
	sUrl, err := s.toAbsUrl(rawurl, nil)
	if err != nil {
		return ""
	}

	pUrl, err := url.Parse(sUrl)

	return (&url.URL{
		Path: path.Join(pUrl.Hostname(), pUrl.Path),
	}).String()
}

func (s *Scraper) process() error {
	for {
		select {
		case resource := <-s.resourceQueue:
			go s.fetch(resource)

		case <-s.ctx.Done():
			s.Log.Info("Context done, stopping fetching.")
			return nil

		}
	}
	return nil
}

func (s *Scraper) toAbsUrl(rawurl string, baseUrl *url.URL) (string, error) {
	relurl, err := url.Parse(rawurl)
	if err != nil {
		return "", err
	}

	if relurl.IsAbs() {
		return relurl.String(), nil
	}

	if baseUrl == nil {
		baseUrl = s.baseUrl
	}

	base, err := url.Parse((&url.URL{
		Scheme: baseUrl.Scheme,
		Host:   baseUrl.Host,
		Path:   baseUrl.Path,
	}).String())
	if err != nil {
		return "", err
	}

	absurl := base.ResolveReference(relurl)
	return absurl.String(), nil
}

func (s *Scraper) Cleanup() error {
	err := os.RemoveAll(s.SnapshotDir)
	if err != nil {
		return err
	}

	return nil
}

func (s *Scraper) ensureFilename(resource Resource, baseUrl *url.URL) (*url.URL, string, error) {
	rUrl, err := s.toAbsUrl(resource.Url(), baseUrl)
	if err != nil {
		return nil, "", err
	}

	parsed, err := url.Parse(rUrl)
	if err != nil {
		return nil, "", err
	}

	rootFileparts := []string{s.SnapshotDir, s.Id}
	var dirPath string

	filepath, filename := path.Split(parsed.Path)

	switch resource.(type) {
	case *index:
		dirPath = path.Join(rootFileparts...)
		filename = "index.html"
	default:
		dirPath = path.Join(append(rootFileparts, parsed.Hostname(), filepath)...)
	}

	err = os.MkdirAll(dirPath, 0755)
	if err != nil {
		return nil, "", err
	}

	return parsed, path.Join(dirPath, filename), nil
}

func (s *Scraper) fetch(resource Resource) {
	defer s.wg.Done()

	if resource.Url() == "" {
		return
	}

	realUrl, filename, err := s.ensureFilename(resource, nil)
	if err != nil {
		s.Log.Error("Error ensuring filename",
			zap.Error(err),
			zap.String("url", resource.Url()),
		)
		return
	}

	out, err := os.Create(filename)
	defer out.Close()
	if err != nil {
		s.Log.Error("Error creating resource file",
			zap.Error(err),
			zap.String("filename", filename),
			zap.String("url", realUrl.String()),
			zap.String("resource url", resource.Url()),
		)
		return
	}

	respReader, err := s.request(realUrl.String())
	if err != nil {
		s.Log.Error("Unable to fetch resource",
			zap.Error(err),
			zap.String("url", realUrl.String()),
		)
		return
	}

	transformedResp, err := resource.Transform(s, realUrl, respReader)
	if err != nil {
		s.Log.Error("Error transforming resource", zap.Error(err))
		return
	}

	_, err = io.Copy(out, transformedResp)
	if err != nil {
		s.Log.Error("Error writing file",
			zap.Error(err),
			zap.String("url", resource.Url()),
			zap.String("filename", filename),
		)
		return
	}
}

func (s *Scraper) queueResource(resource Resource) {
	s.wg.Add(1)
	s.resourceQueue <- resource
}

func NewScraper(ctx context.Context, logger *zap.Logger, id, seed string) *Scraper {
	ctx, canc := context.WithCancel(ctx)

	parsedUrl, err := url.Parse(seed)
	if err != nil {
		panic(err)
	}

	snapshotDir, err := ioutil.TempDir("", "ipfs-archive")
	if err != nil {
		panic(err)
	}

	scraper := &Scraper{
		ctx:           ctx,
		cancel:        canc,
		Log:           logger.Named(fmt.Sprintf("scraper-%s", id)),
		seed:          seed,
		baseUrl:       parsedUrl,
		resourceQueue: make(chan Resource, 5),
		Id:            id,
		SnapshotDir:   snapshotDir,
		client: &http.Client{
			Timeout: time.Second * 10,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
	}

	return scraper
}
