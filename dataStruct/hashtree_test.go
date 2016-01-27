package datastruct

import (
	"fmt"
	"testing"

	"github.com/laohanlinux/assert"
)

func TestHashTree(t *testing.T) {
	sets := make(map[string][]byte)

	sets["a"] = []byte("1")
	sets["b"] = []byte("2")
	sets["c"] = []byte("3")
	sets["d"] = []byte("4")
	sets["e"] = []byte("5")
	sets["f"] = []byte("6")
	sets["g"] = []byte("7")
	sets["h"] = []byte("8")
	sets["i"] = []byte("9")
	sets["g"] = []byte("10")
	sets["k"] = []byte("11")
	sets["l"] = []byte("12")
	sets["m"] = []byte("13")
	sets["n"] = []byte("14")
	sets["o"] = []byte("15")
	sets["p"] = []byte("16")
	sets["q"] = []byte("17")
	sets["r"] = []byte("18")
	sets["广州"] = []byte("是一个好城市")

	hTree := NewHashTree([]byte("a"), sets["a"], hashfunc)

	for k, v := range sets {
		hTree.Insert([]byte(k), v)
	}

	for k, v := range sets {
		value := hTree.Get([]byte(k))
		fmt.Println("v:", string(v), "value:", string(value))
		assert.Equal(t, v, value)
	}

	for k, _ := range sets {
		hTree.Del([]byte(k))
	}

	for k, _ := range sets {
		value := hTree.Get([]byte(k))
		assert.Nil(t, value)
	}

	for k, v := range sets {
		hTree.Insert([]byte(k), v)
	}

	for k, v := range sets {
		value := hTree.Get([]byte(k))
		fmt.Println("v:", string(v), "value:", string(value))
		assert.Equal(t, v, value)
	}

	for k, _ := range sets {
		hTree.Del([]byte(k))
	}

	for k, _ := range sets {
		value := hTree.Get([]byte(k))
		assert.Nil(t, value)
	}
}
