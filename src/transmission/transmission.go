package transmission

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

const rpc_location = "transmission/rpc"
const session_header_key = "X-Transmission-Session-Id"

type request struct {
	u string
	r []byte
	p int
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
	r := new(request)
	r.u = uri
	return r
}
func (r *request) tag(t int64) {
	b := "{'method': 'torrent-add', 'arguments': {'filename': '%s'}, 'tag': %d}"
	r.r = ([]byte)(strings.Replace(fmt.Sprintf(b, r.u, t), "'", "\"", -1))
	log.Printf("%s", r.r)
}
func (q *request) Read(r []byte) (n int, err error) {
	n = len(r)
	if q.p+n > len(q.r) {
		r = r[:len(q.r)-q.p]
		n = len(r)
		err = errors.New("EOF")
	}
	r = q.r[q.p : q.p+n]
	fmt.Println(".", r, ".", q.r, ".", q.p, q.p+n)
	q.p = q.p + n
	return
}

func (c *Client) Call(r *request) error {
	tag := c.n
	c.n++
	r.tag(tag)

	/*
		req, err := http.NewRequest("POST", c.a, r)
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Set(session_header_key, c.s)
		req.Header.Set("Content-Type", "application/json")
		req.Close = true
		log.Printf("Sending %v", req)
	*/

	buffer := make([]byte, 200)
	n, err := r.Read(buffer)
	if err != nil {
		log.Print(err)
	}
	fmt.Printf("<%s>:%T#%d\n", buffer, buffer, n)

	/*
		res, err := c.c.Do(req)
		defer res.Body.Close()
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Sent")

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
