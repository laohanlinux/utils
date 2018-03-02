package storage

import (
	"errors"

	"github.com/emirpasic/gods/trees/btree"
)

var (
	ErrNotFound = errors.New("not found")
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

func NewPatch() *Patch { return &Patch{cgs: btree.NewWithStringComparator(5)} }

type Patch struct {
	cgs *btree.Tree
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

func (path *Patch) Iterator() interface{} {
	return path.cgs.Iterator()
}

func (patch *Patch) IsEmpty() bool { return patch.cgs.Size() == 0 }

func (patch *Patch) Len() int { return patch.cgs.Size() }

type ForkOp int

const (
	Stored      = ForkOp(0)
	Replaced    = ForkOp(1)
	Inserted    = ForkOp(2)
	Deleted     = ForkOp(3)
	MissDeleted = ForkOp(4)
	Finished    = ForkOp(5)
)

type ForkOperator struct {
	Name   string
	Value  []byte
	Option *Change
}

type Fork struct {
	snapshot  Snapshot
	patch     *Patch
	changeLog *ForkOperator
	logged    bool
}

func (fork *Fork) Get(name string, key StorageKey) (value StorageValue, err error) {
	var found bool
	var change = fork.patch.Changes(name)
	if change == nil {
		err = ErrNotFound
		return
	}
	value, found = change.data.Get(key)
	if !found {
		err = ErrNotFound
	}
	return
}

func (fork *Fork) contains(name string, key StorageKey) (found bool) {
	change := fork.patch.Changes(name)
	if change == nil {
		return false
	}
	_, found = change.data.Get(key)
	return
}

// func (fork *Fork) Iterator(name string, from []byte) interface{} {
// 	change := for
// }

type Snapshot interface {
	Get(name string, key StorageKey) (value StorageValue, err error)
	Contains(name string, key StorageKey) (included bool)
	Iterator() (iterator interface{})
}
