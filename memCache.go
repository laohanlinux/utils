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

	THRESHOLD_FREE_OS_MEMORY = 268435456
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
	SetBufferSize(uint64)
}

var nbc *NonBlockingChan

// define a []byte type memcache pool
type NonBlockingChan struct {
	send      chan []byte
	recv      chan []byte
	freeMem   chan byte
	blockSize uint64
}

// memCache Object
type noBuffferObj struct {
	b    []byte
	used time.Time
}

var memCacheOnce sync.Once

// blockSize: memcache block size, default 4k
func NewNonBlockingChan(blockSize ...int) (<-chan []byte, chan<- []byte) {
	memCacheOnce.Do(func() {
		nbc = &NonBlockingChan{
			send:      make(chan []byte),
			recv:      make(chan []byte),
			freeMem:   make(chan byte),
			blockSize: 1024 * 4,
		}
		go nbc.doWork()
		go nbc.freeOldMemCache()
	})
	return nbc.send, nbc.recv
}

// Very Block is 4kb
//func (nbc NonBlockingChan) makeBuffer() []byte { return make([]byte, nbc.blockSize) }
func (nbc *NonBlockingChan) makeBuffer() []byte { return make([]byte, nbc.blockSize) }

func (nbc *NonBlockingChan) bufferSize() uint64 { return 0 }

func (nbc *NonBlockingChan) SetBufferSize(blockSize uint64) { nbc.blockSize = blockSize }

func (nbc *NonBlockingChan) doWork() {
	defer func() {
		debug.FreeOSMemory()
	}()

	items := list.New()
	//timeout := time.NewTimer(time.Minute)
	timeout := time.NewTimer(time.Second * 5)
	for {
		if items.Len() == 0 {
			items.PushFront(noBuffferObj{
				b:    nbc.makeBuffer(),
				used: time.Now(),
			})
		}
		select {
		case e := <-nbc.recv:
			items.PushBack(noBuffferObj{
				b:    e,
				used: time.Now(),
			})
		case nbc.send <- items.Remove(items.Front()).(noBuffferObj).b:
		case <-timeout.C:
			// do clear memcache work, because no one are use the memcache pool
			itemsLen := items.Len()
			fmt.Println("do clear memCache work:", time.Now().Unix())
			if itemsLen > 0 {
				for i := 0; i < itemsLen/2; i++ {
					items.Remove(items.Front())
				}
				debug.FreeOSMemory()
			}
			timeout.Reset(time.Second * 5)
		case <-nbc.freeMem:
			// free too old memcached
			fmt.Println("开始释放旧数据=>", items.Len())
			item := items.Front()
			if item != nil {
				fmt.Printf("内存对象最新使用时间:%d\n", item.Value.(noBuffferObj).used.Unix())
			}
			var freeSize uint64
			for item != nil {
				nItem := item.Next()
				fmt.Printf("内存对象最新使用时间:%d\n", item.Value.(noBuffferObj).used.Unix())
				if time.Since(item.Value.(noBuffferObj).used) > (time.Second * time.Duration(10)) {
					fmt.Println("free memCache obj:", item)
					items.Remove(item)
				} else {
					break
				}
				item = nItem
				freeSize += nbc.blockSize
			}
			if freeSize > THRESHOLD_FREE_OS_MEMORY {
				debug.FreeOSMemory()
			}
		}
	}
}

// free old memcache object, timeout = 1 minute not to be used
func (nbc *NonBlockingChan) freeOldMemCache() {
	//timeout := time.NewTimer(time.Minute * 5)
	timeout := time.NewTicker(time.Second * 10)
	for {
		select {
		case <-timeout.C:
			fmt.Println("free memCacheOnce timeout:", time.Now().Unix())
			nbc.freeMem <- 'f'
		}
	}
}
