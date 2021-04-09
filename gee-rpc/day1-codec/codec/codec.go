package codec

import (
	"io"
)

type Header struct {
	ServiceMethod string // 服务名和方法名。格式： "Service.Method"
	Seq           uint64 // 是请求的序号，也可以认为是某个请求的 ID，用来区分不同的请求
	Error         string // 错误信息，客户端置为空，服务端如果如果发生错误，将错误信息置于 Error 中。
}

// 解码编码的接口
type Codec interface {
	io.Closer
	ReadHeader(*Header) error
	ReadBody(interface{}) error
	Write(*Header, interface{}) error
}

// 构造解码编码接口实例的函数
type NewCodecFunc func(io.ReadWriteCloser) Codec

// 编码类型
type Type string

const (
	GobType  Type = "application/gob"
	JsonType Type = "application/json" // not implemented
)

// 编码类型和其解码编码的映射
var NewCodecFuncMap map[Type]NewCodecFunc

func init() {
	NewCodecFuncMap = make(map[Type]NewCodecFunc)
	NewCodecFuncMap[GobType] = NewGobCodec
}
