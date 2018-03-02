package storage

type Storage interface {
	Put(key StorageKey, value StorageValue) error
	Get(key StorageKey) (value StorageValue, err error)
}

type Iterator interface {
	Next() bool
}

type StorageKey interface {
	// AsBytes() []byte
	// Encode() []byte
	// Decode() StorageKey
}

type StorageValue interface {
	// AsBytes() []byte
	// Encode() []byte
	// Decode() StorageValue
}

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
