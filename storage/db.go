package storage

import (
	"github.com/emirpasic/gods/maps/hashmap"
	"github.com/emirpasic/gods/trees/btree"
)

func NewChanges() *Changes {
	return &Changes{data: new(btree.Tree)}
}

type Changes struct {
	data *btree.Tree
}

func (changes *Changes) Iterator() interface{} {
	return changes.data.Iterator
}

type Change struct {
	Put    []byte
	Delete struct{}
}

func NewPatch() *Patch { return &Patch{cgs: hashmap.New()} }

type Patch struct {
	cgs *hashmap.Map
}

func (patch *Patch) Changes(name string) *Changes {
	changes := patch.cgs
	change, ok := changes.Get(name)
	if !ok {
		return nil
	}
	return change.(*Changes)
}

func (patch *Patch) InsertChanges(name string, changes *Changes) {
	patch.cgs.Put(name, changes)
}

func (patch *Patch) IsEmpty() bool { return patch.cgs.Size() == 0 }

func (patch *Patch) Len() int { return patch.cgs.Size() }
