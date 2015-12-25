package memCache

import (
	"github.com/laohanlinux/memCache/memCachePool"
)

//
const (
	NoBlokingBytesType     = 0
	NoBlokingChanBytesType = 1
)

// MemCacheFactory ...
type MemCacheFactory struct {
}

//GetMemCachePool ....
func (m MemCacheFactory) GetMemCachePool(memCacheType, poolSize int) MemCachePool {
	switch memCacheType {
	case NoBlokingBytesType:
		return memCachePool.NewNoBlockingBytesChan(poolSize)
	case NoBlokingChanBytesType:
		return memCachePool.NewNoBlockingChan(poolSize)
	default:
		return nil
	}
}

// MemCachePool is a interface for different type memecache
type MemCachePool interface {
}
