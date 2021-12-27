package cachestruct

// Data 作为数据用于存储值，实现了linkedlist的value接口
type Data struct {
	B []byte //实际缓存值，用来保存数据
}

func (v Data) Length() int {
	return len(v.B)
}

func (v Data) ByteSlice() []byte {
	return CloneBytes(v.B)
}

func (v Data) String() string {
	return string(v.B)
}

func CloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
