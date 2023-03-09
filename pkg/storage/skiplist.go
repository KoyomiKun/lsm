package storage

import (
	"fmt"
	"math/rand"
	"sync"
)

type Key []byte
type Value []byte

type Option func(*SkipList)

func (k Key) Less(k2 Key) bool {
	for i, b := range k {
		if i >= len(k2) {
			return false
		}
		if b < k2[i] {
			return true
		}
		if b > k2[i] {
			return false
		}
	}
	return false
}

func (k Key) Equals(k2 Key) bool {
	if len(k) != len(k2) {
		return false
	}
	for i, b := range k {
		if b != k2[i] {
			return false
		}
	}
	return true
}

func WithMaxLevel(level int) Option {
	return func(sl *SkipList) {
		sl.maxLevel = level
	}
}

func WithRatio(ratio int) Option {
	return func(sl *SkipList) {
		sl.ratio = ratio
	}
}

type Node struct {
	key   Key
	value Value
	nxt   []*Node
}

type SkipList struct {
	Node

	tmp      []*Node
	maxLevel int
	ratio    int // 1 / ratio nodes exists on next level

	currentLevel int
	currentLen   int

	rmu *sync.RWMutex
}

func NewSkipList(opts ...Option) *SkipList {
	l := &SkipList{
		maxLevel:     32,
		ratio:        4,
		currentLevel: 0,
		currentLen:   0,

		rmu: &sync.RWMutex{},
	}

	for _, opt := range opts {
		opt(l)
	}

	l.Node.nxt = make([]*Node, l.maxLevel)
	l.tmp = make([]*Node, l.maxLevel)
	return l
}

func (sl *SkipList) Get(key Key) Value {
	sl.rmu.RLock()
	defer sl.rmu.RUnlock()

	prev := &sl.Node

	var nxt *Node
	for i := sl.currentLevel - 1; i >= 0; i-- {
		nxt = prev.nxt[i]
		// find the first key >= target at nxt position
		for nxt != nil && nxt.key.Less(key) {
			prev = nxt
			nxt = prev.nxt[i]
		}
	}
	if nxt != nil && nxt.key.Equals(key) {
		return nxt.value
	}

	return nil
}

func (sl *SkipList) Set(key Key, value Value) {
	sl.rmu.Lock()
	defer sl.rmu.Unlock()

	prev := &sl.Node
	var nxt *Node
	for i := sl.currentLevel - 1; i >= 0; i-- {
		nxt = prev.nxt[i]
		for nxt != nil && nxt.key.Less(key) {
			prev = nxt
			nxt = prev.nxt[i]
		}
		sl.tmp[i] = prev
	}

	if nxt != nil && nxt.key.Equals(key) {
		nxt.value = value
		return
	}

	level := sl.getLevel()

	if level > sl.currentLevel {
		level = sl.currentLevel + 1
		sl.currentLevel = level
		sl.tmp[level-1] = &sl.Node
	}

	node := &Node{
		key:   key,
		value: value,
		nxt:   make([]*Node, level),
	}

	for i := level - 1; i >= 0; i-- {
		node.nxt[i] = sl.tmp[i].nxt[i]
		sl.tmp[i].nxt[i] = node
	}

	sl.currentLen++

}

func (sl *SkipList) Delete(key Key) Value {
	sl.rmu.Lock()
	defer sl.rmu.Unlock()

	prev := &sl.Node
	var nxt *Node
	for i := sl.currentLevel - 1; i >= 0; i-- {
		nxt = prev.nxt[i]
		for nxt != nil && nxt.key.Less(key) {
			prev = nxt
			nxt = prev.nxt[i]
		}

		sl.tmp[i] = prev
	}

	if nxt == nil || !nxt.key.Equals(key) {
		return nil
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
