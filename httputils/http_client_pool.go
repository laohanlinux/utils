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

func NewHTTPClientPool() *HttpClientPool {
	return &HttpClientPool{
		httpClientLock: &sync.Mutex{},
	}
}

type HttpClientPool struct {
	httpClientLock *sync.Mutex
	dialFn         DialFunc
	freeClient     *HTTPClient
}

func (hcp *HttpClientPool) GetHttpClient() *HTTPClient {
	hcp.httpClientLock.Lock()
	c := hcp.freeClient
	if c == nil {
		c = NewHTTPClient(hcp.dialFn)
	} else {
		hcp.freeClient = c.next
	}
	hcp.httpClientLock.Unlock()
	return c
}

func (hcp *HttpClientPool) Recycle(c *HTTPClient) {
	hcp.httpClientLock.Lock()
	c.next = hcp.freeClient
	hcp.freeClient = c
	hcp.httpClientLock.Unlock()
}

func (hcp *HttpClientPool) SetDialFn(dialFn DialFunc) {
	hcp.dialFn = dialFn
}

type DialFunc func(network, addr string) (net.Conn, error)

type HTTPClient struct {
	C    *http.Client
	next *HTTPClient
}

func NewHTTPClient(dialFn DialFunc) *HTTPClient {
	if dialFn == nil {
		dialFn = DefaultDial
	}
	return &HTTPClient{
		C: &http.Client{
			Transport: &http.Transport{
				Dial:              dialFn,
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
