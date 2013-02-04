package main

import (
	"config"
	"fmt"
	"io"
	"net/http"
	"os"
	"net/rpc"
	"os/signal"
	"syscall"
	"time"
)

type sigRec chan os.Signal
type cfg config.Config
type cl rpc.Client

type rpc_command struct {
	method string
	arguments struct {
		filename string
	}
}
func (r *cl) add(file string) error {
	rpc := new(rpc_command)
	rpc.method = "add_torrent"
	rpc.arguments.filename = file
	return (*rpc(r)).Call(rpc)
}

func main() {
	c, err := config.FromCmdl()
	if err != nil {
		panic(err)
	}
	defer (*cfg(c)).Close()

	s := make(chan os.Signal)
	signal.Notify(s)
	sigRec(s).circleOfLife((*cfg)(c))
}

func (s sigRec) circleOfLife(c *cfg) {
	c.poll_all()
	run := false // true
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
	es, err := u.Parser(readAll(response.Body))
	if err != nil {
		fmt.Println(err)
		end <- false
		return
	}
	if len(es) == 0 {
		end <- false
		return
	}
	c.History.Select(u.URI)
	for _, e := range es {
		id, url := e[0], e[1]
		if c.History.Exists(id) {
			fmt.Printf("%s already added\n", id)
			continue
		}
		if err = c.Transmission.add(url); err == nil {
			if err = c.History.Add(id); err != nil {
				fmt.Printf("SQL Add Error »%s«: %v\n", id, err)
			}
		}
	}
	end <- true
	return
}

func (c *cfg) Close() {
	c.History.Close()
	c.Transmission.Close()
}

func readAll(r io.Reader) string {
	const slicesize = 1 << 10
	text := make([]byte, slicesize)
	n, err := r.Read(text)
	for err == nil {
		tmp := make([]byte, slicesize)
		n, err = r.Read(tmp)
		text = append(text, tmp[:n]...)
	}
	return string(text)
}
