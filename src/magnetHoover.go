package main

import (
	"config"
	"encoding/json"
)

var config Config.config

const SIGKILL = 9
const SIGHUP = 5

var signal chan byte

func main() {
	readFile := config.FromCmdl(&config)
	if  readFile {
		config.FromCmdl(&config)
	}
	run := true
	for run {
		select {
		case <- config.Intervall.After():
			poll()
		case s := <-signal:
			switch s {
			case SIGKILL:
				run = false
			case SIGHUP:
				readConfigFromFile()
			}
		}
	}
}

func poll() {

}