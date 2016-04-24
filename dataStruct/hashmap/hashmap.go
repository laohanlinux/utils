package hashmap

import (
	"hash/fnv"
	//"fmt"
	"fmt"
)

type HashFunc func(int, []byte) int

/**

|---------|
| bucket1 | ---> [] [] [] [] []
| bucket2 | ---> [] [] [] [] []
| bucket3 | ---> [] [] [] [] []
-----
*/
type HashNode struct {
	key   []byte
	value []byte
	Prev  *HashNode
	Next  *HashNode
}
type HashMap struct {
	buckets  int
	HashFunc HashFunc
	nodes    []*HashNode
}

func NewHashMap(buckets int, hashFunc HashFunc) *HashMap {
	hashMap := new(HashMap)
	hashMap.buckets = buckets
	hashMap.HashFunc = hashFunc
	hashMap.nodes = make([]*HashNode, buckets, buckets)
	return hashMap
}

func (hm *HashMap) LookupEntry(key []byte) []byte {
	node := hm.GetNodeByKey(key)
	if node == nil {
		return nil
	}
	return node.value
}

func (hm *HashMap) AddEntry(key, value []byte) {

	node := HashNode{
		key:   key,
		value: value,
		Prev:  nil,
		Next:  nil,
	}
	bucketNode :=  hm.bucket(key)

	if *bucketNode == nil {
		*bucketNode = &node
		return
	}

//	fmt.Println(string((*bucketNode).key), string(key))
	bucketNode1 := *bucketNode
	node.Next, node.Prev =  bucketNode1, nil
	bucketNode1.Prev = &node
	*bucketNode = &node
//	disply(*bucketNode)
	return
}

// GetNodeByKey return a HashNode
// return value: nil(not fould), not nil(found)
func (hm *HashMap) GetNodeByKey(key []byte) *HashNode {
	bucketNode := *hm.bucket(key)
	for bucketNode != nil {
		if string(bucketNode.key) == string(key) {
			return bucketNode
		}
		bucketNode = bucketNode.Next
	}

	return bucketNode
}

func (hm *HashMap) DelEntry(key []byte) {
	bucketNode := hm.bucket(key)
	for *bucketNode != nil {
		if string((*bucketNode).key) == string(key) {
			break
		}
		*bucketNode = (*bucketNode).Next
	}

	// the bucket root node is empty
	if *bucketNode == nil {
		return
	}

	// del mid or last node
	if (*bucketNode).Prev != nil {
		(*bucketNode).Prev.Next = (*bucketNode).Next
	}

	// clear old data
	(*bucketNode).Next, (*bucketNode).Prev = nil, nil
	*bucketNode = nil
	return
}

// return the buckets node
func (hm *HashMap) bucket(key []byte) **HashNode {
	idx := hm.HashFunc(hm.buckets, key)
	return &hm.nodes[idx]
}

func Hash32Func(buckets int, key []byte) int {
	if key == nil || len(key) == 0 {
		return 0
	}

	hash32 := fnv.New32()
	hash32.Write(key)
	return int(hash32.Sum32()) % buckets
}


func disply(node *HashNode) {
	for node != nil {
		fmt.Println("key:", string(node.key), "value:", node.value)
		node = node.Next
	}
}