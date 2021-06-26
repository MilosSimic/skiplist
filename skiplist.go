package skiplist

import (
	"errors"
	"math/rand"
	"sync"
	"time"
)

type SkipList struct {
	maxHeight int
	height    int
	size      int
	head      *SkipListNode
	mutex     *sync.RWMutex
}

type SkipListNode struct {
	key       string
	value     []byte
	timestamp int64
	tombstone bool
	next      []*SkipListNode
}

type Entry struct {
	key       string
	value     []byte
	timestamp int64
	tombstone bool
}

func (s *SkipListNode) Key() string {
	return s.key
}

func (s *SkipListNode) Value() []byte {
	return s.value
}

func (s *SkipListNode) Timestamp() int64 {
	return s.timestamp
}

func (s *SkipListNode) Tombstone() bool {
	return s.tombstone
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

func (s *SkipList) Add(key string, value []byte) Entry {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	level := s.roll()
	newNode := &SkipListNode{
		key:       key,
		value:     value,
		timestamp: time.Now().Unix(),
		tombstone: false,
		next:      make([]*SkipListNode, level+1), // if level is 0 than we need one more for pointers
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
	return Entry{
		key:       key,
		value:     value,
		tombstone: false,
		timestamp: newNode.timestamp,
	}
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

func (s *SkipList) Get(key string) (Entry, error) {
	rez := s.search(key)
	if rez != nil && !rez.tombstone {
		return Entry{
			key:       rez.key,
			value:     rez.value,
			tombstone: rez.tombstone,
			timestamp: rez.timestamp,
		}, nil
	}
	return Entry{}, errors.New("Not existing key")
}

func (s *SkipList) TombstoneIt(key string) (Entry, error) {
	rez := s.search(key)
	if rez != nil {
		rez.value = nil
		rez.timestamp = time.Now().Unix()
		rez.tombstone = true
		s.size--
		return Entry{
			key:       rez.key,
			value:     nil,
			tombstone: true,
			timestamp: rez.timestamp,
		}, nil
	}
	return Entry{}, errors.New("Not existing key")
}

func (s *SkipList) Update(key string, value []byte) (Entry, error) {
	rez := s.search(key)
	if rez != nil {
		rez.value = value
		rez.timestamp = time.Now().Unix()
		return Entry{
			key:       rez.key,
			value:     rez.value,
			tombstone: false,
			timestamp: rez.timestamp,
		}, nil
	}
	return Entry{}, errors.New("Not existing key")
}

func (s *SkipList) Remove(key string) (Entry, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	curr := s.head
	for i := s.height; i >= 0; i-- {
		for ; curr.next[i] != nil; curr = curr.next[i] {
			if curr.next[i].key > key {
				break
			} else if curr.next[i].key == key {
				del := curr.next[i]
				curr.next[i] = curr.next[i].next[i]
				s.size--
				return Entry{
					key:       del.key,
					value:     del.value,
					tombstone: del.tombstone,
					timestamp: del.timestamp,
				}, nil
			}
		}
	}
	return Entry{}, errors.New("Not existing key")
}

func (s *SkipList) Size() int {
	return s.size
}

func (s *SkipList) ToMap(list []*SkipListNode, seen map[string]Entry) {
	for _, n := range list {
		if n == nil {
			continue
		} else {
			if _, ok := seen[n.key]; !ok {
				seen[n.key] = Entry{
					key:       n.key,
					value:     n.value,
					timestamp: n.timestamp,
					tombstone: n.tombstone,
				}
				s.ToMap(n.next, seen)
			} else {
				continue
			}
		}
	}
}
