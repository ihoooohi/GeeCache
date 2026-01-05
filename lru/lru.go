package lru

import (
	"container/list"
)

type Cache struct {
	maxbytes  int64
	nbytes    int64
	ll        *list.List
	cache     map[string]*list.Element
	OnEvicted func(key string, val CacheValue)
}

type CacheValue interface {
	Len() int
}

type entry struct {
	key string
	val CacheValue
}

func New(maxbytes int64, OnEvicted func(string, CacheValue)) *Cache {
	return &Cache{
		maxbytes:  maxbytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: OnEvicted,
	}
}

// 查找
// 1.返回entry的值
// 2.将节点element移到队尾
func (c *Cache) Get(key string) (CacheValue, bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.val, true
	}
	return nil, false
}

// 删除
// 1.删除链表节点
// 2.删除map[key]
// 3.更新nbytes
// 4.触发回调函数
func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.nbytes -= int64(len(kv.key)) + int64(kv.val.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.val)
		}
	}

}

// 增加&修改
func (c *Cache) Add(key string, val CacheValue) {
	if ele, ok := c.cache[key]; ok {
		kv := ele.Value.(*entry)
		c.nbytes += int64(val.Len()) - int64(kv.val.Len())
		kv.val = val
		c.ll.MoveToFront(ele)
	} else {
		e2 := &entry{
			key: key,
			val: val,
		}
		ele := c.ll.PushFront(e2)
		c.nbytes += int64(e2.val.Len()) + int64(len(e2.key))
		c.cache[key] = ele
	}
	//maxbytes == 0 是无限缓存模式
	for c.maxbytes != 0 && c.nbytes > c.maxbytes {
		c.RemoveOldest()
	}

}

// 计算Cache有多少条记录
func (c *Cache) CacheLen() int {
	return c.ll.Len()
}
