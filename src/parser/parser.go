package parser

type ParserFunc func(string) ([][]string, error)
type parserList map[string]ParserFunc

var parsers map[string]ParserFunc

func register(n string, p ParserFunc) {
	if parsers == nil {
		parsers = make(parserList, 10)
	}
	parsers[n] = p
}

func defaultParserFunc(body string) ([][]string, error) {
	return nil, nil
}

func Parser(p string) ParserFunc {
	p, ok := parsers[p]
	if !ok {
		return parsers["Default"]
	}
	return p
}

func init() {
	register("Default", defaultParserFunc)
}
