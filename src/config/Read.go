package config

import (
	"flag"
)

var configFile = "config.json"

type Config struct {
	Intervall Intervall
	Urls      Urls
}

func Default() {
	return Config{
		DefaultIntervall(),
		nil
	}
}

func FromCmdl(config* Config) (fromFile bool) {
	config = Default()
	flag.BoolVar(&fromFile, "r", true, "read config file")
	flag.StringVar(&configFile, "c", "config file")
	flag.Var(&config.Intervall, "i", "intervall between to polls")
	flag.Var(&config.Urls, "u", "url to poll")
	flag.Parse()
	if config.Intervall == 0 {
		config.Intervall = DefaultIntervall()
	}
	return
}

func FromFile(config* Config) error {
	file, err := os.Open(configFile)
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return err
	}
	return nil
}

