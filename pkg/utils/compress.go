package utils

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"io"
	"sync"
)

var gzipWriterPool = sync.Pool{
	New: func() interface{} {
		return gzip.NewWriter(io.Discard)
	},
}

var gzipReaderPool = sync.Pool{
	New: func() interface{} {
		return new(gzip.Reader)
	},
}

var flateWriterPool = sync.Pool{
	New: func() interface{} {
		w, _ := flate.NewWriter(io.Discard, flate.DefaultCompression)
		return w
	},
}

func GzipCompress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	w := gzipWriterPool.Get().(*gzip.Writer)
	w.Reset(&buf)
	_, err := w.Write(data)
	if err != nil {
		gzipWriterPool.Put(w)
		return nil, err
	}
	err = w.Close()
	if err != nil {
		gzipWriterPool.Put(w)
		return nil, err
	}
	gzipWriterPool.Put(w)
	return buf.Bytes(), nil
}

func GzipDecompress(data []byte) ([]byte, error) {
	buf := bytes.NewBuffer(data)
	r := gzipReaderPool.Get().(*gzip.Reader)
	err := r.Reset(buf)
	if err != nil {
		gzipReaderPool.Put(r)
		return nil, err
	}
	defer func() {
		r.Close()
		gzipReaderPool.Put(r)
	}()
	return io.ReadAll(r)
}

func FlateCompress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	w := flateWriterPool.Get().(*flate.Writer)
	w.Reset(&buf)
	_, err := w.Write(data)
	if err != nil {
		flateWriterPool.Put(w)
		return nil, err
	}
	err = w.Close()
	if err != nil {
		flateWriterPool.Put(w)
		return nil, err
	}
	flateWriterPool.Put(w)
	return buf.Bytes(), nil
}

func FlateDecompress(data []byte) ([]byte, error) {
	buf := bytes.NewBuffer(data)
	r := flate.NewReader(buf)
	defer func() {
		r.Close()
	}()
	return io.ReadAll(r)
}
