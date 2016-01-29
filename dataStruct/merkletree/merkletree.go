package merkletree

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

// MerkleTree ...
type MerkleTree struct {
	TreeHeight uint
	Root       *merkleTreeNode
}

type merkleTreeNode struct {
	hash  string
	left  *merkleTreeNode
	right *merkleTreeNode
}

// NewMerkleTree ...
func NewMerkleTree(data []byte, pieces int) *MerkleTree {
	if data == nil || len(data) == 0 || pieces == 0 {
		return &MerkleTree{
			TreeHeight: 0,
			Root:       nil,
		}
	}

	size := len(data) / pieces
	dataPieces := [][]byte{}
	for i := 0; i < pieces; i++ {
		dataPieces = append(dataPieces, data[i*size:(i*size+1)])
	}
	mt := &MerkleTree{
		TreeHeight: 0,
		Root:       nil,
	}
	mt.BuildMerkleTree(dataPieces)
	return mt
}

func (mt *MerkleTree) DisplayMerkleTreeNodeDetail() {
	display(mt.Root)
}

func display(n *merkleTreeNode) {
	if n == nil {
		return
	}
	if n.left == nil && n.right == nil {
		fmt.Println(n.hash)
		return
	}
	if n.left != nil {
		display(n.left)
	}
	if n.right != nil {
		display(n.right)
	}
}

// BuildMerkleTree ...
func (mt *MerkleTree) BuildMerkleTree(data [][]byte) {
	if data == nil || len(data) == 0 {
		return
	}

	height := uint(1)

	mtns := generateLeaves(data)
	for len(mtns) > 1 {
		mtns = levelUp(mtns)
		height++
	}

	// set root node
	mt.Root = mtns[0]
	mt.TreeHeight = height
}

/**

A
--------------------------
			@
			|
		|-------|
		@		@
		|		|
	|-------|	----|
	@		@		@
	|
----
|
@
--------------------------
B
--------------------------
			@
			|
		|-------|
		@		@
		|		|
	|-------|	----|
	@		@		@
	|
---------
|		|
@		@

---------------------------

1. compare height
2. compare data
3. compare left and right node

**/

// CompareMerkleTree ...
func (mt *MerkleTree) CompareMerkleTree(tmpMT *MerkleTree) bool {
	if mt.TreeHeight != tmpMT.TreeHeight {
		return false
	}

	q1 := []*merkleTreeNode{mt.Root}
	q2 := []*merkleTreeNode{tmpMT.Root}

	for len(q1) > 0 {
		// 1.
		if q1[0].hash != q2[0].hash {
			return false
		}
		//TODO
		// use list to storage mt
		q1 = updateQueue(q1)
		q2 = updateQueue(q2)
	}
	return true
}

// FindDiff ...
func findDiff(mt *MerkleTree, tmpMT *MerkleTree) [][]*merkleTreeNode {
	var mtns [][]*merkleTreeNode

	st1, st2 := diffTree(mt, tmpMT)
	if st1.Root.left == nil && st2.Root.right == nil {
		mtns = append(mtns, []*merkleTreeNode{st1.Root})
	} else {
		leftSt1 := &MerkleTree{Root: st1.Root.left}
		leftSt2 := &MerkleTree{Root: st2.Root.left}
		rightSt1 := &MerkleTree{Root: st1.Root.right}
		rightSt2 := &MerkleTree{Root: st2.Root.right}
		mtns = append(mtns, findDiff(leftSt1, leftSt2)...)
		mtns = append(mtns, findDiff(rightSt1, rightSt2)...)
	}
	return mtns
}

func diffTree(mt1, mt2 *MerkleTree) (st1, st2 *MerkleTree) {
	q1 := []*merkleTreeNode{mt1.Root}
	q2 := []*merkleTreeNode{mt2.Root}

	for len(q1) > 0 {
		if q1[0].hash != q2[0].hash {
			return &MerkleTree{Root: q1[0]}, &MerkleTree{Root: q2[0]}
		}
		q1 = updateQueue(q1)
		q2 = updateQueue(q2)
	}

	return nil, nil
}

/**
			@
			|
		|-------|
		@		@
		|		|
	|-------|	-----
	@		@		@

**/
func levelUp(mtns []*merkleTreeNode) []*merkleTreeNode {
	var nLevel []*merkleTreeNode
	m := md5.New()
	for i := 0; i < len(mtns)/2; i++ {
		hashs := mtns[i*2].hash + mtns[i*2+1].hash
		m.Write([]byte(hashs))
		hash := hex.EncodeToString(m.Sum(nil))
		m.Reset()

		tmpMtn := &merkleTreeNode{
			hash:  hash,
			left:  mtns[i*2],
			right: mtns[i*2+1],
		}
		nLevel = append(nLevel, tmpMtn)
	}
	//if mtns len is even, join the last one
	if len(mtns)%2 == 1 {
		tmpMtn := &merkleTreeNode{
			hash:  mtns[len(mtns)-1].hash,
			right: mtns[len(mtns)-1],
		}
		nLevel = append(nLevel, tmpMtn)
	}

	return nLevel
}

func generateLeaves(data [][]byte) []*merkleTreeNode {
	var leaves []*merkleTreeNode
	m := md5.New()
	for _, d := range data {
		m.Write(d)
		hash := hex.EncodeToString(m.Sum(nil))
		mtn := &merkleTreeNode{
			hash:  hash,
			left:  nil,
			right: nil,
		}
		leaves = append(leaves, mtn)
		m.Reset()
	}

	return leaves
}

func updateQueue(q []*merkleTreeNode) []*merkleTreeNode {
	if q[0].left != nil {
		q = append(q, q[0].left)
	}
	if q[0].right != nil {
		q = append(q, q[0].right)
	}

	if len(q) > 1 {
		return q[1:]
	}
	return []*merkleTreeNode{}
}

func compareTowBytes(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}

	return true
}
