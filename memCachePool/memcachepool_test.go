package memCachePool

import (
	"fmt"
	"runtime"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/laohanlinux/assert"
	"github.com/laohanlinux/go-logger/logger"
)

func TestNoBlockingChan(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	memCacheFactory := MemCacheFactory{}

	memCacheObj := memCacheFactory.GetMemCachePool(NoBlockingChanType, 1024)
	//
	m, ok := memCacheObj.(*NoBlockingChan)
	assert.Equal(t, ok, true)
	send := m.SendChan()
	recycle := m.RecycleChan()

	// asks for a block
	b := <-send
	b[0] = 'A'
	recycle <- b

	// asks for a block
	<-send
	b1 := <-send
	chacter := fmt.Sprintf("%d", b[0])
	assert.Equal(t, chacter, "65")
	assert.Equal(t, b, b1)
}

func TestNoBlockingBytesChan(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	memCacheFactory := MemCacheFactory{}
	memCacheObj := memCacheFactory.GetMemCachePool(NoBlockingBytesChanType, 1024)
	m, ok := memCacheObj.(*NoBlockingBytesChan)
	assert.Equal(t, ok, true)

	send := m.SendChan()
	recycle := m.RecycleChan()

	// ask for a block
	var gGroup sync.WaitGroup
	c := <-send

	exitChan := make(chan int)
	for i := 0; i < 1024; i++ {
		gGroup.Add(1)
		go func() {
			defer gGroup.Done()
			for {
				select {
				case bytes := <-c:
					logger.Info(string(bytes))
				case <-exitChan:
					logger.Info("exit")
					return
				}
			}
		}()
	}
	// start to product data
	for i := 0; i < 1024; i++ {
		c <- []byte(strconv.Itoa(i))
	}
	//close enter point
	close(exitChan)

	gGroup.Wait()
	c <- []byte("a")
	//assert.Equal(t, len(c), 1)
	recycle <- c
}

func BenchmarkNoBlockingBytesChan1(b *testing.B) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	memCacheFactory := MemCacheFactory{}
	memCacheObj := memCacheFactory.GetMemCachePool(NoBlockingBytesChanType, 1024)
	m1, _ := memCacheObj.(*NoBlockingBytesChan)

	//
	memCacheObj = memCacheFactory.GetMemCachePool(NoBlockingChanType, 1024)
	m2, _ := memCacheObj.(*NoBlockingChan)

	// ask for a block
	var gGroup sync.WaitGroup

	for i := 0; i < b.N; i++ {
		gGroup.Add(1)
		go func() {
			defer gGroup.Done()
			Bench1(m1, m2, b.N)
			time.Sleep(time.Millisecond * 1)
		}()
	}

	gGroup.Wait()
}

func Bench1(nbbc *NoBlockingBytesChan, nbc *NoBlockingChan, currency int) {
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
			select {
			case bytes := <-c:
				//display info
				logger.Info(string(bytes))
				//recycle
				recycle2 <- bytes
			case <-exitChan:
				return
			}
		}()
	}

	// product binary data
	for i := 0; i < currency; i++ {
		b := <-send2
		b[0] = byte(i)
		c <- b
	}

	close(exitChan)
	gGroup.Wait()
	recycle1 <- c
}
