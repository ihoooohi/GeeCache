package geecache

import pb "geecache/geecachepb"

type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

//节点取值器
type PeerGetter interface {
	Get(in *pb.Request, out *pb.Response) error 
}

