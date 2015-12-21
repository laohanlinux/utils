package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"strconv"
	"sync"

	"github.com/laohanlinux/memCache"
	//	"time"
)

func main() {
	//	_, r := memCache.NewNonBlockingChan()
	runtime.GOMAXPROCS(runtime.NumCPU())
	//	var g sync.WaitGroup
	currentNumbers, _ := strconv.Atoi(os.Args[1])

	//	for i := 0; i < currentNumbers; i++ {
	//		g.Add(1)
	//		go func(idx int) {
	//			defer g.Done()
	//			for i := 0; i < 10000; i++ {
	//				var s1 []byte
	//				fmt.Println("====")
	//				s1 = make([]byte, 1024)
	//				time.Sleep(time.Second * 1)
	//				r <- s1
	//			}
	//		}(i)
	//	}

	//	g.Wait()
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	ben(currentNumbers)
}

func ben(current int) {
	var g sync.WaitGroup
	s, r := memCache.NewNoBlockingChan(10240)
	for i := 0; i < current; i++ {
		g.Add(1)
		go func(idx int) {
			defer g.Done()
			for {
				s1 := <-s
				r <- s1
			}
		}(i)
	}

	g.Wait()
}
