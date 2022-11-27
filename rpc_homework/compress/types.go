package compress

type Compresser interface {
	Compress(in []byte) ([]byte, error)
	Decompress(in []byte) ([]byte, error)
	Code() byte
}
