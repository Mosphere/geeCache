package geeCache

//缓存值的抽象和封装
type ByteView struct {
	b []byte //byte支持任意的数据类型
}

func (v ByteView) Len() int {
	return len(v.b)
}

func (v *ByteView) String() string {
	return string(v.b)
}
func (v *ByteView) ByteSlice() []byte {
	return clone(v.b)
}

func clone(b []byte) []byte {
	c := make([]byte, len(b))
	copy(b, c)
	return c
}
