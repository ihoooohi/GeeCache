package singleflight

import "sync"

type call struct {
	wg sync.WaitGroup
	value any
	err error
}

type Group struct {
	mu sync.Mutex
	m map[string]*call
}

//Do函数的作用是判断这个请求是自己调用还是等待其他协程
//为什么第二个参数是一个函数？这个函数是用来执行调用操作的，执行具体的业务逻辑用的，只是设计的时候为了拓展性设计成通用函数的样子，方便具体业务注入
func (g *Group) Do(key string, fn func() (any, error)) (any, error){
	g.mu.Lock()
	//先延迟初始化g.m
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	//1.先判断等待情况
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait()
		return c.value, nil
	}
	//2.再判断自己调用情况（第一个goroutine）
	c := new(call)
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	c.value, c.err = fn()
	c.wg.Done()

	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()

	return c.value, c.err

}

