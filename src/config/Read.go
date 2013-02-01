package config

import (
	"os"
	"encoding/json"
	"flag"
	"time"
	"strings"
)

var configFile = "config.json"

type Config struct {
	Intervall time.Duration
	Urls      []string
}

func (c *Config) Tic() <-chan time.Time {
	return time.After(c.Intervall)
}

func FromCmdl(config* Config) error {
	var notFromFile bool
	var url string
	
	flag.BoolVar(&notFromFile, "n", false, "don't read config file")
	flag.BoolVar(&notFromFile, "noread", false, "don't read config file")
	flag.StringVar(&configFile, "c", configFile, "config file")
	flag.StringVar(&configFile, "config", configFile, "config file")
	flag.DurationVar(&config.Intervall, "i", 5*time.Minute, "intervall between to polls")
	flag.DurationVar(&config.Intervall, "intervall", 5*time.Minute, "intervall between to polls")
	flag.StringVar(&url, "u", "", "url to poll")
	flag.StringVar(&url, "url", "", "url to poll")
	flag.Parse()
	config.Urls = strings.Split(url, ",")
	if notFromFile {
		return nil
	}
	return FromFile(config)
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
