package geecache

//接口型函数
type Getter interface {
	Get(key string) ([]byte, error)
}

type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error){
	return f(key)
}

func main() {
	var g Getter = GetterFunc(func(key string) ([]byte, error) {

	} )

	g.Get(key)
}


//Group结构
