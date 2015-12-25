package memCache

//	"errors"

// var BinaryMemCachePool = errors.New("binarymemcachepool")
//
// var UnKownMemCachePoolType = errors.New("unkown memCachePool type")

// ThresholdFreeOsMemory (256M) for memCache size to free to os
const (
	ThresholdFreeOsMemory = 268435456
)

const (
	NoBlokingBytesType     = 0
	NoBlokingChanBytesType = 1
)

//
type MemCacheFactory struct {
}

func (m MemCacheFactory) GetMemCachePool(memCacheType, poolSize int) MemCachePool {
	switch memCacheType {
	case NoBlokingBytesType:

	case NoBlokingChanBytesType:
	default:
		return
	}

}

// MemCachePool is a interface for different type memecache
type MemCachePool interface {
}

var nbc *NoBlockingChan

// memCache Object
type noBuffferObj struct {
	b    []byte
	used int64
}
