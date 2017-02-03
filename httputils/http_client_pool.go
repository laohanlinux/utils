package httputils

import (
	"net"
	"net/http"
	"sync"
	"time"
)

const (
	DefaultAddr    = "localhost:80"
	DefaultNetWork = "tcp"
)

type HttpClientPool struct {
	httpClientLock *sync.Mutex
	freeClient     *HTTPClient
}

func (hcp *HttpClientPool) GetHttpClient() *HTTPClient {
	hcp.httpClientLock.Lock()
	c := hcp.freeClient
	if c == nil {
		c = new(HTTPClient)
	} else {
		hcp.freeClient = c.next
	}
	hcp.httpClientLock.Unlock()
	return c
}

func (hcp *HttpClientPool) recycle(c *HTTPClient) {
	hcp.httpClientLock.Lock()
	c.next = hcp.freeClient
	hcp.freeClient = c
	hcp.httpClientLock.Unlock()
}

type HTTPClient struct {
	C    *http.Client
	next *HTTPClient
}

func NewHTTPClient() *HTTPClient {
	// conn, err := DefaultDial(DefaultNetWork, DefaultAddr)
	// if err != nil {
	// 	log.Crit("err", err)
	// }
	return &HTTPClient{
		C: &http.Client{
			Transport: &http.Transport{
				Dial:              DefaultDial,
				DisableKeepAlives: false,
			},
		},
	}
}

func DefaultDial(network, addr string) (net.Conn, error) {
	dial := net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	conn, err := dial.Dial(network, addr)
	if err != nil {
		return nil, err
	}
	return conn, err
}
