package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/laohanlinux/go-logger/logger"
	. "github.com/laohanlinux/utils/memCachePool"
)

//	"time"

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	memCacheFactory := MemCacheFactory{}
	memCacheObj := memCacheFactory.GetMemCachePool(NoBlockingBytesChanType, 1024)
	m1, _ := memCacheObj.(*NoBlockingBytesChan)

	go func() {
		fmt.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	//
	memCacheObj = memCacheFactory.GetMemCachePool(NoBlockingChanType, 1024)
	m2, _ := memCacheObj.(*NoBlockingChan)

	// ask for a block
	var gGroup sync.WaitGroup

	currency, _ := strconv.Atoi(os.Args[1])
	for i := 0; i < currency; i++ {
		gGroup.Add(1)
		go func() {
			defer gGroup.Done()
			bench1(m1, m2, currency)
			time.Sleep(time.Millisecond * 1)
		}()
	}

	gGroup.Wait()
	fmt.Println("等待总控制进程退出")
	time.Sleep(time.Minute * 10)
}
func bench1(nbbc *NoBlockingBytesChan, nbc *NoBlockingChan, currency int) {
	send1 := nbbc.SendChan()
	recycle1 := nbbc.RecycleChan()
	c := <-send1

	send2 := nbc.SendChan()
	recycle2 := nbc.RecycleChan()

	var gGroup sync.WaitGroup
	exitChan := make(chan int)
	for i := 0; i < currency; i++ {
		gGroup.Add(1)
		go func() {
			defer gGroup.Done()
			for {
				select {
				case bytes := <-c:
					//display info
					logger.Info(string(bytes))
					//recycle
					recycle2 <- bytes
				case <-exitChan:
					logger.Info("exit=> bench1")
					return
				}
			}
		}()
	}

	// product binary data
	for i := 0; i < currency; i++ {
		b := <-send2
		b[0] = byte(i)
		c <- b
	}

	// for i := 0; i < currency; i++ {
	// 	exitChan <- i
	// }
	close(exitChan)
	gGroup.Wait()
	recycle1 <- c
	//
	fmt.Println("退出bench1 controller")
}

// func main() {
// 	//	_, r := memCache.NewNonBlockingChan()
// 	runtime.GOMAXPROCS(runtime.NumCPU())
// 	//	var g sync.WaitGroup
// 	currentNumbers, _ := strconv.Atoi(os.Args[1])
//
// 	//	for i := 0; i < currentNumbers; i++ {
// 	//		g.Add(1)
// 	//		go func(idx int) {
// 	//			defer g.Done()
// 	//			for i := 0; i < 10000; i++ {
// 	//				var s1 []byte
// 	//				fmt.Println("====")
// 	//				s1 = make([]byte, 1024)
// 	//				time.Sleep(time.Second * 1)
// 	//				r <- s1
// 	//			}
// 	//		}(i)
// 	//	}
//
// 	//	g.Wait()
// 	go func() {
// 		log.Println(http.ListenAndServe("localhost:6060", nil))
// 	}()
// 	ben(currentNumbers)
// }

// func ben(current int) {
// 	var g sync.WaitGroup
// 	s, r := memCache.NewNoBlockingChan(10240)
// 	for i := 0; i < current; i++ {
// 		g.Add(1)
// 		go func(idx int) {
// 			defer g.Done()
// 			for {
// 				s1 := <-s
// 				r <- s1
// 			}
// 		}(i)
// 	}
//
// 	g.Wait()
// }
