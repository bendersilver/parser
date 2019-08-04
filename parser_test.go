package parser

import (
	"strings"
	"testing"
)

func TestParser(t *testing.T) {
	page, err := NewPage("https://github.com/bendersilver/parser")
	if err != nil {
		t.Fatal(err)
	}
	if page.Text(".author") != "bendersilver" {
		t.Error("Invalid value Text")
	}
	if page.Attr(".author a", "href") != "/bendersilver" {
		t.Error("Invalid value Attr")
	}
	elem := page.AllNode(".js-navigation-item")[0]
	if !strings.Contains(elem.RootAttr("class"), "js-navigation-item") {
		t.Error("Invalid value RootAttr")
	}
}
