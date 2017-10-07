package scraper

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"regexp"
	"strings"
)

var styleUrlRegex = regexp.MustCompile(`url\((.*?)\)`)

func trimQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') || (string(s[0]) == "'" && string(s[len(s)-1]) == "'") {
			return s[1 : len(s)-1]
		}
	}

	return strings.TrimSpace(s)
}

type Resource interface {
	Url() string
	Transform(scraper *Scraper, srcUrl *url.URL, reader io.Reader) (io.Reader, error)
}

type index struct {
	url string
}

func (i *index) Url() string {
	return i.url
}

func (i *index) Transform(scraper *Scraper, srcUrl *url.URL, reader io.Reader) (io.Reader, error) {
	buf, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(buf), nil
}

func NewIndex(url string) *index {
	return &index{
		url: url,
	}
}

type asset struct {
	url string
}

func (a *asset) Url() string {
	return a.url
}

func (a *asset) Transform(scraper *Scraper, srcUrl *url.URL, reader io.Reader) (io.Reader, error) {
	buf, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(buf), nil
}

func NewAsset(url string) *asset {
	return &asset{
		url: url,
	}
}

type stylesheet struct {
	url string
}

func (i *stylesheet) Url() string {
	return i.url
}

func (i *stylesheet) Transform(scraper *Scraper, srcUrl *url.URL, reader io.Reader) (io.Reader, error) {
	buf, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	styleUrlRegex.ReplaceAllFunc(buf, func(match []byte) []byte {
		urlMatch := string(match[4 : len(match)-1])
		rUrl := trimQuotes(urlMatch)
		if strings.HasPrefix(rUrl, "data:") {
			return match
		}

		asset := NewAsset(rUrl)

		absUrl, filename, err := scraper.ensureFilename(asset, srcUrl)
		if err != nil {
			return match
		}

		scraper.queueResource(NewAsset(absUrl.String()))

		return []byte(fmt.Sprintf("url(%s)", filename))
	})

	return bytes.NewReader(buf), nil
}

func NewStylesheet(url string) *stylesheet {
	return &stylesheet{
		url: url,
	}
}

type script struct {
	url string
}

func (s *script) Url() string {
	return s.url
}

func (s *script) Transform(scraper *Scraper, srcUrl *url.URL, reader io.Reader) (io.Reader, error) {
	buf, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(buf), nil
}

func NewScript(url string) *script {
	return &script{
		url: url,
	}
}

type image struct {
	url string
}

func (i *image) Url() string {
	return i.url
}

func (i *image) Transform(scraper *Scraper, srcUrl *url.URL, reader io.Reader) (io.Reader, error) {
	buf, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(buf), nil
}

func NewImage(url string) *image {
	return &image{
		url: url,
	}
}
