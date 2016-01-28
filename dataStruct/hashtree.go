/*=============================================================================

 Copyright (C) 2016 All rights reserved.

 Author: Rg

 Email: daimaldd@gmail.com

 Last modified: 2016-01-28 10:32

 Filename: hashtree.go

 Description:

=============================================================================*/

//
package datastruct

import (
	"hash/crc32"
	"sync"

	"github.com/laohanlinux/go-logger/logger"
)

const (
	// HashTreeSize ...
	HashTreeSize = 32
)

// HashKeyType ...
type HashKeyType int

// HashFunc ...
type HashFunc func([]byte) HashKeyType

// HashTree ...
type HashTree struct {
	hp     [10]HashKeyType
	hf     HashFunc
	root   *node
	locker *sync.RWMutex
}

// NewHashTree ...
func NewHashTree(key []byte, value []byte, hf HashFunc) *HashTree {
	ht := &HashTree{
		hp:     [10]HashKeyType{2, 3, 5, 7, 11, 13, 17, 19, 23, 29},
		hf:     hf,
		root:   nil,
		locker: &sync.RWMutex{},
	}

	nKey := ht.hf(key)
	ht.root = &node{
		data:   value,
		key:    nKey,
		leaf:   nil,
		exists: true,
	}
	return ht
}

// Insert ...
func (htree *HashTree) Insert(key []byte, value []byte) {
	nKey := htree.hf(key)
	htree.locker.Lock()
	insert(nKey, value, htree.root, htree.hp, 0)
	htree.locker.Unlock()
}

// Get
// if return nil indicate no found the key's value
func (htree *HashTree) Get(key []byte) []byte {
	nKey := htree.hf(key)
	htree.locker.RLock()
	data := get(nKey, htree.root, htree.hp, 0)
	htree.locker.RUnlock()
	//logger.Info(data)
	value, _ := data.([]byte)
	return value
}

// Del ...
func (htree *HashTree) Del(key []byte) {
	nKey := htree.hf(key)
	htree.locker.Lock()
	del(nKey, htree.root, htree.hp, 0)
	htree.locker.Unlock()
}

type node struct {
	data   interface{}
	key    HashKeyType
	leaf   []*node
	exists bool
}

func insert(key HashKeyType, value interface{}, n *node, hp [10]HashKeyType, level int) {
	if n == nil {
		logger.Info(key, value)
		n = &node{
			data:   value,
			key:    key,
			exists: true,
		}
		return
	}

	if n.key == key {
		n.key = key
		n.data = value
		n.exists = true
		return
	}

	// if not exists
	if n.exists == false {
		n.key = key
		n.data = value
		n.exists = true
		return
	}

	// get level node leaf
	idx := key % hp[level]
	level++
	if n.leaf == nil {
		n.leaf = make([]*node, hp[level], hp[level])
		for i := 0; i < int(hp[level]); i++ {
			n.leaf[i] = &node{
				exists: false,
			}
		}
	}
	insert(key, value, n.leaf[idx], hp, level)
}

// 删除所有和key相等的节点
func del(key HashKeyType, n *node, hp [10]HashKeyType, level int) {
	if n == nil {
		return
	}

	// 查到key
	// 把data也置为nil， 降低内存
	if n.key == key {
		n.exists = false
		n.data = nil
		return
	}

	idx := key % hp[level]

	// 如果有字节点，则继续检索
	if n.leaf == nil || n.leaf[idx] == nil {
		return
	}
	level++
	del(key, n.leaf[idx], hp, level)
}

func get(key HashKeyType, n *node, hp [10]HashKeyType, level int) interface{} {

	if n == nil {
		return nil
	}
	if n.key == key {
		// logger.Debug("get:", key)
		if n.exists {
			return n.data
		}
		return nil
	}

	idx := key % hp[level]
	if n.leaf == nil || n.leaf[idx] == nil {
		return nil
	}

	level++
	return get(key, n.leaf[idx], hp, level)
}

func hashfunc(key []byte) HashKeyType {
	return HashKeyType(crc32.ChecksumIEEE(key))
}
