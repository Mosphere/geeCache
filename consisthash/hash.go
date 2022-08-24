package consisthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func([]byte) uint32

type Map struct {
	hash     Hash
	replicas int            //虚拟节点个数
	keys     []int          //节点hash值得排序数组(从小到大), 为环状结构，
	hashMap  map[int]string //虚拟节点和真实节点映射：虚拟节点hash值为key，真实节点值为value
}

func New(replicas int, fn Hash) *Map {
	m := &Map{
		hash:     fn,
		replicas: replicas,
		hashMap:  make(map[int]string),
	}

	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}

	return m
}

//
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			//虚拟节点格式：真实节点名称_虚拟节点序号
			hash := m.hash([]byte(key + "_" + strconv.Itoa(i)))
			m.keys = append(m.keys, int(hash))
			m.hashMap[int(hash)] = key
		}
	}
	sort.Ints(m.keys)
}

//根据一个真实节点
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}

	hash := int(m.hash([]byte(key)))

	// 当hash比m.keys都大时，idx == len(m.keys), 这时应该选择 m.keys[0]
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})

	//
	index := m.keys[idx%len(m.keys)]
	//根据hash值获取

	return m.hashMap[index]
}
