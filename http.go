package geecache

import (
	"fmt"
	"geecache/consistenthash"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	pb "geecache/geecachepb"
	"google.golang.org/protobuf/proto"
)

const ( 
	defaultBasePath = "/_geecache/"
	defaultReplicas = 50
)

type HTTPPool struct {
	self string //记录地址，包括ip/主机名 + 端口号
	basePath string //节点通信前缀
	mu sync.Mutex
	peers *consistenthash.HashRing
	httpGetters map[string]*httpGetter

}

type httpGetter struct {
	baseURL string
}
// -----------------------初始化------------------
func NewHTTPPool(s string) *HTTPPool {
	return &HTTPPool{
		self: s,
		basePath: defaultBasePath,
	}
}

//新建/更新httppool中的哈希环
func (h *HTTPPool) Set(peers ...string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.peers = consistenthash.New(defaultReplicas, nil)
	h.peers.Add(peers...)
	h.httpGetters = make(map[string]*httpGetter, len(peers))
	for _, peer := range peers {
		h.httpGetters[peer] = &httpGetter{baseURL: peer + h.basePath}
	}
}
//---------------------------------------------

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

	//使用protobuf
	body, err := proto.Marshal(&pb.Response{Value: value.ByteSlice()})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")

	//裸传字节
	// w.Write(value.ByteSlice())

	//传protobuf
	w.Write(body)
	

}



//挑选节点
func (h *HTTPPool) PickPeer(key string) (PeerGetter, bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if peer := h.peers.Get(key); peer != "" && peer != h.self {
		return h.httpGetters[peer], true
	}

	return nil, false
}



//---------------用protobuf前--------------------

// func (h *httpGetter) Get(group string, key string) ([]byte, error) {
// 	u := fmt.Sprintf(
// 		"%s%s/%s",
// 		h.baseURL,
// 		url.QueryEscape(group),
// 		url.QueryEscape(key),
// 	)

// 	res, err:= http.Get(u)
// 	if err != nil {
// 		return nil, err
// 	}

// 	defer res.Body.Close()

// 	if res.StatusCode != http.StatusOK {
// 		return nil, fmt.Errorf("server returned:%s", res.Status)
// 	}

// 	bytes, err := io.ReadAll(res.Body)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return bytes, nil

// }

//----------------用protobuf后---------------

func (h *httpGetter) Get(in *pb.Request, out *pb.Response) error {
	u := fmt.Sprintf(
		"%v%v/%v",
		h.baseURL,
		url.QueryEscape(in.GetGroup()),
		url.QueryEscape(in.GetKey()),
	)
	res, err := http.Get(u)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned: %v", res.Status)
	}

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %v", err)
	}

	if err = proto.Unmarshal(bytes, out); err != nil {
		return fmt.Errorf("decoding response body: %v", err)
	}

	return nil
}


//检查HTTPPool有没有实现PeerPicker接口
var _ PeerPicker = 	(*HTTPPool)(nil)