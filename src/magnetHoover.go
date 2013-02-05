package main

import (
	"config"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type sigRec chan os.Signal
type cfg config.Config

func main() {
	c, err := config.FromCmdl()
	if err != nil {
		panic(err)
	}
	defer c.Close()

	s := make(chan os.Signal)
	signal.Notify(s)
	sigRec(s).circleOfLife((*cfg)(c))
}

func (s sigRec) circleOfLife(c *cfg) {
	c.poll_all()
	run := true
	for run {
		select {
		case <-c.Tic():
			c.poll_all()
		case si := <-s:
			switch si {
			case syscall.SIGINT:
				run = false
			case syscall.SIGHUP:
				fmt.Print("Rereading config file ")
				c_tmp, err := config.FromFile()
				if err == nil {
					c = (*cfg)(c_tmp)
					fmt.Println("Sucess")
				} else {
					fmt.Println("FAIL")
				}
			}
		}
	}
}

func (c *cfg) Tic() <-chan time.Time {
	return time.After(c.Intervall)
}

func (c *cfg) poll_all() {
	println("Tic")
	max, end := 0, make(chan bool)
	for _, u := range c.URIs {
		max++
		go c.poll(u, end)
	}
	for max > 0 {
		select {
		case <-end:
			max--
		}
	}
	return
}

func (c *cfg) poll(u config.URI, end chan bool) {
	response, err := http.Get(u.URI)
	if err != nil {
		fmt.Println(err)
	}
	defer response.Body.Close()
	bs, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		end <- false
		return
	}
	es, err := u.Parser(string(bs))
	if err != nil {
		fmt.Println(err)
		end <- false
		return
	}
	if len(es) == 0 {
		end <- false
		return
	}
	for _, e := range es {
		id, url := e[0], e[1]
		if c.History.Exists(id) {
			fmt.Printf("%s already added\n", id)
			continue
		}

		if err = c.Transmission.Add(url); err != nil {
			fmt.Printf("Transmission Add Error »%s«: %v\n", url, err)
			panic(err)
			c.History.Add(id)
		}
	}
	end <- true
	return
}
