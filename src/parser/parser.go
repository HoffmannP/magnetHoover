package parser

import (
	"log"
	"strings"
)

type ParserRawFunc func(string) ([][]string, error)
type _ struct{} // only for syntax highlighting in GoSublime
type ParserFunc func() ([][]string, error)
type _ struct{} // only for syntax highlighting in GoSublime

type UrlFunc func(string) (string, error)

type parserPlugin struct {
	p ParserRawFunc
	u UrlFunc
}

var parsers map[string]parserPlugin

func register(n string, p ParserRawFunc, u UrlFunc) {
	if parsers == nil {
		parsers = make(map[string]parserPlugin, 10)
	}
	parsers[n] = parserPlugin{p, u}
}

func Parser(id string) ParserFunc {
	// Which Parser
	name := "Default"
	parts := strings.Split(id, "§")
	if len(parts) > 1 {
		name = parts[0]
		id = parts[1]
	}
	pl, ok := parsers[name]
	if !ok {
		pl = parsers["Default"]
	}

	// Extract uri from page/feed id
	uri, err := pl.u(id)
	if err != nil {
		log.Println(err)
	}
	return func() ([][]string, error) {
		return pl.p(uri)
	}
}

func init() {
	register("Default",
		func(body string) ([][]string, error) { return nil, nil },
		func(uri string) (string, error) { return uri, nil })
}
