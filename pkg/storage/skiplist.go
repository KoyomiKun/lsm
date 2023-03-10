package storage

import (
	"fmt"
	"math/rand"
	"sync"
)

type Serializable interface {
	Serialize() []byte
	Deserialize([]byte)
}

type Option[K Serializable, V Serializable] func(*SkipList[K, V])

func WithMaxLevel[K Serializable, V Serializable](level int) Option[K, V] {
	return func(sl *SkipList[K, V]) {
		sl.maxLevel = level
	}
}

func WithRatio[K Serializable, V Serializable](ratio int) Option[K, V] {
	return func(sl *SkipList[K, V]) {
		sl.ratio = ratio
	}
}

type Node[K Serializable, V Serializable] struct {
	key   K
	value V
	nxt   []*Node[K, V]
}

type SkipList[K Serializable, V Serializable] struct {
	Node[K, V]

	tmp      []*Node[K, V]
	maxLevel int
	ratio    int // 1 / ratio nodes exists on next level

	currentLevel int
	currentLen   int

	comparator func(k1 K, k2 K) int

	rmu *sync.RWMutex
}

func NewSkipList[K Serializable, V Serializable](opts ...Option[K, V]) *SkipList[K, V] {
	l := &SkipList[K, V]{
		maxLevel:     32,
		ratio:        4,
		currentLevel: 0,
		currentLen:   0,

		rmu: &sync.RWMutex{},
	}

	for _, opt := range opts {
		opt(l)
	}

	l.Node.nxt = make([]*Node[K, V], l.maxLevel)
	l.tmp = make([]*Node[K, V], l.maxLevel)
	return l
}

func (sl *SkipList[K, V]) Get(key K) interface{} {
	sl.rmu.RLock()
	defer sl.rmu.RUnlock()

	prev := &sl.Node

	var nxt *Node[K, V]
	for i := sl.currentLevel - 1; i >= 0; i-- {
		nxt = prev.nxt[i]
		// find the first key >= target at nxt position
		for nxt != nil && sl.comparator(nxt.key, key) == 1 {
			prev = nxt
			nxt = prev.nxt[i]
		}
	}
	if nxt != nil && sl.comparator(nxt.key, key) == 0 {
		return nxt.value
	}

	return nil
}

func (sl *SkipList[K, V]) Set(key K, value V) {
	sl.rmu.Lock()
	defer sl.rmu.Unlock()

	prev := &sl.Node
	var nxt *Node[K, V]
	for i := sl.currentLevel - 1; i >= 0; i-- {
		nxt = prev.nxt[i]
		for nxt != nil && sl.comparator(nxt.key, key) == 1 {
			prev = nxt
			nxt = prev.nxt[i]
		}
		sl.tmp[i] = prev
	}

	if nxt != nil && sl.comparator(nxt.key, key) == 0 {
		nxt.value = value
		return
	}

	level := sl.getLevel()

	if level > sl.currentLevel {
		for i := sl.currentLevel; i < level; i++ {
			sl.tmp[i] = &sl.Node
		}
		sl.currentLevel = level
	}

	node := &Node[K, V]{
		key:   key,
		value: value,
		nxt:   make([]*Node[K, V], level),
	}

	for i := level - 1; i >= 0; i-- {
		node.nxt[i] = sl.tmp[i].nxt[i]
		sl.tmp[i].nxt[i] = node
	}

	sl.currentLen++

}

func (sl *SkipList[K, V]) Delete(key K) V {
	sl.rmu.Lock()
	defer sl.rmu.Unlock()

	prev := &sl.Node
	var nxt *Node[K, V]
	for i := sl.currentLevel - 1; i >= 0; i-- {
		nxt = prev.nxt[i]
		for nxt != nil && sl.comparator(nxt.key, key) == 1 {
			prev = nxt
			nxt = prev.nxt[i]
		}

		sl.tmp[i] = prev
	}

	if nxt == nil || sl.comparator(nxt.key, key) != 0 {
		var zero V
		return zero // ugly
	}

	for i, v := range nxt.nxt {
		if sl.tmp[i].nxt[i] == nxt {
			sl.tmp[i].nxt[i] = v
			if sl.Node.nxt[i] == nil {
				sl.currentLen--
			}
		}
	}

	sl.currentLen--
	return nxt.value
}

func (sl *SkipList) Len() int {
	return sl.currentLen
}

func (sl *SkipList) Print() {
	fmt.Printf("current len: %d, current level: %d\n", sl.currentLen, sl.currentLevel)
	for i := 0; i < sl.currentLevel; i++ {
		nxt := sl.Node.nxt[i]
		for nxt != nil {
			fmt.Printf("%v,", nxt.value)
			nxt = nxt.nxt[i]
		}
		fmt.Println()
	}
}

func (sl *SkipList) getLevel() int {
	i := 1
	for ; i < sl.maxLevel; i++ {
		if rand.Int()%sl.ratio != 0 {
			break
		}
	}
	return i
}

func (sl *SkipList) Iter() *SkipListIter {
	return &SkipListIter{
		currentNode: &sl.Node,
	}

}

// READ ONLY!
type SkipListIter struct {
	currentNode *Node
}

func (it *SkipListIter) Key() Comparable {
	if it.valid() {
		return it.currentNode.key
	}
	return nil
}

func (it *SkipListIter) Value() interface{} {
	if it.valid() {
		return it.currentNode.value
	}
	return nil
}

func (it *SkipListIter) valid() bool {
	return it.currentNode != nil
}

// if current position is nil, return false; else return true
func (it *SkipListIter) Next() bool {
	if !it.valid() {
		return false
	}
	it.currentNode = it.currentNode.nxt[0]
	return true
}
