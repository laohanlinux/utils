package main

import (
	"fmt"
	"github.com/laohanlinux/memCache"
	"os"
	"strconv"
	"sync"
	"time"
)

func main() {
	s, r := memCache.NewNonBlockingChan()

	var g sync.WaitGroup
	currentNumbers, _ := strconv.Atoi(os.Args[1])

	for i := 0; i < currentNumbers; i++ {
		g.Add(1)
		go func(idx int) {
			defer g.Done()
			for i := 0; i < 10000; i++ {
				var s1 []byte
				if i%2 == 0 {
					fmt.Println(idx, "申请到一块内存")
					s1 = <-s
				} else {
					fmt.Println("====")
					s1 = make([]byte, 1024)
				}
				time.Sleep(time.Second * 1)
				r <- s1
			}
		}(i)
	}

	g.Wait()
}
