package main

import (
	"config"
	"os/signal"
	"os"
	"syscall"
	"fmt"
	"net/http"
)

func main() {
	var cfg config.Config
	config.FromCmdl(&cfg)
	
	sig := make(chan os.Signal)
	signal.Notify(sig)

	run := true
	poll_all(cfg.Urls)
	for run {
		select {
		case <-cfg.Tic():
			poll_all(cfg.Urls)
		case si := <-sig:
			switch si {
			case syscall.SIGINT:
				run = false
			case syscall.SIGHUP:
				fmt.Println("Rereading config file")
				config.FromFile(&cfg)
			}
		}
	}
}

func poll_all(urls []string) {
	if len(urls) == 0 {
		return
	}
	for _, url := range urls {
		go poll(url)
	}
	return
}

func poll(url string) {
	response, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer response.Body.Close()
	fmt.Println(response, response.ContentLength)
}