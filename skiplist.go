package skiplist

import (
	"math/rand"
	"sync"
)

type SkipList struct {
	maxHeight int
	height    int
	size      int
	head      *SkipListNode
	mutex     *sync.RWMutex
}

type SkipListNode struct {
	key     string
	value   []byte
	deleted bool
	next    []*SkipListNode
}

type Value struct {
	value   []byte
	deleted bool
}

func (s *SkipListNode) Key() string {
	return s.key
}

func (s *SkipListNode) Value() []byte {
	return s.value
}

func New(maxHeight int, seed int64) *SkipList {
	rand.Seed(seed)
	return &SkipList{
		maxHeight: maxHeight,
		size:      0,
		height:    0,
		head: &SkipListNode{
			key:   "",
			value: nil,
			next:  make([]*SkipListNode, maxHeight),
		},
		mutex: &sync.RWMutex{},
	}
}

func (s *SkipList) roll() int {
	level := 0 // alwasy start from level 0

	// We roll until we don't get 1 from rand function and we did not
	// outgrow maxHeight. BUT rand can give us 0, and if that is the case
	// than we will just increase level, and wait for 1 from rand!
	for ; rand.Int31n(2) == 1; level++ {
		if level > s.height {
			// When we get 1 from rand function and we did not
			// outgrow maxHeight, that number becomes new height
			s.height = level
			return level
		}
	}
	return level
}

func (s *SkipList) Add(key string, value []byte) string {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	level := s.roll()
	newNode := &SkipListNode{
		key:   key,
		value: value,
		next:  make([]*SkipListNode, level+1), // if level is 0 than we need one more for pointers
	}

	curr := s.head
	// This loop goes over height of pointers
	// starting from head and goes down level by level
	for i := s.height; i >= 0; i-- {

		// This loop follow the pointers...
		// If there is no next pointer => next is nil than exit this loop
		// If there is next pointer, test key if smaller than go lowet and test again
		for ; curr.next[i] != nil; curr = curr.next[i] {

			// Key is lower, than break tis loop,
			// go down the pointers and test next, than test key
			if curr.next[i].key > key {
				break
			}
		}

		// We might not insert on that level so we need to be sure!
		if i <= level {
			newNode.next[i] = curr.next[i]
			curr.next[i] = newNode
			s.size++
		}
	}
	return key
}

func (s *SkipList) search(key string) *SkipListNode {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	curr := s.head
	for i := s.height; i >= 0; i-- {
		for ; curr.next[i] != nil; curr = curr.next[i] {
			if curr.next[i].key > key {
				break
			} else if curr.next[i].key == key {
				return curr.next[i]
			}
		}
	}
	return nil
}

func (s *SkipList) Contains(key string) bool {
	if s.search(key) == nil {
		return false
	}
	return true
}

func (s *SkipList) Get(key string) []byte {
	rez := s.search(key)
	if rez != nil {
		return rez.value
	}
	return nil
}

func (s *SkipList) Remove(key string) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	curr := s.head
	founded := false
	for i := s.height; i >= 0; i-- {
		for ; curr.next[i] != nil; curr = curr.next[i] {
			if curr.next[i].key > key {
				break
			} else if curr.next[i].key == key {
				curr.next[i] = curr.next[i].next[i]
				founded = true
			}
		}
	}
	return founded
}

func (s *SkipList) Size() int {
	return s.size
}

func (s *SkipList) ToMap(list []*SkipListNode, seen map[string][]*SkipListNode) {
	for _, n := range list {
		if n == nil {
			continue
		} else {
			if _, ok := seen[n.key]; !ok {
				seen[n.key] = n.next
				s.ToMap(n.next, seen)
			} else {
				continue
			}
		}
	}
}

func (s *SkipList) Prep(list []*SkipListNode, seen map[string]Value) {
	for _, n := range list {
		if n == nil {
			continue
		} else {
			if _, ok := seen[n.key]; !ok {
				seen[n.key] = Value{
					n.value,
					n.deleted,
				}
				s.Prep(n.next, seen)
			} else {
				continue
			}
		}
	}
}
