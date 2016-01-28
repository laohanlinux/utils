package merkletree

import (
	"fmt"
	"testing"
)

func TestDiff(t *testing.T) {
	data1 := []byte("abcdefghijklmnopqrstuvwxyz")
	data2 := []byte("abcdefhijklmnopqrstuvwxyz")
	mt1 := NewMerkleTree(data1, 9)
	mt2 := NewMerkleTree(data2, 9)
	difs := findDiff(mt1, mt2)
	fmt.Println(difs)
}
