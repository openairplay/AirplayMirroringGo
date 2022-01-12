/**
Created by: Joseph Han <joseph.bing.han@gmail.com>.
Created at: 2020-12-31 09:04
*/

package airplay

import (
	"bytes"
	"fmt"
	"go.uber.org/atomic"
	"log"
	"net"
	"net/http"
	"time"
)

type AppleTVClient struct {
	url     string
	Stream  chan []byte
	client  http.Client
	start   atomic.Bool
	pending atomic.Bool
}

func NewAppleTVClient(ip net.IP) *AppleTVClient {
	c := AppleTVClient{
		url:    fmt.Sprintf("http://%s:7000/photo", ip),
		Stream: make(chan []byte),
	}

	c.start.Store(false)
	c.pending.Store(false)
	go c.serve()
	return &c
}
func (c *AppleTVClient) serve() {

	for {
		if c.start.Load() {
			select {
			case buf := <-c.Stream:
				if c.pending.Load() {
					time.Sleep(250 * time.Millisecond)
				} else {
					go func() {
						c.pending.Store(true)
						c.sendPhoto(buf)
						c.pending.Store(false)
					}()
				}
			default:
				time.Sleep(10 * time.Millisecond)
			}
		} else {
			time.Sleep(100 * time.Millisecond)
		}

	}
}

func (c *AppleTVClient) SetAppleTVIP(ip net.IP) {
	c.url = fmt.Sprintf("http://%s:7000/photo", ip)
}
func (c *AppleTVClient) sendPhoto(buf []byte) {
	var (
		req *http.Request
		err error
	)
	req, err = http.NewRequest("PUT", c.url, bytes.NewReader(buf))
	if err != nil {
		log.Fatalf("Create http request fail: %v", err)
	}
	req.ContentLength = int64(len(buf))

	_, err = c.client.Do(req)
	if err != nil {
		log.Fatalf("Create http request fail: %v", err)
	}
}
func (c *AppleTVClient) Start() {
	c.start.Store(true)

}

func (c *AppleTVClient) Stop() {
	c.start.Store(false)
}
