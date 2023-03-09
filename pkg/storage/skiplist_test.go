package storage

import (
	"math"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetAndGet(t *testing.T) {
	type KV struct {
		K Key
		V Value
	}
	tests := []struct {
		name string
		kvs  []KV
		find Key
		res  Value
	}{
		{
			name: "insert in order",
			kvs: []KV{
				{
					K: Key{1},
					V: Value{1},
				},
				{
					K: Key{2},
					V: Value{2},
				},
				{
					K: Key{3},
					V: Value{3},
				},
				{
					K: Key{4},
					V: Value{4},
				},
				{
					K: Key{5},
					V: Value{5},
				},
			},
			find: Key{3},
			res:  Value{3},
		},
		{
			name: "insert randomly",
			kvs: []KV{
				{
					K: Key{2},
					V: Value{2},
				},
				{
					K: Key{7},
					V: Value{7},
				},
				{
					K: Key{4},
					V: Value{4},
				},
				{
					K: Key{1},
					V: Value{1},
				},
			},
			find: Key{7},
			res:  Value{7},
		},
		{
			name: "element not exist",
			kvs: []KV{
				{
					K: Key{2},
					V: Value{2},
				},
				{
					K: Key{7},
					V: Value{7},
				},
				{
					K: Key{4},
					V: Value{4},
				},
				{
					K: Key{1},
					V: Value{1},
				},
			},
			find: Key{5},
			res:  nil,
		},
		{
			name: "same key",
			kvs: []KV{
				{
					K: Key{2},
					V: Value{2},
				},
				{
					K: Key{1},
					V: Value{1},
				},
				{
					K: Key{4},
					V: Value{4},
				},
				{
					K: Key{1},
					V: Value{2},
				},
			},
			find: Key{1},
			res:  Value{2},
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
					K: Key{1},
					V: Value{1},
				},
				{
					K: Key{2},
					V: Value{2},
				},
				{
					K: Key{3},
					V: Value{3},
				},
				{
					K: Key{4},
					V: Value{4},
				},
				{
					K: Key{5},
					V: Value{5},
				},
			},
			deleted: Key{3},
		},
		{
			name: "delete key not exist",
			kvs: []KV{
				{
					K: Key{2},
					V: Value{2},
				},
				{
					K: Key{7},
					V: Value{7},
				},
				{
					K: Key{4},
					V: Value{4},
				},
				{
					K: Key{1},
					V: Value{1},
				},
			},
			deleted: Key{8},
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
		sk.Set(Key{byte(rand.Uint32())}, Value{byte(rand.Uint32())})
	}
}

func BenchmarkGet(b *testing.B) {
	const N = math.MaxUint8
	sk := NewSkipList(WithMaxLevel(32), WithRatio(4))
	for i := 0; i < N; i++ {
		sk.Set(Key{byte(i)}, Value{byte(i)})
	}
	for i := 0; i < b.N; i++ {
		sk.Get(Key{byte(rand.Uint32() % N)})
	}
}

func BenchmarkDelete(b *testing.B) {
	const N = math.MaxUint8
	sk := NewSkipList(WithMaxLevel(32), WithRatio(4))
	for i := 0; i < N; i++ {
		sk.Set(Key{byte(i)}, Value{byte(i)})
	}
	for i := 0; i < b.N; i++ {
		sk.Delete(Key{byte(rand.Uint32() % N)})
	}
}
