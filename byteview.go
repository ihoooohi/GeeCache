package geecache

//this Byteview need to be immutable
type Byteview struct {
	b []byte
}

func (v Byteview) Len() int {
	return len(v.b)
}

func (v Byteview) ByteSlice() []byte {
	return  cloneBytes(v.b)
	
}

func (v Byteview) String() string {
	return string(v.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
} 

