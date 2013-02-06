package main

import (
	"config"
	"log"
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
				log.Print("Rereading config file…")
				c_tmp, err := config.FromFile()
				if err == nil {
					c = (*cfg)(c_tmp)
					log.Println("Success")
					c.poll_all()
				} else {
					log.Println("Fail", err)
				}
			default:
				log.Printf("Received signal »%s«\n", si)
			}
		}
	}
}

func (c *cfg) Tic() <-chan time.Time {
	return time.After(c.Intervall)
}

func (c *cfg) poll_all() {
	// log.Println("Tic")
	max, end := 0, make(chan bool)
	for _, p := range c.URIs {
		max++
		go c.poll(p, end)
	}
	for max > 0 {
		select {
		case <-end:
			max--
		}
	}
	// log.Println("Tac")
	return
}

func (c *cfg) poll(p func() ([][]string, error), end chan bool) {
	es, err := p()
	if err != nil {
		log.Println(err)
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
			// log.Printf("%s already added\n", id)
			continue
		}

		if err = c.Transmission.Add(url); err != nil {
			log.Printf("Transmission Add Error »%s«: %v\n", url, err)
		}
		c.History.Add(id)
	}
	end <- true
	return
}
