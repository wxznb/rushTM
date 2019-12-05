package main

import (
	"container/list"
	"fmt"
	"sync"
)

type LRUCache struct {
	cap   int
	cache map[interface{}]*list.Element
	list  *list.List
	mutex sync.Mutex
}

type Pair struct {
	key   string
	value string
}

func New(cap int) LRUCache {
	return LRUCache{
		cap:   cap,
		cache: make(map[interface{}]*list.Element),
		list:  list.New(),
	}
}

func (this *LRUCache) Get(key string) string {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	if elem, ok := this.cache[key]; ok {
		this.list.MoveToBack(elem)
		return elem.Value.(Pair).value
	}
	return ""
}

func (this *LRUCache) Put(key string, value string) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	if elem, ok := this.cache[key]; ok {
		elem.Value = Pair{key, value}
		this.list.MoveToBack(elem)
	} else {
		if this.list.Len() >= this.cap {
			delete(this.cache, this.list.Front().Value.(Pair).key)
			this.list.Remove(this.list.Front())
		}
		this.list.PushBack(Pair{key, value})
		this.cache[key] = this.list.Back()
	}
}

func main() {
	lru := New(3)
	lru.Get("hello")
	lru.Put("ok", "hello")
	fmt.Println(lru.Get("ok"))
}
