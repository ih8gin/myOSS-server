package rs

import (
	"MyOSS/apiServer/objectstream"
	"fmt"
	"github.com/klauspost/reedsolomon"
	"io"
)

const (
	DATA_SHARDS     = 4
	PARITY_SHARDS   = 2
	ALL_SHARDS      = DATA_SHARDS + PARITY_SHARDS
	BLOCK_PER_SHARD = 8000
	BLOCK_SIZE      = BLOCK_PER_SHARD * DATA_SHARDS
)

type RSPutStream struct {
	*encoder
}

func NewRSPutStream(dataServers []string, hash string, size int64) (*RSPutStream, error) {
	if len(dataServers) != ALL_SHARDS {
		return nil, fmt.Errorf("dataServers number mismatch")
	}
	perShard := (size + DATA_SHARDS - 1) / DATA_SHARDS
	writers := make([]io.Writer, ALL_SHARDS)
	var e error
	for i := range writers {
		writers[i], e = objectstream.NewTempPutStream(dataServers[i], fmt.Sprintf("%s.%d", hash, i), perShard)
		if e != nil {
			return nil, e
		}
	}
	enc := NewEncoder(writers)
	return &RSPutStream{enc}, nil
}

type encoder struct {
	writers []io.Writer
	enc     reedsolomon.Encoder
	cache   []byte
}

func NewEncoder(writers []io.Writer) *encoder {
	enc, _ := reedsolomon.New(DATA_SHARDS, PARITY_SHARDS)
	return &encoder{writers, enc, nil}
}

func (enc *encoder) Write(p []byte) (n int, e error) {
	length := len(p)
	current := 0
	for length != 0 {
		next := BLOCK_SIZE - len(enc.cache)
		if next > length {
			next = length
		}
		enc.cache = append(enc.cache, p[current:current+next]...)
		if len(enc.cache) == BLOCK_SIZE {
			enc.Flush()
		}
		current += next
		length -= next
	}
	return len(p), nil
}

func (enc *encoder) Flush() {
	if len(enc.cache) == 0 {
		return
	}
	shards, _ := enc.enc.Split(enc.cache)
	enc.enc.Encode(shards)
	for i := range shards {
		enc.writers[i].Write(shards[i])
	}
	enc.cache = []byte{}
}

func (enc *encoder) Commit(success bool) {
	enc.Flush()
	for i := range enc.writers {
		enc.writers[i].(*objectstream.TempPutStream).Commit(success)
	}
}
