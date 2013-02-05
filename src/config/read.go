package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"history"
	"os"
	"plugin"
	"strings"
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
type URI struct {
	Parser plugin.ParserFunc
	URI    string
}
type Config struct {
	Intervall    time.Duration
	History      *history.History
	Transmission *transmission.Client
	URIs         []URI
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
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cf)
	if err != nil {
		return nil, err
	}

	if c.Intervall, err = time.ParseDuration(cf.Intervall); err != nil {
		fmt.Println(err)
		c.Intervall = 5 * time.Minute
	}
	if c.History, err = history.New(cf.Database); err != nil {
		return nil, err
	}
	if c.Transmission, err = transmission.NewClient(cf.Transmission.SSL, cf.Transmission.Host, cf.Transmission.Port); err != nil {
		return nil, err
	}
	for _, uri := range cf.URIs {
		parser := plugin.Default
		parts := strings.Split(uri, "§")
		if len(parts) > 1 {
			parser = plugin.Parser(parts[0])
			uri = parts[1]
		}
		c.URIs = append(c.URIs, URI{parser, uri})
	}
	return
}

func (c *Config) Close() {
	c.History.Close()
}
