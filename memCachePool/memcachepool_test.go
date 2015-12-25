package MemCachePool

import (
	"fmt"
	"github.com/laohanlinux/assert"
	"testing"
)

func TestMemCachePool(t *testing.T) {
	memCacheFactory := MemCacheFactory{}
	memCacheObj := memCacheFactory.GetMemCachePool(NoBlokingChanBytesType, 1024)
	assert.NotNil(t, memCacheObj)
	fmt.Printf("%#v\n", memCacheObj)
}
