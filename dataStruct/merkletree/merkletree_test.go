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
	//	difs := findDiff(mt1, mt2)
	//	fmt.Println(difs)

	fmt.Println(mt1.Root, mt2)
	mt1.DisplayMerkleTreeNodeDetail()
	fmt.Println("===================")
	mt2.DisplayMerkleTreeNodeDetail()

	fmt.Println("mt1 is equal to mt2: ", mt1.CompareMerkleTree(mt2))
}
