package memCachePool

import (
	"container/list"
	"runtime/debug"
	"sync"
	"time"
)

var bytesChanOnce sync.Once

var nbbc *NoBlockingBytesChan

// ThresholdFreeBytesChan is the threshold that frees memory to os.
// LifeTimeChan is the live time of every block NoBlockingBytesChan.
const (
	ThresholdFreeBytesChan = 268435456
	LifeTimeBytesChan      = 180
)

type noBytesObj struct {
	b    chan []byte
	used int64
}

// NoBlockingBytesChan is a no block channel for memory cache.
// the recycle time is 1 minute ;
// the recycle threshold of total memory is 268435456;
// the recycle threshold of ervry block timeout is 3 minutes
type NoBlockingBytesChan struct {
	send      chan chan []byte //
	recv      chan chan []byte //
	freeMem   chan byte        //
	blockSize uint64           //
}

// NewNoBlockingBytesChan for create a no blocking chan with size block
func NewNoBlockingBytesChan(blockSizes ...uint64) *NoBlockingBytesChan {
	var blockSize uint64
	if len(blockSizes) > 0 {
		blockSize = blockSizes[0]
	} else {
		blockSize = 4096
	}
	bytesChanOnce.Do(func() {
		nbbc = &NoBlockingBytesChan{
			send:      make(chan chan []byte),
			recv:      make(chan chan []byte),
			freeMem:   make(chan byte),
			blockSize: blockSize,
		}
		go nbbc.doWork()
		go nbbc.freeOldMemCache()
	})
	return nbbc
}

// SendChan is the only read channel for NoBlockingBytesChan
func (nbbc *NoBlockingBytesChan) SendChan() <-chan chan []byte {
	return nbbc.send
}

// RecycleChan is the recycle channel for NoBlockingBytesChan
func (nbbc *NoBlockingBytesChan) RecycleChan() chan<- chan []byte {
	return nbbc.recv
}

// SetBufferSize used to set no blocking channel into blockSize
func (nbbc *NoBlockingBytesChan) SetBufferSize(blockSize uint64) {
	nbbc.blockSize = blockSize
}

// Very Block is 4kb
func (nbbc *NoBlockingBytesChan) makeBuffer() chan []byte { return make(chan []byte, nbbc.blockSize) }

func (nbbc *NoBlockingBytesChan) doWork() {
	defer func() {
		debug.FreeOSMemory()
	}()

	var freeSize uint64
	items := list.New()
	for {
		if items.Len() == 0 {
			items.PushBack(noBytesObj{
				b:    nbbc.makeBuffer(),
				used: time.Now().Unix(),
			})
		}
		e := items.Front()
		select {
		case item := <-nbbc.recv:
			//must sure clear the dirty data
			for len(item) != 0 {
				<-item
			}
			items.PushBack(noBytesObj{
				b:    item,
				used: time.Now().Unix(),
			})
		case nbbc.send <- e.Value.(noBytesObj).b:
			items.Remove(e)
		case <-nbbc.freeMem:
			// free too old memcached
			item := items.Front()
			freeTime := time.Now().Unix()
			if items.Len() > 1 {
				for item != nil {
					nItem := item.Next()
					if (freeTime - item.Value.(noBytesObj).used) > LifeTimeBytesChan {
						items.Remove(item)
						item.Value = nil
					} else {
						break
					}
					item = nItem
					freeSize += nbbc.blockSize
				}
			}
			// if needed free memory more than ThresholdFreeBytesChan, call the debug.FreeOSMemory
			if freeSize > ThresholdFreeBytesChan {
				debug.FreeOSMemory()
				freeSize = 0
			}
		}
	}
}

func (nbbc *NoBlockingBytesChan) freeOldMemCache() {
	timeout := time.NewTicker(time.Second * 60)
	for {
		select {
		case <-timeout.C:
			nbbc.freeMem <- 'f'
		}
	}
}
