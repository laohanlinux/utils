package httputils

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"testing"
	"time"

	log "github.com/laohanlinux/utils/gokitlog"
)

func TestHTTPClientPool(t *testing.T) {
	p := NewHTTPClientPool()
	// use can set the dianFunction
	p.SetDialFn(DefaultDial)
	var g sync.WaitGroup
	n := 10

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
						continue
					}
				} else {
					req, err = http.NewRequest("POST", "http://www.baidu.com/post", nil)
					if err != nil {
						log.Error("err", err)
						continue
					}
					req.Header.Add("Content-Type", "application/json")

					resp, err = c.C.Do(req)
					if err != nil {
						log.Error("err", err)
						continue
					}
				}

				b, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					log.Crit(err)
				}
				fmt.Println(string(b))
				time.Sleep(time.Second)
				p.Recycle(c)
			}
		}()
	}

	g.Wait()
}
