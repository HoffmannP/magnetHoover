package plugin

type ParserFunc func(string) ([][]string, error)
type parserList map[string]ParserFunc

var parsers map[string]ParserFunc

func register(n string, p ParserFunc) {
	if parsers == nil {
		parsers = make(parserList, 10)
	}
	parsers[n] = p
}

func Default(body string) ([][]string, error) {
	return nil, nil
}

func Parser(n string) ParserFunc {
	p, ok := parsers[n]
	if !ok {
		return parsers["Default"]
	}
	return p
}

func init() {
	register("Default", Default)
}
