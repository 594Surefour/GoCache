package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// Hash maps bytes to uint32
type Hash func(data []byte) uint32

// Consistence constains all hashed keys
type Consistence struct {
	hash     Hash           //哈希函数
	replicas int            //虚拟节点倍数
	keys     []int          // 哈希环
	hashMap  map[int]string //虚拟节点到哈希节点的映射
}

// New creates a Consistence instance
func New(replicas int, fn Hash) *Consistence {
	m := &Consistence{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// Add adds some keys to the hash.
func (c *Consistence) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < c.replicas; i++ {
			hash := int(c.hash([]byte(strconv.Itoa(i) + key)))
			c.keys = append(c.keys, hash)
			c.hashMap[hash] = key
		}
	}
	sort.Ints(c.keys)
}

// Get gets the closest item in the hash to the provided key.
func (c *Consistence) GetPeer(key string) string {
	if len(c.keys) == 0 {
		return ""
	}

	hashValue := int(c.hash([]byte(key)))
	// Binary search for appropriate replica.
	idx := sort.Search(len(c.keys), func(i int) bool {
		return c.keys[i] >= hashValue
	})

	return c.hashMap[c.keys[idx%len(c.keys)]]
}
