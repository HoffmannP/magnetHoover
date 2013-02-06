package transmission

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

const rpc_location = "transmission/rpc"
const session_header_key = "X-Transmission-Session-Id"

type request struct {
	u string
	r io.Reader
	p int
}

type response struct {
	Result    string
	Arguments interface{}
	Tag       int64
}

type Client struct {
	a string
	c *http.Client
	s string
	n int64
}

func NewClient(ssl bool, host string, port int) (*Client, error) {
	var protocoll string
	if ssl {
		protocoll = "https"
	} else {
		protocoll = "http"
	}

	c := &Client{
		fmt.Sprintf("%s://%s:%d/%s", protocoll, host, port, rpc_location),
		&http.Client{Transport: &http.Transport{DisableKeepAlives: true}},
		"",
		time.Now().Unix() * 1000}

	r, err := http.DefaultClient.Head(c.a)
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
	r := new(request)
	r.u = uri
	return r
}
func (r *request) tag(t int64) {
	b := " {'method': 'torrent-add', 'arguments': {'filename': '%s'}, 'tag': %d}"
	r.r = strings.NewReader(strings.Replace(fmt.Sprintf(b, r.u, t), "'", "\"", -1))
}
func (q *request) Read(r []byte) (n int, err error) {
	return q.r.Read(r)
}

func (c *Client) Call(r *request) error {
	tag := c.n
	c.n++
	r.tag(tag)

	req, err := http.NewRequest("POST", c.a, r)
	if err != nil {
		return err
	}
	req.Header.Set(session_header_key, c.s)

	res, err := c.c.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var s response
	err = json.NewDecoder(res.Body).Decode(&s)
	res.Body.Close()
	if err != nil {
		return err
	}
	if s.Tag != c.n || s.Result != "success" {
		return errors.New(fmt.Sprintf("RPC Torrent Add Error »%s«", s.Result))
	}
	return nil
}

func (c Client) Add(uri string) error {
	err := c.Call(request_add(uri))
	if err != nil {
		log.Println(err)
	}
	return nil
}
