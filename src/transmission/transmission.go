package transmission

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const rpc_location = "transmission/rpc"
const session_header_key = "X-Transmission-Session-Id"

type request struct {
	u string
	t int64
}

type response struct {
	Result    string
	Arguments interface{}
	Tag       int64
}

type Client struct {
	a string
	s string
	c *http.Client
	n int64
}

func NewClient(ssl bool, host string, port int) (*Client, error) {
	var protocoll string
	if ssl {
		protocoll = "https"
	} else {
		protocoll = "http"
	}

	c := &Client{fmt.Sprintf("%s://%s:%d/%s", protocoll, host, port, rpc_location), "", new(http.Client), time.Now().Unix() * 1000}

	r, err := c.c.Head(c.a)
	if err != nil {
		return nil, err
	}
	if r.StatusCode != 409 {
		return nil, errors.New(fmt.Sprintf("Unexpected status return code »%d« (wrong address?)", r.StatusCode))
	}
	c.s = r.Header.Get(session_header_key)
	if c.s == "" {
		return nil, errors.New(fmt.Sprintf("Missing session header key »%s« (wrong address?)", session_header_key))
	}
	return c, nil
}

func request_add(uri string) *request {
	return &request{uri, 0}
}
func (r *request) tag(t int64) {
	r.t = t
}
func (r *request) Read(d []byte) (int, error) {
	t := "{'method': 'torrent-add', 'arguments': {'filename': '%s'}, 'tag': %d}"
	b := bytes.NewBufferString(strings.Replace(fmt.Sprintf(t, r.u, r.t), "'", "\"", -1))
	n, err := b.Read(d)
	return n, err
}

func (c *Client) Call(r *request) error {
	tag := c.n
	c.n++
	r.tag(tag)

	req, err := http.NewRequest("POST", c.a, io.Reader(r))
	if err != nil {
		return err
	}
	req.Header.Set(session_header_key, c.s)
	req.Header.Set("Content-Type", "application/json")
	fmt.Printf("PreSend: %v\n", req)

	res, err := c.c.Do(req)
	if err != nil {
		return err
	}

	if err != nil {
		log.Fatal(err)
	}
	robots, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", robots)

	/*
		var s response
		err = json.NewDecoder(res.Body).Decode(&s)
		res.Body.Close()
		if err != nil {
			return err
		}
		if s.Tag != c.n || s.Result != "success" {
			return errors.New(fmt.Sprintf("RPC Torrent Add Error »%s«", s.Result))
		}
	*/
	return nil
}

func (c Client) Add(uri string) error {
	return c.Call(request_add(uri))
}
