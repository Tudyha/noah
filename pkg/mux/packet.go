package mux

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"noah/pkg/utils"

	"github.com/gookit/goutil/arrutil"
)

const (
	magicNumber         = 0x12345678
	version             = 0x01
	compressNone  uint8 = 0x00
	compressGzip  uint8 = 0x01
	compressFlate uint8 = 0x02
)

var (
	compresses []uint8 = []uint8{compressNone, compressGzip, compressFlate}
)

type packet struct {
	flag     flag
	data     []byte
	connId   uint32
	compress uint8
}

type flag uint8

const (
	flagFirst   flag = 0x00
	flagNewConn flag = 0x01
	flagConnOk  flag = 0x02
	flagData    flag = 0x03
	flagClose   flag = 0x04
	flagPing    flag = 0x05
	flagPong    flag = 0x06
)

// readPacket 从 reader 中读取并解析一个完整的数据包
func readPacket(reader io.Reader) (*packet, error) {
	// 读取魔数
	var mn uint32
	err := binary.Read(reader, binary.BigEndian, &mn)
	if err != nil {
		return nil, errors.New("error reading magic number: " + err.Error())
	}

	// 读取版本号
	var ver uint8
	err = binary.Read(reader, binary.BigEndian, &ver)
	if err != nil {
		return nil, errors.New("error reading version: " + err.Error())
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

	// 读取压缩方式
	var compress uint8
	err = binary.Read(reader, binary.BigEndian, &compress)
	if err != nil {
		return nil, errors.New("error reading compress: " + err.Error())
	}
	if !arrutil.Contains(compresses, compress) {
		return nil, errors.New("invalid compress")
	}

	// 读取数据长度
	var length uint32
	err = binary.Read(reader, binary.BigEndian, &length)
	if err != nil {
		return nil, errors.New("error reading length: " + err.Error())
	}

	// 读取数据内容
	data := make([]byte, length)
	_, err = io.ReadFull(reader, data)
	if err != nil {
		return nil, errors.New("error reading data: " + err.Error())
	}

	// 检查魔数和版本号
	if magicNumber != mn {
		return nil, errors.New("invalid magic number")
	}
	if version != ver {
		return nil, errors.New("invalid version")
	}

	// 解压
	switch compress {
	case compressGzip:
		data, err = utils.GzipDecompress(data)
		if err != nil {
			return nil, errors.New("error decompressing data: " + err.Error())
		}
	case compressFlate:
		data, err = utils.FlateDecompress(data)
		if err != nil {
			return nil, errors.New("error decompressing data: " + err.Error())
		}
	}

	return &packet{
		flag:     flag,
		data:     data,
		connId:   connId,
		compress: compress,
	}, nil
}

// buildPacket 构建一个完整的数据包
func buildPacket(flag flag, connId uint32, data []byte, compress uint8) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, uint32(magicNumber))
	if err != nil {
		return nil, err
	}
	err = binary.Write(buf, binary.BigEndian, uint8(version))
	if err != nil {
		return nil, err
	}
	err = binary.Write(buf, binary.BigEndian, flag)
	if err != nil {
		return nil, err
	}
	err = binary.Write(buf, binary.BigEndian, connId)
	if err != nil {
		return nil, err
	}
	if !arrutil.Contains(compresses, compress) {
		return nil, errors.New("invalid compress")
	}
	err = binary.Write(buf, binary.BigEndian, compress)
	if err != nil {
		return nil, err
	}

	// 压缩
	switch compress {
	case compressGzip:
		data, err = utils.GzipCompress(data)
		if err != nil {
			return nil, errors.New("error compressing data: " + err.Error())
		}
	case compressFlate:
		data, err = utils.FlateCompress(data)
		if err != nil {
			return nil, errors.New("error compressing data: " + err.Error())
		}
	}

	err = binary.Write(buf, binary.BigEndian, uint32(len(data)))
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(data)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
