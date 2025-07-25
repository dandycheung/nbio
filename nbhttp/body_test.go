package nbhttp

import (
	"bytes"
	"crypto/rand"
	"io"
	"testing"

	"github.com/lesismal/nbio/mempool"
)

func TestBodyReaderPool(t *testing.T) {
	br := bodyReaderPool.Get().(*BodyReader)
	buf := make([]byte, 10)
	pbuf := &buf
	br.buffers = append(br.buffers, pbuf)
	*br = emptyBodyReader
	bodyReaderPool.Put(br)

	for i := 0; i < 1000; i++ {
		br2 := bodyReaderPool.Get().(*BodyReader)
		if br2.buffers != nil {
			t.Fatal("len>0")
		}
		buf = make([]byte, 10)
		pbuf = &buf
		br2.buffers = append(br.buffers, pbuf)
		*br2 = emptyBodyReader
		bodyReaderPool.Put(br)
	}
}

func TestBodyReader(t *testing.T) {
	engine := NewEngine(Config{
		BodyAllocator: mempool.NewAligned(),
	})
	var (
		b0 []byte
		b1 = make([]byte, 2049)
		b2 = make([]byte, 1132)
		b3 = make([]byte, 11111)
	)
	_, _ = rand.Read(b1)
	_, _ = rand.Read(b2)
	_, _ = rand.Read(b3)

	allBytes := append(b0, b1...)
	allBytes = append(allBytes, b2...)
	allBytes = append(allBytes, b3...)

	newBR := func() *BodyReader {
		br := NewBodyReader(engine)
		_ = br.append(b1)
		_ = br.append(b2)
		_ = br.append(b3)
		return br
	}

	br1 := newBR()
	body1, err := io.ReadAll(br1)
	if err != nil {
		t.Fatalf("io.ReadAll(br1) failed: %v", err)
	}
	if !bytes.Equal(allBytes, body1) {
		t.Fatalf("!bytes.Equal(allBytes, body1)")
	}
	_ = br1.Close()

	br2 := newBR()
	body2 := make([]byte, len(allBytes))
	for i := range body2 {
		_, err := br2.Read(body2[i : i+1])
		if err != nil {
			t.Fatalf("br2.Readbody2[%d:%d] failed: %v", i, i+1, err)
		}
	}
	if !bytes.Equal(allBytes, body2) {
		t.Fatalf("!bytes.Equal(allBytes, body2)")
	}
	_ = br2.Close()
}
