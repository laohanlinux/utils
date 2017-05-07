package hashmap

import (
	"fmt"
	"testing"

	"github.com/laohanlinux/assert"
	//"strconv"
	"strconv"
)

func TestHashMap(t *testing.T) {
	hm := NewHashMap(16, Hash32Func)
	k1, v1 := []byte("foo"), []byte("bar")
	k2, v2 := []byte("foo"), []byte("cc")
	k3, v2 := []byte("foo1"), []byte("cc")

	hm.AddEntry([]byte("foo"), []byte("bar"))
	value := hm.LookupEntry([]byte("foo"))
	assert.Equal(t, value, v1)
	fmt.Println(string(value))
	//fmt.Printf("%#v\n", hm.nodes)

	//del
	hm.DelEntry(k1)
	assert.Nil(t, hm.LookupEntry(k1))

	// add
	hm.AddEntry(k1, v1)
	assert.Equal(t, hm.LookupEntry(k1), v1)
	hm.AddEntry(k2, v2)
	assert.Equal(t, hm.LookupEntry(k1), v2)
	fmt.Printf("%#v\n", hm.nodes)

	hm.AddEntry(k3, v2)
	fmt.Printf("%#v\n", hm.nodes)
}

func TestHashMap1(t *testing.T) {
	//ky := make(map[string]string)
	hm := NewHashMap(160, Hash32Func)
	for i := 0; i < 1024; i++ {
		k := strconv.Itoa(i)
		hm.AddEntry([]byte(k), []byte(k))
	}
	fmt.Printf("%#v\n", hm.nodes)
	for i := 0; i < 1024; i++ {
		k := strconv.Itoa(i)
		v := hm.LookupEntry([]byte(k))
		//fmt.Println("TL:",string(k))
		assert.NotNil(t, v)
		assert.Equal(t, string(v), k)
	}

	for i := 0; i < 1024; i++ {
		k := strconv.Itoa(i)
		hm.DelEntry([]byte(k))
	}
	for i := 0; i < 1024; i++ {
		k := strconv.Itoa(i)
		v := hm.LookupEntry([]byte(k))
		//fmt.Println("TL:",string(k))
		assert.Nil(t, v)
	}

	fmt.Printf("%#v\n", hm.nodes)

	for i := 0; i < 1024; i++ {
		k := strconv.Itoa(i)
		hm.AddEntry([]byte(k), []byte(k))
	}
	fmt.Printf("%#v\n", hm.nodes)
	for i := 0; i < 1024; i++ {
		k := strconv.Itoa(i)
		v := hm.LookupEntry([]byte(k))
		//fmt.Println("TL:",string(k))
		assert.NotNil(t, v)
		assert.Equal(t, string(v), k)
	}

	for i := 0; i < 1024; i++ {
		k := strconv.Itoa(i)
		hm.DelEntry([]byte(k))
	}
	for i := 0; i < 1024; i++ {
		k := strconv.Itoa(i)
		v := hm.LookupEntry([]byte(k))
		//fmt.Println("TL:",string(k))
		assert.Nil(t, v)
	}

	fmt.Printf("%#v\n", hm.nodes)
}

func TestHashMap2(t *testing.T) {
	k5, v5 := []byte("5"), []byte("5")
	k10, v10 := []byte("10"), []byte("10")
	hm := NewHashMap(16, Hash32Func)
	hm.AddEntry(k5, v5)
	hm.AddEntry(k10, v10)
	assert.Equal(t, hm.LookupEntry([]byte(k10)), v10)
	fmt.Printf("%#v\n", hm.nodes)
}
