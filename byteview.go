package geecache

//this byteview need to be immutable
type byteview struct {
	b []byte
}

func (v byteview) Len() int {
	return len(v.b)
}

func (v byteview) ByteSlice() []byte {
	return  cloneBytes(v.b)
	
}

func (v byteview) String() string {
	return string(v.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
} 

