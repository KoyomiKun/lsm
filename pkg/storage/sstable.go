package storage

import (
	"math"

	"lsm/pkg/encoding"
)

const (
	MaxElemPerTable = math.MaxUint64
)

type SSTable struct {
	TableHeader

	Keys   []Key
	Values []Value

	// TODO:  streaming //
	//marshaledHeader []byte
	//marshaledData   []byte
	//marshaledIndex  []byte
}

type Index struct {
	ElemOffset uint64
	ElemSize   uint64
}

type TableHeader struct {
	// TODO: compression //
	DataOffset  uint64
	IndexOffset uint64
	DataSize    uint64
	IndexSize   uint64

	// TODO: support different marshal type //
	//KeyMarshalType
	//ValueMarshalType
}

func (t *SSTable) MarshalData(dst []byte) {
	dst := encoding.MarshalUint64()

}
