package scraper

import (
	"context"
	"io"
	"net/http"

	"os"

	"net/url"
	"path"

	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/pborman/uuid"
	"go.uber.org/zap"
)

type Resource interface {
	Url() string
	Index() bool
}

type index struct {
	url string
}

func (i *index) Url() string {
	return i.url
}

func (i *index) Index() bool {
	return true
}

type asset struct {
	url string
}

func (a *asset) Url() string {
	return a.url
}

func (a *asset) Index() bool {
	return false
}

type Scraper struct {
	ctx           context.Context
	cancel        context.CancelFunc
	wg            sync.WaitGroup
	seed          string
	baseUrl       *url.URL
	resourceQueue chan Resource
	Log           *zap.Logger
	rootDir       string
}

func (s *Scraper) Scrape() error {
	doc, err := goquery.NewDocument(s.seed)
	if err != nil {

		return err
	}

	go s.process()

	doc.Find("img").Each(func(i int, sel *goquery.Selection) {
		img, exists := sel.Attr("src")
		if !exists {
			s.Log.Error("img tag with no src",
				zap.String("tag", sel.Text()),
			)
			return
		}

		sel.SetAttr("src", s.rewriteAssetUrl(img))

		s.wg.Add(1)
		s.resourceQueue <- &asset{
			url: img,
		}
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

		s.wg.Add(1)
		s.resourceQueue <- &asset{
			url: link,
		}
	})

	doc.Find("script").Each(func(i int, sel *goquery.Selection) {
		scriptSrc, exists := sel.Attr("src")
		if !exists {
			s.Log.Error("script tag with no src",
				zap.String("tag", sel.Text()),
			)
			return
		}

		sel.SetAttr("src", s.rewriteAssetUrl(scriptSrc))
		sel.RemoveAttr("crossorigin")

		s.wg.Add(1)
		s.resourceQueue <- &asset{
			url: scriptSrc,
		}
	})

	s.wg.Wait()

	indexFilename, err := s.ensureFilename(&index{
		url: s.seed,
	})

	out, err := os.Create(indexFilename)
	defer out.Close()
	if err != nil {
		s.Log.Error("Error creating file",
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
	sUrl, err := s.toAbsUrl(rawurl)
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
			s.fetch(resource)

		case <-s.ctx.Done():
			s.Log.Info("Context done, stopping fetching.")
			return nil
		}
	}
	return nil
}

func (s *Scraper) toAbsUrl(rawurl string) (string, error) {
	relurl, err := url.Parse(rawurl)
	if err != nil {
		return "", err
	}

	if relurl.IsAbs() {
		return relurl.String(), nil
	}

	base, err := url.Parse((&url.URL{
		Scheme: s.baseUrl.Scheme,
		Host:   s.baseUrl.Host,
	}).String())
	if err != nil {
		return "", err
	}

	absurl := base.ResolveReference(relurl)
	return absurl.String(), nil
}

func (s *Scraper) ensureFilename(resource Resource) (string, error) {
	rUrl, err := s.toAbsUrl(resource.Url())
	if err != nil {
		return "", err
	}

	parsed, err := url.Parse(rUrl)
	if err != nil {
		return "", err
	}

	rootFileparts := []string{"snapshots", s.rootDir}
	var dirPath string

	filepath, filename := path.Split(parsed.Path)

	if resource.Index() {
		dirPath = path.Join(rootFileparts...)
	} else {
		dirPath = path.Join(append(rootFileparts, parsed.Hostname(), filepath)...)
	}

	err = os.MkdirAll(dirPath, 0755)
	if err != nil {
		return "", err
	}

	if resource.Index() {
		return path.Join(dirPath, "index.html"), nil
	}

	return path.Join(dirPath, filename), nil
}

func (s *Scraper) fetch(resource Resource) {
	defer s.wg.Done()

	filename, err := s.ensureFilename(resource)
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
		s.Log.Error("Error creating file",
			zap.Error(err),
			zap.String("filename", filename),
			zap.String("url", resource.Url()),
		)
		return
	}

	resp, err := http.Get(resource.Url())
	defer resp.Body.Close()
	if err != nil {
		s.Log.Error("Error grabbing url",
			zap.Error(err),
			zap.String("url", resource.Url()),
		)
		return
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		s.Log.Error("Error writing file",
			zap.Error(err),
			zap.String("url", resource.Url()),
			zap.String("filename", filename),
		)
		return
	}
}

func NewScraper(ctx context.Context, seed string) *Scraper {
	logger, ok := ctx.Value("logger").(*zap.Logger)
	if !ok {
		panic("not a logger")
	}
	ctx, canc := context.WithCancel(ctx)

	parsedUrl, err := url.Parse(seed)
	if err != nil {
		panic(err)
	}

	scraper := &Scraper{
		ctx:           ctx,
		cancel:        canc,
		Log:           logger.Named("scraper"),
		seed:          seed,
		baseUrl:       parsedUrl,
		resourceQueue: make(chan Resource, 5),
		rootDir:       uuid.NewUUID().String(),
	}

	return scraper
}
