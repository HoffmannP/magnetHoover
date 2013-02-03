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

	run := false // true
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
	var max int
	end := make(chan bool)
	for _, url := range urls {
		max++
		go poll(url, end)
	}
	for max > 0 {
		select {
		case <-end:
			max--
		}
	}
	return
}

func poll(url string, end chan bool) {
	response, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer response.Body.Close()
	body := make([]byte, 1<<10)
	n, err := response.Body.Read(body)
	for err == nil {
		tmp := make([]byte, 1<<10)
		n, err = response.Body.Read(tmp)
		body = append(body, tmp[:n]...)
	}
	fmt.Printf("%s\n", body)
	end<- true
}
