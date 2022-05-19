package share_clipboard

import (
	"sync"
	"unsafe"
)

const (
	ActionTypeCopy = iota
	ActionTypePaste
)

const (
	StatusOK = "200 OK\n"
)

const (
	defaultBufSize = 1024*1024*10 + 1
)

var bufferPool sync.Pool

func GetBuffer() []byte {
	return bufferPool.Get().([]byte)
}

func PutBuffer(b []byte) {
	bufferPool.Put(b)
}

func init() {
	bufferPool.New = func() interface{} {
		return make([]byte, defaultBufSize)
	}
}

func ToString(p []byte) string {
	return *(*string)(unsafe.Pointer(&p))
}

func ToBytes(str string) []byte {
	return *(*[]byte)(unsafe.Pointer(&str))
}
