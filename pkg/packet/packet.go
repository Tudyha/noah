package packet

import (
	"bufio"
	"encoding/binary"
	"io"
)

type MessageType uint16

const (
	// 通用消息
	MessageType_Unknown MessageType = iota
	MessageType_Ping
	MessageType_Error

	// 业务消息 - 鉴权相关
	MessageType_Login = iota + 1000
	MessageType_LoginAck
	MessageType_Logout

	// 业务消息 - 多路复用相关
	MessageType_Stream_Create = iota + 2000
	MessageType_Stream_CreateAck
	MessageType_Stream_Data
	MessageType_Stream_Close

	// tunnel 控制消息
	MessageType_Tunnel = iota + 3000
	MessageType_Tunnel_Pty
	MessageType_Tunnel_PtyAck
)

type CodecType uint8

const (
	CodecType_Unkown CodecType = iota
	CodecType_Protobuf
)

type Codec interface {
	Decode(r io.Reader) (*Packet, error)
	Encode(w io.Writer, p *Packet) (int, error)
}

type Packet struct {
	MessageType MessageType
	CodecType   CodecType
	length      uint64
	Data        []byte
}

type codec struct {
}

func NewCodec() Codec {
	return &codec{}
}

func (c *codec) Decode(r io.Reader) (*Packet, error) {
	p := new(Packet)
	buf := bufio.NewReader(r)
	if err := binary.Read(buf, binary.BigEndian, &p.MessageType); err != nil {
		return nil, err
	}
	if err := binary.Read(buf, binary.BigEndian, &p.CodecType); err != nil {
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

func (c *codec) Encode(w io.Writer, p *Packet) (int, error) {
	p.length = uint64(len(p.Data))
	if err := binary.Write(w, binary.BigEndian, p.MessageType); err != nil {
		return 0, err
	}
	if err := binary.Write(w, binary.BigEndian, p.CodecType); err != nil {
		return 0, err
	}
	if err := binary.Write(w, binary.BigEndian, p.length); err != nil {
		return 0, err
	}
	return int(p.length), binary.Write(w, binary.BigEndian, p.Data)
}
