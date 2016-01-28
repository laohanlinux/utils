/*=============================================================================

 Copyright (C) 2016 All rights reserved.

 Author: Rg

 Email: daimaldd@gmail.com

 Last modified: 2016-01-28 10:33

 Filename: memcachepool.go

 Description:

=============================================================================*/

//
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
