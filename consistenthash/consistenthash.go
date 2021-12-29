package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func(data []byte) uint32

// Map 包含所有节点
type Map struct {
	hash     Hash
	replicas int            //虚拟节点的倍数，代表每个真实节点对应几个虚拟节点
	keys     []int          //经过排序，代表一个哈希环
	hashMap  map[int]string //虚拟节点和真实节点的映射，虚拟节点即为一个hash值，真实节点为一个地址
}

// New 创建一个Map
func New(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {

		//默认的哈希函数
		m.hash = crc32.ChecksumIEEE

	}
	return m
}

// Add 向Map中添加真实节点
func (m *Map) Add(addrs ...string) {
	for _, addr := range addrs {

		//添加虚拟节点
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + addr)))

			//将节点的哈希值加入数组
			m.keys = append(m.keys, hash)

			//将虚拟节点映射到真实节点地址
			m.hashMap[hash] = addr
		}

	}
	sort.Ints(m.keys) //使数组有序
}

// Get 根据key值来确定真实节点地址
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}
	hash := int(m.hash([]byte(key)))

	//寻找大于key的哈希值的最小虚拟节点哈希值
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})

	//取模是因为key的hash可能大于所有节点的hash，这时应定位到第一个节点
	return m.hashMap[m.keys[idx%len(m.keys)]]
}
