package packet

import (
	"bufio"
	"encoding/binary"
	"io"
	"sync"
)

type MessageType uint16

const (
	// 通用消息
	MessageType_Unknown MessageType = iota
	MessageType_Error

	// 业务消息 - 鉴权相关
	MessageType_Login = iota + 1000
	MessageType_LoginAck

	// 业务消息 - 多路复用相关
	MessageType_Stream_Create = iota + 2000
	MessageType_Stream_CreateAck
	MessageType_Stream_Data
	MessageType_Stream_Close
)

type Codec interface {
	Decode(r io.Reader) (*Packet, error)
	Encode(w io.Writer, p *Packet) (int, error)
	Release()
}

var pool = sync.Pool{
	New: func() any {
		return &ProtoCodec{}
	},
}

type Packet struct {
	MessageType MessageType
	length      uint64
	Data        []byte
}

// protobuf编解码器
type ProtoCodec struct {
}

func NewCodec() Codec {
	codec := pool.Get().(*ProtoCodec)
	return codec
}

func (c *ProtoCodec) Decode(r io.Reader) (*Packet, error) {
	p := new(Packet)
	buf := bufio.NewReader(r)
	if err := binary.Read(buf, binary.BigEndian, &p.MessageType); err != nil {
		return nil, err
	}
	if err := binary.Read(buf, binary.BigEndian, &p.length); err != nil {
		return nil, err
	}
	p.Data = make([]byte, p.length)
	if _, err := io.ReadFull(buf, p.Data); err != nil {
		return nil, err
	}
	return p, nil
}

func (c *ProtoCodec) Encode(w io.Writer, p *Packet) (int, error) {
	p.length = uint64(len(p.Data))
	if err := binary.Write(w, binary.BigEndian, p.MessageType); err != nil {
		return 0, err
	}
	if err := binary.Write(w, binary.BigEndian, p.length); err != nil {
		return 0, err
	}
	return int(p.length), binary.Write(w, binary.BigEndian, p.Data)
}

func (c *ProtoCodec) Release() {
	pool.Put(c)
}
