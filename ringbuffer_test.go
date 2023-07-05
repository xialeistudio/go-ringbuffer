package ringbuffer

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	a := assert.New(t)
	t.Run("size < 8", func(t *testing.T) {
		rb := New(7)
		a.Equal(8, len(rb.buf))
		a.Equal(7, rb.mask)
	})
	t.Run("size = 8", func(t *testing.T) {
		rb := New(8)
		a.Equal(8, len(rb.buf))
		a.Equal(7, rb.mask)
	})
	t.Run("size > 8", func(t *testing.T) {
		rb := New(63)
		a.Equal(64, len(rb.buf))
		a.Equal(63, rb.mask)
	})
}

func TestRingBuffer_Operate(t *testing.T) {
	a := assert.New(t)
	t.Run("empty buffer", func(t *testing.T) {
		rb := New(8)
		data := make([]byte, 1)
		n, err := rb.Read(data)
		a.Equal(io.EOF, err)
		a.Equal(0, n)
	})
	t.Run("no available data", func(t *testing.T) {
		rb := New(8)
		data := []byte{1}
		n, err := rb.Write(data)
		a.Nil(err)
		a.Equal(1, n)

		newData := make([]byte, 4)
		n, err = rb.Read(newData)
		a.Equal(1, n)
		a.Equal(io.EOF, err)
		a.Equal(byte(1), newData[0])
	})
	t.Run("readIndex <= writeIndex", func(t *testing.T) {
		rb := New(8)
		for i := 0; i < 8; i++ {
			n, err := rb.Write([]byte{'a'})
			a.Nil(err)
			a.Equal(1, n)
		}
		a.Equal(8, rb.ReadableBytes())
		a.Equal(0, rb.WritableBytes())
		data := make([]byte, 4)
		n, err := rb.Read(data)
		a.Nil(err)
		a.Equal("aaaa", string(data))
		a.Equal(4, n)
		a.Equal(8, rb.writeIndex)
		a.Equal(4, rb.readIndex)
		a.Equal(4, rb.ReadableBytes())
		a.Equal(4, rb.WritableBytes())
	})
	t.Run("readIndex > writeIndex", func(t *testing.T) {
		rb := New(8)
		for i := 0; i < 8; i++ {
			n, err := rb.Write([]byte{'a'})
			a.Nil(err)
			a.Equal(1, n)
		}
		a.Equal(0, rb.WritableBytes())
		data := make([]byte, 6)
		n, err := rb.Read(data)
		a.Nil(err)
		a.Equal("aaaaaa", string(data))
		a.Equal(6, n)
		a.Equal(6, rb.readIndex)
		a.Equal(8, rb.writeIndex)
		a.Equal(6, rb.WritableBytes())
		// write new data
		n, err = rb.Write([]byte{'b', 'b', 'b', 'b'})
		a.Nil(err)
		a.Equal(4, n)
		a.Equal(4, rb.writeIndex)
		// read data
		data = make([]byte, 6)
		n, err = rb.Read(data)
		a.Nil(err)
		a.Equal(6, n)
		a.Equal([]byte{'a', 'a', 'b', 'b', 'b', 'b'}, data)
		a.Equal(4, rb.readIndex)
		a.Equal(8, rb.WritableBytes())
	})
}

func TestRingBuffer_WritableBytes(t *testing.T) {
	a := assert.New(t)

	rb := New(8)
	n, err := rb.Write([]byte{'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a'})
	a.Nil(err)
	a.Equal(8, n)

	data := make([]byte, 6)
	n, err = rb.Read(data)
	a.Nil(err)
	a.Equal(6, n)
	a.Equal([]byte{'a', 'a', 'a', 'a', 'a', 'a'}, data)
	a.Equal(6, rb.WritableBytes())

	n, err = rb.Write([]byte{'b', 'b', 'b', 'b'})
	a.Nil(err)
	a.Equal(4, n)
	a.Equal(2, rb.WritableBytes())
}

func TestRingBuffer_Write(t *testing.T) {
	a := assert.New(t)
	rb := New(8)
	n, err := rb.Write([]byte{'a', 'a', 'a', 'a', 'a', 'a'})
	a.Nil(err)
	a.Equal(6, n)
	a.Equal(2, rb.WritableBytes())
	a.Equal(6, rb.ReadableBytes())

	n, err = rb.Write([]byte{'b', 'b', 'b', 'b', 'b', 'b', 'b', 'b', 'b', 'b', 'b', 'b', 'b', 'b', 'b', 'b', 'b', 'b'})
	a.Nil(err)
	a.Equal(18, n)

	a.Equal(32, len(rb.buf))
	a.Equal(24, rb.ReadableBytes())
	a.Equal(8, rb.WritableBytes())

	data := make([]byte, 24)
	n, err = rb.Read(data)
	a.Nil(err)
	a.Equal(24, n)
	a.Equal("aaaaaabbbbbbbbbbbbbbbbbb", string(data))
}

func TestRingBuffer_Clear(t *testing.T) {
	a := assert.New(t)

	rb := New(8)
	n, err := rb.Write([]byte{'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a'})
	a.Nil(err)
	a.Equal(8, n)
	a.Len(rb.buf, 8)
	a.Equal(0, rb.readIndex)
	a.Equal(8, rb.writeIndex)

	rb.Clear()

	a.Len(rb.buf, 8)
	a.Equal(7, rb.mask)
	a.Equal(0, rb.readIndex)
	a.Equal(0, rb.writeIndex)
}

func TestReaderWriter(t *testing.T) {
	a := assert.New(t)
	rb := New(8)
	n, err := io.WriteString(rb, "helloworld")
	a.Nil(err)
	a.Equal(10, n)

	data, err := io.ReadAll(rb)
	a.Nil(err)
	a.Equal("helloworld", string(data))
}

func BenchmarkWrite(b *testing.B) {
	data := make([]byte, 64)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rb := New(64)
		_, _ = rb.Write(data)
	}
}

func BenchmarkRead(b *testing.B) {
	data := make([]byte, 64)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rb := New(64)
		_, _ = rb.Write(data)
		_, _ = rb.Read(data)
	}
}
