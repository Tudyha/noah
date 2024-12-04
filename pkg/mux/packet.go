package mux

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

const (
	MAGIC_NUMBER = 0x12345678
	VERSION      = 0x01
)

type packet struct {
	Flag        flag
	Data        []byte
	MagicNumber uint32
	Version     uint8
	Length      uint32
	ConnId      uint32
}

type flag uint8

const (
	Flag_First    flag = 0x00
	Flag_New_Conn flag = 0x01
	Flag_Conn_Ok  flag = 0x02
	Flag_Data     flag = 0x03
	Flag_Close    flag = 0x04
	Flag_Ping     flag = 0x05
	Flag_Pong     flag = 0x06
)

// readPacket 从 reader 中读取并解析一个完整的数据包
func readPacket(reader io.Reader) (*packet, error) {
	// 读取魔数
	var magicNumber uint32
	err := binary.Read(reader, binary.BigEndian, &magicNumber)
	if err != nil {
		return nil, errors.New("error reading magic number: " + err.Error())
	}

	// 读取版本号
	var version uint8
	err = binary.Read(reader, binary.BigEndian, &version)
	if err != nil {
		return nil, errors.New("error reading version: " + err.Error())
	}

	// 读取数据长度
	var length uint32
	err = binary.Read(reader, binary.BigEndian, &length)
	if err != nil {
		return nil, errors.New("error reading length: " + err.Error())
	}

	// 读取标志位
	var flag flag
	err = binary.Read(reader, binary.BigEndian, &flag)
	if err != nil {
		return nil, errors.New("error reading flag: " + err.Error())
	}

	// 读取连接ID
	var connId uint32
	err = binary.Read(reader, binary.BigEndian, &connId)
	if err != nil {
		return nil, errors.New("error reading session ID: " + err.Error())
	}

	// 读取数据内容
	data := make([]byte, length)
	_, err = io.ReadFull(reader, data)
	if err != nil {
		return nil, errors.New("error reading data: " + err.Error())
	}

	// 检查魔数和版本号
	if magicNumber != MAGIC_NUMBER {
		return nil, errors.New("invalid magic number")
	}
	if version != VERSION {
		return nil, errors.New("invalid version")
	}

	return &packet{
		MagicNumber: magicNumber,
		Version:     version,
		Length:      length,
		Flag:        flag,
		Data:        data,
		ConnId:      connId,
	}, nil
}

// buildPacket 构建一个完整的数据包
func buildPacket(flag flag, connId uint32, data []byte) ([]byte, error) {
	packet := &packet{
		MagicNumber: MAGIC_NUMBER,
		Version:     VERSION,
		Length:      uint32(len(data)),
		Flag:        flag,
		Data:        data,
		ConnId:      connId,
	}

	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, packet.MagicNumber)
	if err != nil {
		return nil, err
	}
	err = binary.Write(buf, binary.BigEndian, packet.Version)
	if err != nil {
		return nil, err
	}
	err = binary.Write(buf, binary.BigEndian, packet.Length)
	if err != nil {
		return nil, err
	}
	err = binary.Write(buf, binary.BigEndian, packet.Flag)
	if err != nil {
		return nil, err
	}
	err = binary.Write(buf, binary.BigEndian, packet.ConnId)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(packet.Data)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func ReadFirstPacket(reader io.Reader) ([]byte, error) {
	p, err := readPacket(reader)
	if p.Flag != Flag_First {
		return nil, errors.New("invalid first packet")
	}
	return p.Data, err
}

func BuildFirstPacket(data []byte) ([]byte, error) {
	return buildPacket(Flag_First, 0, data)
}
