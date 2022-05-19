package share_clipboard

import (
	"encoding/json"
	"io/ioutil"
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

type Config struct {
	ServerAddr          string `json:"server_addr"`
	SendClipboardHotKey string `json:"send_clipboard_hotkey"`
	LoadClipboardHotKey string `json:"load_clipboard_hotkey"`
}

func ReadConfig(path string) (*Config, error) {
	config := new(Config)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
