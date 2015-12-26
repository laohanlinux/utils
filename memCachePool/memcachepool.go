package memCachePool

import "fmt"

//
const (
	NoBlockingChanType      = 0
	NoBlockingBytesChanType = 1
)

// MemCacheFactory ...
type MemCacheFactory struct {
}

//GetMemCachePool ....
func (m MemCacheFactory) GetMemCachePool(memCacheType, poolSize int) MemCachePool {
	switch memCacheType {
	case NoBlockingChanType:
		fmt.Println("申请一个NoBlockingChanType的内存缓冲")
		return NewNoBlockingChan(poolSize)
	case NoBlockingBytesChanType:
		return NewNoBlockingBytesChan(poolSize)
	default:
		return nil
	}
}

// MemCachePool is a interface for different type memecache
type MemCachePool interface {
}
