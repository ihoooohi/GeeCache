package geecache

import(
	"net/http"
)

const defaultBasePath = "_geecache"

type HTTPPool struct {
	self string //记录地址，包括ip/主机名 + 端口号
	basePath string //节点通信前缀
}

func NewHTTPPool(s string) *HTTPPool {
	return &HTTPPool{
		self: s,
		basePath: defaultBasePath,
	}
}

