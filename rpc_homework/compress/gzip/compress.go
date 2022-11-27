package gzip

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
)

type Compresser struct {
}

func (c *Compresser) Compress(in []byte) ([]byte, error) {
	if len(in) == 0 {
		return in, nil
	}
	buffer := &bytes.Buffer{}
	w := gzip.NewWriter(buffer)
	_, err := w.Write(in)
	if err != nil {
		return nil, err
	}
	err = w.Close()
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func (c *Compresser) Decompress(in []byte) ([]byte, error) {
	if len(in) == 0 {
		return in, nil
	}
	br := bytes.NewReader(in)
	r, err := gzip.NewReader(br)
	if err != nil {
		return nil, err
	}
	out, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Compresser) Code() byte {
	return 1
}

func NewCompresser() *Compresser {
	return &Compresser{}
}
