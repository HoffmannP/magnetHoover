package main

import (
	// "github.com/HoffmannP/magnetHoover/config"
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
	log.Println("Starting magnetHoover")
	c, err := config.FromCmdl()
	if err != nil {
		if c == nil {
			log.Fatal(err)
		} else {
			c.Logger.Print("..")
		}
	}
	defer c.Close()

	s := make(chan os.Signal)
	signal.Notify(s, signal.)
	sigRec(s).circleOfLife((*cfg)(c))
	log.Println("Closing magnetHoover")
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
			case syscall.SIGTERM:
				run = false
			case syscall.SIGINT:
				run = false
			case syscall.SIGHUP:
				c.Logger.Print("Rereading config file…")
				c_tmp, err := config.FromFile()
				if err == nil {
					c = (*cfg)(c_tmp)
					c.Logger.Print("Successfully reload config file")
					c.poll_all()
				} else {
					c.Logger.Print("Error reloading config file (using old values)", err)
				}
			default:
				c.Logger.Printf("Received signal »%s«", si)
			}
		}
	}
}

func (c *cfg) Tic() <-chan time.Time {
	return time.After(c.Intervall)
}

func (c *cfg) poll_all() {
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
	return
}

func (c *cfg) poll(p func() ([][]string, error), end chan bool) {
	es, err := p()
	if err != nil {
		c.Logger.Print(err)
		end <- false
		return
	}
	if len(es) == 0 {
		end <- false
		return
	}
	for _, e := range es {
		id, url := e[0], e[1]
		switch exists, err := c.History.Exists(id); {
		case err != nil:
			c.Logger.Printf(err.Error())
			continue
		case exists:
			c.Logger.Printf("%s already added\n", id)
			continue
		}
		if err = c.Transmission.Add(url); err != nil {
			c.Logger.Printf("Transmission Add Error »%s«: %v\n", url, err)
		} else {
			c.Logger.Printf("Added »%s«", url)
		}

		c.History.Add(id)
	}
	end <- true
	return
}
