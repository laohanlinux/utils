package memCachePool

import (
	"container/list"
	"runtime/debug"
	"sync"
	"time"
)

// ThresholdFreeOsMemory (256M) for memCache size to free to os
const (
	ThresholdFreeNoBlockChan = 268435456
	LifeTimeChan             = 180
)

var noblockOnece sync.Once

var nbc *NoBlockingChan

// memCache Object
type noBuffferObj struct {
	b    []byte
	used int64
}

// NoBlockingChan is a no block channel for memory cache.
// the recycle time is 1 minute ;
// the recycle threshold of total memory is 268435456;
// the recycle threshold of ervry block timeout is 3 minutes
type NoBlockingChan struct {
	send      chan []byte //
	recv      chan []byte //
	freeMem   chan byte   //
	blockSize uint64      //
}

// NewNoBlockingChan for create a no blocking chan bytes with size block
func NewNoBlockingChan(blockSizes ...uint64) *NoBlockingChan {
	var blockSize uint64
	if len(blockSizes) > 0 {
		blockSize = blockSizes[0]
	} else {
		blockSize = 4096
	}
	noblockOnece.Do(func() {
		nbc = &NoBlockingChan{
			send:      make(chan []byte),
			recv:      make(chan []byte),
			freeMem:   make(chan byte),
			blockSize: blockSize,
		}
		go nbc.doWork()
		go nbc.freeOldMemCache()
	})
	return nbc
}

// SendChan ...
func (nbc *NoBlockingChan) SendChan() <-chan []byte { return nbc.send }

// RecycleChan ...
func (nbc *NoBlockingChan) RecycleChan() chan<- []byte { return nbc.recv }

// SetBufferSize used to set no blocking channel into blockSize
func (nbc *NoBlockingChan) SetBufferSize(blockSize uint64) {
	nbc.blockSize = blockSize
}

// Very Block is 4kb
func (nbc *NoBlockingChan) makeBuffer() []byte { return make([]byte, nbc.blockSize) }

func (nbc *NoBlockingChan) doWork() {
	defer func() {
		debug.FreeOSMemory()
	}()

	var freeSize uint64
	items := list.New()
	for {
		if items.Len() == 0 {
			items.PushBack(noBuffferObj{
				b:    nbc.makeBuffer(),
				used: time.Now().Unix(),
			})
		}
		e := items.Front()
		select {
		case item := <-nbc.recv:
			items.PushBack(noBuffferObj{
				b:    item,
				used: time.Now().Unix(),
			})
		case nbc.send <- e.Value.(noBuffferObj).b:
			items.Remove(e)
		case <-nbc.freeMem:
			// free too old memcached
			item := items.Front()
			freeTime := time.Now().Unix()
			if items.Len() > 1 {
				for item != nil {
					nItem := item.Next()
					if (freeTime - item.Value.(noBuffferObj).used) > LifeTimeChan {
						items.Remove(item)
						item.Value = nil
					} else {
						break
					}
					item = nItem
					freeSize += nbc.blockSize
				}
			}

			// if needed free memory more than ThresholdFreeOsMemory, call the debug.FreeOSMemory
			if freeSize > ThresholdFreeNoBlockChan {
				debug.FreeOSMemory()
				freeSize = 0
			}
		}
	}
}

func (nbc *NoBlockingChan) freeOldMemCache() {
	timeout := time.NewTicker(time.Second * 60)
	for {
		select {
		case <-timeout.C:
			nbc.freeMem <- 'f'
		}
	}
}
