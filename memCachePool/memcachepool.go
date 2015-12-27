package memCachePool

//
const (
	NoBlockingChanType      = 0
	NoBlockingBytesChanType = 1
)

// MemCacheFactory ...
type MemCacheFactory struct {
}

// GetMemCachePool ....
func (m MemCacheFactory) GetMemCachePool(memCacheType, poolSize uint64) MemCachePool {
	switch memCacheType {
	case NoBlockingChanType:
		return NewNoBlockingChan(poolSize)
	case NoBlockingBytesChanType:
		return NewNoBlockingBytesChan(poolSize)
	default:
		return nil
	}
}

// MemCachePool is a interface for different type memecache
type MemCachePool interface {
	SetBufferSize(uint64)
}
