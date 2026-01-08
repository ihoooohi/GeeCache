package geecache

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const defaultBasePath = "/_geecache/"

type HTTPPool struct {
	self string //记录地址，包括ip/主机名 + 端口号
	basePath string //节点通信前缀
}

type httpGetter struct {
	baseURL string
}

func NewHTTPPool(s string) *HTTPPool {
	return &HTTPPool{
		self: s,
		basePath: defaultBasePath,
	}
}

func (h *HTTPPool) Log(format string, v ...any) {
	log.Printf("[Server %s] %s", h.self,fmt.Sprintf(format, v...) )

}

func (h *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, h.basePath) {
		panic(fmt.Errorf("HTTPPool serving unexpected path: %s", r.URL.Path))
	}
	h.Log("%s, %s", r.Method, r.URL.Path)
	//<basePath>/<group>/<key>
	parts := strings.SplitN(r.URL.Path[len(h.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return //为什么要加return，因为http.Error不会停止程序，只会返回响应
	}

	groupName := parts[0]
	key := parts[1]

	group := GetGroup(groupName)
	if group == nil {
		http.Error(w, "no such group: "+groupName, http.StatusNotFound)
		return
	}

	value, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(value.ByteSlice())
	

}

func (h *httpGetter) Get(group string, key string) ([]byte, error) {
	u := fmt.Sprintf(
		"%s%s/%s",
		h.baseURL,
		url.QueryEscape(group),

	)
}