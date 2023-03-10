package storage

import (
	"math"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetAndGet(t *testing.T) {
	type KV struct {
		K IntKey
		V int
	}
	tests := []struct {
		name string
		kvs  []KV
		find IntKey
		res  interface{}
	}{
		{
			name: "insert in order",
			kvs: []KV{
				{
					K: 1,
					V: 1,
				},
				{
					K: 2,
					V: 2,
				},
				{
					K: 3,
					V: 3,
				},
				{
					K: 4,
					V: 4,
				},
				{
					K: 5,
					V: 5,
				},
			},
			find: 3,
			res:  3,
		},
		{
			name: "insert randomly",
			kvs: []KV{
				{
					K: 2,
					V: 2,
				},
				{
					K: 7,
					V: 7,
				},
				{
					K: 4,
					V: 4,
				},
				{
					K: 1,
					V: 1,
				},
			},
			find: 7,
			res:  7,
		},
		{
			name: "element not exist",
			kvs: []KV{
				{
					K: 2,
					V: 2,
				},
				{
					K: 7,
					V: 7,
				},
				{
					K: 4,
					V: 4,
				},
				{
					K: 1,
					V: 1,
				},
			},
			find: 5,
			res:  nil,
		},
		{
			name: "same key",
			kvs: []KV{
				{
					K: 2,
					V: 2,
				},
				{
					K: 1,
					V: 1,
				},
				{
					K: 4,
					V: 4,
				},
				{
					K: 1,
					V: 2,
				},
			},
			find: 1,
			res:  2,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sk := NewSkipList(WithMaxLevel(32), WithRatio(4))
			for _, kv := range test.kvs {
				sk.Set(kv.K, kv.V)
			}

			assert.Equal(t, test.res, sk.Get(test.find))
		})
	}
}

func TestSetAndDelete(t *testing.T) {
	type KV struct {
		K Key
		V Value
	}
	tests := []struct {
		name    string
		kvs     []KV
		deleted Key
	}{
		{
			name: "delete key",
			kvs: []KV{
				{
					K: 1,
					V: 1,
				},
				{
					K: 2,
					V: 2,
				},
				{
					K: 3,
					V: 3,
				},
				{
					K: 4,
					V: 4,
				},
				{
					K: 5,
					V: 5,
				},
			},
			deleted: 3,
		},
		{
			name: "delete key not exist",
			kvs: []KV{
				{
					K: 2,
					V: 2,
				},
				{
					K: 7,
					V: 7,
				},
				{
					K: 4,
					V: 4,
				},
				{
					K: 1,
					V: 1,
				},
			},
			deleted: 8,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sk := NewSkipList(WithMaxLevel(32), WithRatio(4))
			for _, kv := range test.kvs {
				sk.Set(kv.K, kv.V)
			}
			sk.Delete(test.deleted)

			assert.Equal(t, Value(nil), sk.Get(test.deleted))
		})
	}
}

func BenchmarkSet(b *testing.B) {
	sk := NewSkipList(WithMaxLevel(32), WithRatio(4))
	for i := 0; i < b.N; i++ {
		sk.Set(rand.Uint32(), rand.Uint32())
	}
}

func BenchmarkGet(b *testing.B) {
	const N = math.MaxUint8
	sk := NewSkipList(WithMaxLevel(12), WithRatio(4))
	for i := 0; i < N; i++ {
		sk.Set(i, i)
	}
	for i := 0; i < b.N; i++ {
		sk.Get(rand.Uint32() % N)
	}
}

func BenchmarkDelete(b *testing.B) {
	const N = math.MaxUint8
	sk := NewSkipList(WithMaxLevel(32), WithRatio(4))
	for i := 0; i < N; i++ {
		sk.Set(i, i)
	}
	for i := 0; i < b.N; i++ {
		sk.Delete(rand.Uint32() % N)
	}
}
