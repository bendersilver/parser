package parser

import (
	"net/http"
	"net/http/cookiejar"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html/charset"
)

type page struct {
	Doc *goquery.Selection
}

// Parser -
type Parser interface {
	setDoc(*goquery.Selection)
	Text(string) string
	AllNode(string) []Parser
	AllText(string) []string
	RootAttr(string) string
	Attr(string, string) string
	AllAttr(string, string) []string
}

// setDoc -
func (p *page) setDoc(doc *goquery.Selection) {
	p.Doc = doc
}

// AllNode -
func (p *page) AllNode(path string) (vals []Parser) {
	p.Doc.Find(path).Each(func(i int, s *goquery.Selection) {
		vals = append(vals, Parser(&page{s}))
	})
	return
}

// Text -
func (p *page) Text(path string) (val string) {
	arr := p.AllText(path)
	if len(arr) > 0 {
		val = arr[0]
	}
	return
}

// AllText -
func (p *page) AllText(path string) (vals []string) {
	p.Doc.Find(path).Each(func(i int, s *goquery.Selection) {
		vals = append(vals, strings.TrimSpace(s.Text()))
	})
	return
}

// Attr -
func (p *page) RootAttr(attr string) (val string) {
	if val, ok := p.Doc.Attr(attr); ok {
		return val
	}
	return
}

// Attr -
func (p *page) Attr(path, attr string) (val string) {
	arr := p.AllAttr(path, attr)
	if len(arr) > 0 {
		val = arr[0]
	}
	return
}

// AllAttr -
func (p *page) AllAttr(path, attr string) (vals []string) {
	p.Doc.Find(path).Each(func(i int, s *goquery.Selection) {
		if val, ok := s.Attr(attr); ok {
			vals = append(vals, strings.TrimSpace(val))
		}
	})
	return
}

// NewPage -
func NewPage(url string) (Parser, error) {
	parser := Parser(new(page))
	options := cookiejar.Options{
		// PublicSuffixList: publicsuffix.List,
	}
	jar, err := cookiejar.New(&options)
	if err != nil {
		return parser, err
	}
	client := &http.Client{Jar: jar}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return parser, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return parser, err
	}
	defer resp.Body.Close()
	b, err := charset.NewReader(resp.Body, resp.Header.Get("Content-Type"))
	if err != nil {
		return parser, err
	}
	doc, err := goquery.NewDocumentFromReader(b)

	doc.Find("body").Each(func(i int, s *goquery.Selection) {
		parser.setDoc(s)
	})
	return parser, err
}
