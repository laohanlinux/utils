package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"time"

	"math/rand"
	"net"
	"runtime"
	"strconv"
	"sync"

	"github.com/laohanlinux/memCache/kgredis"
)

func doSetWork(rdpc *kgredis.RedisPoolChan) {
	var g sync.WaitGroup
	for i := 0; i < numbers; i++ {
		g.Add(1)
		go func(key int) {
			defer g.Done()
			for {
				time.Sleep(time.Second * 1)
				r, err := rdpc.Get(1000)

				if err != nil {
					fmt.Println(err)
					continue
				}

				err = r.ReadReply().Err
				if err != nil {
					fmt.Printf("Get Reply:%s\n", err)
					rdpc.Put(r, err)
				}
				err = r.Cmd("SET", strconv.Itoa(key), strconv.Itoa(key)).Err
				if err != nil {
					fmt.Println(err)
				}
				rdpc.Put(r, err)
			}
		}(i)
	}

	g.Wait()

}

func doGetWork(rdpc *kgredis.RedisPoolChan) {
	var g sync.WaitGroup
	for i := 0; i < numbers; i++ {
		g.Add(1)
		go func() {
			defer g.Done()
			for {
				time.Sleep(time.Second * 1)
				idx := rand.Intn(1024)
				r, err := rdpc.Get(1000)
				if err != nil {
					fmt.Println(err)
					continue
				}
				//err = r.ReadReply()

				reply := r.Cmd("GET", strconv.Itoa(idx))

				if reply.Err != nil {
					if _, ok := reply.Err.(*net.OpError); ok {
						fmt.Printf("reply: %v\n", reply.Err)
						rdpc.Put(r, err)
						continue
					}
				}

				value, err := reply.Bytes()
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Printf("value: %s\n", value)
				}
				rdpc.Put(r, nil)
			}
		}()
	}
	g.Wait()

}

var numbers = 1

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())
	//NewRedisPoolConn(address string, maxIdle, maxActive int, idleTimeout, defaultExpire int64)
	rdpc, err := kgredis.NewRedisPoolChan("localhost:6379", 100)
	if err != nil {
		fmt.Println(err)
		return
	}
	go doSetWork(rdpc)
	go doGetWork(rdpc)

	http.ListenAndServe(":6060", nil)
	time.Sleep(time.Second * 30)
}
