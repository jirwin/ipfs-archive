package scraper

import (
	"bytes"
	"io"
	"io/ioutil"
	"regexp"
)

var styleUrlRegex = regexp.MustCompile(`url\((.*?)\)`)

type Resource interface {
	Url() string
	Transform(reader io.Reader) (io.Reader, error)
}

type index struct {
	url string
}

func (i *index) Url() string {
	return i.url
}

func (i *index) Transform(reader io.Reader) (io.Reader, error) {
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

type stylesheet struct {
	url string
}

func (i *stylesheet) Url() string {
	return i.url
}

func (i *stylesheet) Transform(reader io.Reader) (io.Reader, error) {
	buf, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	styleUrlRegex.ReplaceAllFunc(buf, func(match []byte) []byte {
		//fmt.Printf("%s\n", match)
		return match
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

func (s *script) Transform(reader io.Reader) (io.Reader, error) {
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

func (i *image) Transform(reader io.Reader) (io.Reader, error) {
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
