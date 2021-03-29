package lru

import "container/list"

// LRU = Least Recently Used 最近最少使用。 It is not safe for concurrent access.
type Cache struct {
	maxBytes int64 // 允许缓存的最大字节数
	nbytes   int64 // 所有缓存数据的字节数
	ll       *list.List
	cache    map[string]*list.Element

	// 可选的，在条目被清除时执行。
	OnEvicted func(key string, value Value)
}

type entry struct {
	key   string
	value Value
}

// Value use Len to count how many bytes it takes
// Value 定义了 Len() 函数，其用于计算他占用了多少个字节
type Value interface {
	Len() int
}

// New 是 Cache  的构造函数
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// 添加一个值到缓存
// 如果已经存在，将该元素移动到最前面，否则添加改元素到第一个位置。
// 如果缓存的内容的大小（nbytes）超过了允许的最大字节数（maxBytes），那么删除最老的一个缓存
func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nbytes += int64(len(key)) + int64(value.Len())
	}
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}
}

// 根据 key 获取值
// 如果获取到了，该元素移动到链表的最前面
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

// 移除最老的  cache
// 重新计算缓存总字节数
// 如果定义了 OnEvicted，那么触发它
func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// 返回缓存中元素的的个数
func (c *Cache) Len() int {
	return c.ll.Len()
}
