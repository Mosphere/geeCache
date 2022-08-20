package lru

import (
	"container/list"
)

type Cache struct {
	maxBytes  int64
	nBytes    int64
	ll        *list.List //doubly linked list
	elements  map[string]*list.Element
	OnEvicted func(key string, val Value)
}

type Value interface {
	Len() int
}

type entry struct {
	key   string
	value Value
}

func New(max int64, evicted func(key string, val Value)) *Cache {
	return &Cache{
		maxBytes:  max,
		ll:        list.New(),
		elements:  make(map[string]*list.Element),
		OnEvicted: evicted,
	}
}

func (c *Cache) Get(key string) (Value, bool) {
	if ele, ok := c.elements[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return nil, false
}

func (c *Cache) Remove() {
	ele := c.ll.Back()
	if ele != nil {
		kv := ele.Value.(*entry) //方便通过key删除
		c.ll.Remove(ele)
		delete(c.elements, kv.key)
		c.nBytes -= int64(len(kv.key) + kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

func (c *Cache) Add(key string, val Value) {
	if ele, ok := c.elements[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		nBytes := int64(len(kv.key) + kv.value.Len())
		c.nBytes += nBytes
	} else {
		ele := c.ll.PushFront(&entry{key, val})
		c.elements[key] = ele
		c.nBytes += int64(len(key) + val.Len())
	}

	//
	for c.maxBytes != 0 && c.nBytes > c.maxBytes {
		c.Remove()
	}
}

func (c *Cache) Len() int {
	return c.ll.Len()
}
