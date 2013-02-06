package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"history"
	"log"
	"os"
	"parser"
	"time"
	"transmission"
)

var configFile = "config.json"

type ConfigFile struct {
	Intervall    string
	Database     string
	Transmission struct {
		Host string
		Port int
		SSL  bool
	}
	URIs []string
}
type Config struct {
	Intervall    time.Duration
	History      *history.History
	Transmission *transmission.Client
	URIs         []parser.ParserFunc
}

func FromCmdl() (*Config, error) {
	flag.StringVar(&configFile, "config", configFile, "config file")
	flag.Parse()
	return FromFile()
}

func FromFile() (c *Config, err error) {
	var cf ConfigFile
	c = new(Config)
	file, err := os.Open(configFile)
	if err != nil {
		return nil, err
	}
	if err := json.NewDecoder(file).Decode(&cf); err != nil {
		return nil, err
	}
	if c.Intervall, err = time.ParseDuration(cf.Intervall); err != nil {
		fmt.Println(err)
		c.Intervall = 5 * time.Minute
	}
	if c.History, err = history.New(cf.Database); err != nil {
		log.Fatal("History: ", err)
	}
	if c.Transmission, err = transmission.NewClient(cf.Transmission.SSL, cf.Transmission.Host, cf.Transmission.Port); err != nil {
		log.Fatal("Transmission: ", err)
	}
	for _, uri := range cf.URIs {
		c.URIs = append(c.URIs, parser.Parser(uri))
	}
	return
}

func (c *Config) Close() {
	c.History.Close()
}
