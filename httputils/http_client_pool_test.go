package httputils

import (
	"io/ioutil"
	"net/http"
	"sync"
	"testing"
	"time"

	"fmt"
	log "github.com/laohanlinux/utils/gokitlog"
	"runtime"
)

func TestHTTPClient(t *testing.T) {
	c := NewHTTPClient(DefaultDial)

	for {
		resp, err := c.C.Get("http://www.baidu.com")
		if err != nil {
			log.Crit("err", err)
		}
		for k, headers := range resp.Header {
			fmt.Println(k, headers)
		}
		_, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Crit("err", err)
		}
	}
}

func TestHTTPClientPool(t *testing.T) {
	p := NewHTTPClientPool(time.Second * 3)
	runtime.GOMAXPROCS(runtime.NumCPU())
	// use can set the dianFunction
	p.SetDialFn(DefaultDial)
	sleepTime := time.Millisecond * 1000
	var g sync.WaitGroup
	n := 4

	for ; n > 0; n-- {
		g.Add(1)
		go func() {
			defer g.Done()
			for {
				var (
					resp *http.Response
					err  error
					c    = p.GetHttpClient()
					req  *http.Request
				)
				if time.Now().UnixNano()%2 == 0 {
					resp, err = c.C.Get("http://www.baidu.com")
					if err != nil {
						log.Error("err", err)
						time.Sleep(sleepTime)
						continue
					}
				} else {
					req, err = http.NewRequest("POST", "http://www.baidu.com/post", nil)
					if err != nil {
						log.Error("err", err)
						time.Sleep(sleepTime)
						continue
					}
					req.Header.Add("Content-Type", "application/json")

					resp, err = c.C.Do(req)
					if err != nil {
						log.Error("err", err)
						time.Sleep(sleepTime)
						continue
					}
				}

				_, err = ioutil.ReadAll(resp.Body)
				if err != nil {
					log.Crit(err)
				}
				//fmt.Println(string(b))
				time.Sleep(sleepTime)
				log.Debug("Recycle")
				p.Recycle(c, resp)
			}
		}()
	}

	g.Wait()
}
