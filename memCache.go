package memCache

import (
	"container/list"
	"errors"
	"fmt"
	"runtime/debug"
	"sync"
	"time"
)

var BinaryMemCachePool = errors.New("binarymemcachepool")

var UnKownMemCachePoolType = errors.New("unkown memCachePool type")

const (
	THRESHOLD_FREE_OS_MEMORY = 268435456
)

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
	for {
		if items.Len() == 0 {
			items.PushFront(noBuffferObj{
				b:    nbc.makeBuffer(),
				used: time.Now(),
			})
		}
		e := items.Front()
		select {
		case item := <-nbc.recv:
			fmt.Println("memCache收到回收的需要回收的内存")
			items.PushBack(noBuffferObj{
				b:    item,
				used: time.Now(),
			})
			fmt.Println("item len:", items.Len())
		case nbc.send <- e.Value.(noBuffferObj).b:
			oldLen := items.Len()
			items.Remove(e)
			newLen := items.Len()
			fmt.Println("oldLen:", oldLen, "newLen:", newLen)
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
					fmt.Println("free memCache obj:", item.Value.(noBuffferObj).used.Unix())
					items.Remove(item)
					item.Value = nil
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
	timeout := time.NewTicker(time.Second * 2)
	for {
		select {
		case <-timeout.C:
			fmt.Println("free memCacheOnce timeout:", time.Now().Unix())
			nbc.freeMem <- 'f'
		}
	}
}
