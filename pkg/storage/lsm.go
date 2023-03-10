package storage

type MemTable interface {
	Iter() Iterator
}

type Iterator interface {
	Key() Key
	Value() Value
	Next() bool
}

func marshalMemTable(mt MemTable) SSTable {
	iter := mt.Iter()
	for iter.Next() {
		iter.Key()

	}

}
