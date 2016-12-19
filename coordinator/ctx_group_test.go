package coordinator

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/laohanlinux/assert"
)

func TestCoordinator(t *testing.T) {
	g := &sync.WaitGroup{}
	coor := NewCtxGCoordinator(g)
	runtime.GOMAXPROCS(runtime.NumCPU())
	total := int32(0)
	n := 1024
	for i := 0; i < n; i++ {
		childCoor := coor.StartIncr()
		go func(c *CtxGCoordinator, idx int) {
			//defer fmt.Println("goroutine", idx, "exit")
			defer atomic.AddInt32(&total, 1)
			for {
				if childCoor.TimeoutSub() {
					return
				}
				runtime.Gosched()
			}
		}(childCoor, i)
	}

	n++
	go func(c *CtxGCoordinator) {
		defer atomic.AddInt32(&total, 1)
		select {
		case <-c.Done():
			fmt.Println("Done=>")
		case <-time.After(time.Second * 3):
			fmt.Println("timeout")
		}
	}(coor)
	coor.Stop()
	assert.Equal(t, int32(n), total)

	g = &sync.WaitGroup{}
	coor = NewCtxGCoordinator(g)
	childCoor := coor.StartIncr()
	go func(c *CtxGCoordinator) {
		defer atomic.AddInt32(&total, 1)
		defer c.Sub()
		select {
		case <-c.Done():
			fmt.Println("Done=>")
		case <-time.After(time.Second * 2):
			fmt.Println("timeout")
		}
	}(childCoor)
	time.Sleep(time.Second * 3)
	coor.Stop()
}
