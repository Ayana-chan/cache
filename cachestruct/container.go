package cachestruct

// Container 作为数据的容器用于存储值，实现了linkedlist的value接口
type Container struct {
	b []byte //实际缓存值，用来保存数据
}

func (v Container) Length() int {
	return len(v.b)
}

func (v Container) ByteSlice() []byte {
	return cloneBytes(v.b)
}

func (v Container) String() string {
	return string(v.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
