package config

import (
	"flag"
)

func Do() (intervall Intervall, urls Urls, readFromFile bool) {
	flag.BoolVar(&readFromFile, "r", true, "read config file")
	flag.Var(&intervall, "i", "intervall between to polls")
	flag.Var(&urls, "u", "url to poll")
	flag.Parse()
	if intervall == 0 {
		intervall = DefaultIntervall()
	}
	return
}
