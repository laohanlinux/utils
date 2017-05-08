package httputils

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	DefaultAddr    = "localhost:80"
	DefaultNetWork = "tcp"
)

var (
	defaultConnKeepTime = time.Duration(time.Second * 30)
)

func NewHTTPClientPool(connKeepTime time.Duration) *HttpClientPool {
	defaultConnKeepTime = connKeepTime
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

// Recycle should recycle the client that is not error
func (hcp *HttpClientPool) Recycle(c *HTTPClient, resp ...*http.Response) {
	if len(resp) > 0 && resp[0] != nil {
		if value := resp[0].Header.Get("Connection"); value != "" {
			value = strings.ToLower(value)
			if !strings.Contains(value, "keep-alive") ||
				strings.Contains(value, "close") {
				resp[0].Close = true
				c = nil
				return
			}
		}
	}

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
		Timeout:   defaultConnKeepTime,
		KeepAlive: defaultConnKeepTime,
	}
	conn, err := dial.Dial(network, addr)
	if err != nil {
		return nil, err
	}
	return conn, err
}
