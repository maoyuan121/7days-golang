package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// Hash 将 bytes 转换为 uint32
type Hash func(data []byte) uint32

// Map 包含了所有 hash  过的 key
type Map struct {
	hash     Hash           // 函数函数
	replicas int            // 虚拟节点数
	keys     []int          // 哈希环 Sorted, 虚拟节点的哈希值集合
	hashMap  map[int]string // 虚拟节点与真实节点的映射表。 键是虚拟节点的哈希值，值是真实节点的名称
}

// Map 的构造函数
// @param replicas：虚拟节点倍数
// @param fn：hash 函数
func New(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// 将 key 添加到 hash
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
		}
	}

	sort.Ints(m.keys)
}

// Get gets the closest item in the hash to the provided key.
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}

	hash := int(m.hash([]byte(key)))
	// Binary search for appropriate replica.
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})

	return m.hashMap[m.keys[idx%len(m.keys)]]
}
