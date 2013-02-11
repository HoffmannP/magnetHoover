package config

import (
	"encoding/json"
	"flag"
	// "github.com/HoffmannP/magnetHoover/history"
	"history"
	"os"
	// "github.com/HoffmannP/magnetHoover/parser"
	"log"
	"parser"
	"time"
	"transmission"
	// "github.com/HoffmannP/magnetHoover/transmission"
)

// var configFile = "/etc/magnetHoover.json"
var configFile = "config.json"

// var logFile = "/var/log/magnetHoover.log"
var logFile = "magnetHoover.log"

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
	Logger       *log.Logger
	loggingFile  *os.File
}

func FromCmdl() (*Config, error) {
	flag.StringVar(&configFile, "config", configFile, "config file")
	flag.StringVar(&logFile, "log", logFile, "log file")
	flag.Parse()
	return FromFile()
}

func FromFile() (c *Config, err error) {
	var cf ConfigFile
	c = new(Config)

	// Create logger and logging channel
	loggingFile, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	c.Logger = log.New(loggingFile, "", 0)
	c.loggingFile = loggingFile
	c.Logger.Print("Starting magnetHoover logging")
	var errorLogger chan string
	go func() {
		for true {
			c.Logger.Print(<-errorLogger)
		}
	}()

	// Read configuration 
	file, err := os.Open(configFile)
	if err != nil {
		return
	}
	if err := json.NewDecoder(file).Decode(&cf); err != nil {
		return c, err
	}
	file.Close()

	// Setup intervall
	if c.Intervall, err = time.ParseDuration(cf.Intervall); err != nil {
		c.Intervall = 5 * time.Minute
	}

	// History/DB COnnection
	if c.History, err = history.New(cf.Database, errorLogger); err != nil {
		c.Logger.Fatal("History: ", err)
		return c, err
	}

	// Connection to transmission client
	if c.Transmission, err = transmission.NewClient(cf.Transmission.SSL, cf.Transmission.Host, cf.Transmission.Port); err != nil {
		c.Logger.Fatal("Transmission: ", err)
		return c, err
	}

	// Preparation for polling URLs/Sites
	for _, uri := range cf.URIs {
		pf, err := parser.Parser(uri)
		if err != nil {
			c.Logger.Print(err)
		}
		c.URIs = append(c.URIs, pf)
	}
	return
}

func (c *Config) Close() {
	c.History.Close()
	c.Logger.Print("Closing magnetHoover logging")
	c.loggingFile.Close()
}
