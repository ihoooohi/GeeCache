package geecache

import "sync"

//接口型函数
type Getter interface {
	Get(key string) ([]byte, error)
}

type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error){
	return f(key)
}

//Group结构
type Group struct {
	name string
	getter Getter
	mainCache cache
}

var(
	mu sync.RWMutex
	groups = make(map[string]*Group)
)

func NewGroup(name string, getter Getter, mainCache cache) *Group {
	if getter == nil {
		panic("nil getter")
	}

	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name: name,
		getter: getter,
		mainCache: mainCache,
	}
	groups[name] = g
	return g
}

func GetGroup(name string) *Group {
	mu.RLock()
	defer mu.RUnlock()
	g := groups[name]
	return g
}

