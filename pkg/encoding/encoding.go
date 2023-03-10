package encoding

import (
	"encoding/binary"
	log "lsm/pkg/logger"
)

type MarshalType byte
type TypeHandler func(dst []byte, input interface{}) ([]byte, error)

const (
	MarshalTypeUncompression = MarshalType(0)
)

var (
	marshalTypeToHandler = map[MarshalType]TypeHandler{
		MarshalTypeUncompression: marshalUncompression,
	}
	logger log.Logger
)

func Init() {
	logger = log.GetGlobalLogger().WithModule("encoding")
}

func marshalUncompression(dst []byte, input interface{}) ([]byte, error) {
	switch input.(type) {
	case []byte:
		return append(dst, input.([]byte)...), nil
	case uint64:
		return MarshalUint64(dst, input.(uint64)), nil
	case uint32:
		return MarshalUint32(dst, input.(uint32)), nil
	}

	return nil, ErrUnsupportType
}

// dst for minimize memory allocate
func Marshal(dst []byte, input interface{}, marshalType MarshalType) ([]byte, error) {
	return marshalTypeToHandler[marshalType](dst, input)
}

func MarshalUint32(dst []byte, u uint32) []byte {
	return append(dst, byte(u>>24), byte(u>>16), byte(u>>8), byte(u))
}

func UnmarshalUint32(src []byte) uint32 {
	return binary.BigEndian.Uint32(src)
}

func MarshalUint64(dst []byte, u uint64) []byte {
	return append(dst, byte(u>>56), byte(u>>48), byte(u>>40), byte(u>>32), byte(u>>24), byte(u>>16), byte(u>>8), byte(u))
}

func UnmarshalUint64(src []byte) uint64 {
	return binary.BigEndian.Uint64(src)
}
