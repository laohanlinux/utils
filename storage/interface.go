package storage

type Storage interface {
	Put(key StorageKey, value StorageValue) error
	Get(key StorageKey) (value StorageValue, err error)
}

type Iterator interface {
	Next() bool
}

type StorageKey = []byte
type StorageValue = []byte

// type StorageValue interface {
// 	Read()
// 	Write()
// }

// type StorageKey interface{}

// type Int struct {
// 	int
// }

// func (i *Int) Read(b [4]byte) {

// }
// func (i *Int) Write() {}
