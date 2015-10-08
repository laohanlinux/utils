package memCache

import (
	"container/list"
	//	"errors"
	"fmt"
	"runtime/debug"
	"sync"
	"time"
)

const (
	BinaryMemCachePool = "binarymemcachepool"

	UnKownMemCachePoolType = "unkown memCachePool type"
)

/**
type MemCachePoolFactory struct {
}

func (mcpf MemCachePoolFactory) GetMemCachePool(memCachePoolType string) (*MemCachePool, error) {
	switch {
	case memCachePoolType == BinaryMemCachePool:
		if nbc == nil {
			send, recv := newNonBlockingChan()
			nbc = &NonBlockingChan{
				send: send,
				recv: recv,
			}
			return (*MemCachePool)(nbc), nil
		}
		return nil, nil
	default:
		return nil, errors.New(UnKownMemCachePoolType)
	}
}
*/
type MemCachePool interface {
	bufferSize() uint64
}

var nbc *NonBlockingChan

// define a []byte type memcache pool
type NonBlockingChan struct {
	send chan []byte
	recv chan []byte
	once sync.Once
}

func NewNonBlockingChan() (<-chan []byte, chan<- []byte) {
	if nbc == nil {
		nbc = &NonBlockingChan{
			send: make(chan []byte),
			recv: make(chan []byte),
		}
	}

	if nbc != nil {
		nbc.once.Do(func() { go nbc.doWork() })
	}

	return nbc.send, nbc.recv
}

func (nbc NonBlockingChan) makeBuffer() []byte { return make([]byte, 1024*1024*32) }

func (nbc NonBlockingChan) bufferSize() uint64 { return 0 }

func (nbc *NonBlockingChan) doWork() {
	defer func() {
		if err := recover(); err != nil {
			// do some something
			fmt.Println("MemCachePool unexpectedly terminated, err:", err)
			go nbc.doWork()
		}
		debug.FreeOSMemory()
	}()

	items := list.New()
	timeout := time.NewTimer(time.Minute)
	for {
		if items.Len() == 0 {
			items.PushFront(nbc.makeBuffer())
		}
		select {
		case e := <-nbc.recv:
			items.PushBack(e)
		case nbc.send <- items.Remove(items.Front()).([]byte):
		case <-timeout.C:
			// do clear memcache work, because no one are use the memcache pool
			itemsLen := items.Len()
			if itemsLen > 0 {
				for i := 0; i < itemsLen/2; i++ {
					items.Remove(items.Front())
				}
				debug.FreeOSMemory()
			}
			timeout.Reset(time.Minute)
		}
	}
}
