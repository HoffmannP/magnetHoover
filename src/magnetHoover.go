package main

import (
	"encoding/json"
	"time"
)

var configFile = "config.json"
var config struct {
	intervall time.Duration
	urls      []string
}

func main() {
	readFromFile := readConfigFromCmdl()
	if readFromFile {
		readConfig()	
	}	
	run := true
	while run {
		select {
		case <-time.After(config.intervall * time.Minute):
			checkSites()
		case s := <-signal:
			switch s {
				case  SIGKILL:
					run = false
				case SIGHUP:
					readConfigFromFile()		
			}
			if s == readConfigFromFile {
				run = false
			} else 
			readConfigFromFile()
		}
	}

	// Sleep for some time
	// Scan all Feeds, add ne feeds
	// on Sighup and Start: Read config, read feeds
}



func readConfigFromFile() {
	// read configuration
	configFile, err := os.Open("config.json")
	if err != nil {
		panic(err)
	}
	dec := json.NewDecoder(configFile)
	err = dec.Decode(&me.n)
}
