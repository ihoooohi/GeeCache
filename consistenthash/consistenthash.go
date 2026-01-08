package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type hash func(data []byte) uint32

type HashRing struct {
	hash        hash
	replicas    int
	tokens      []int          //sorted
	tokenToNode map[int]string //为了负载均衡，防止一个节点挂载太多，引入虚拟节点也就是token来解决
}

func New( replicas int, hash hash) *HashRing {
	m := &HashRing{
		hash: hash,
		replicas: replicas,
		tokens: make([]int, 0),
		tokenToNode: make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

func (h *HashRing) Add(nodes ...string) {
	for _, node := range nodes {
		for i := 0; i < h.replicas; i++ {
			token := int(h.hash([]byte((node + strconv.Itoa(i)))))
			h.tokens = append(h.tokens, token)
			h.tokenToNode[token] = node
		}
	}
	sort.Ints(h.tokens)
}

func (h *HashRing) Get(key string) string {
	if len(key) == 0 {
		return ""
	}

	keyHashInt := int(h.hash([]byte(key)))
	idx := sort.Search(len(h.tokens), func(i int) bool {
		return h.tokens[i] > keyHashInt
	})

	return h.tokenToNode[h.tokens[idx%len(h.tokens)]]
}
