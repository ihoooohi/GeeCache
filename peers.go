package geecache

type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

//节点取值器
type PeerGetter interface {
	Get(group string, key string) ([]byte, error)
}

