// Package ringbuffer
package ringbuffer

import (
	"io"
)

type RingBuffer struct {
	buf        []byte
	mask       int
	readIndex  int
	writeIndex int
}

// Read reads up to len(p) bytes into p.
func (r *RingBuffer) Read(p []byte) (n int, err error) {
	readable := r.ReadableBytes()
	for i := 0; i < len(p); i++ {
		if readable == 0 {
			return i, io.EOF
		}
		p[i] = r.buf[r.readIndex]
		r.readIndex = (r.readIndex + 1) & r.mask
		readable--
	}
	return len(p), nil
}

// Write writes len(p) bytes from p to the ring buffer.
func (r *RingBuffer) Write(p []byte) (n int, err error) {
	if r.WritableBytes() < len(p) {
		r.growCapacity(len(r.buf) + len(p) - r.WritableBytes())
	}
	r.writeIndex = r.writeIndex & r.mask

	for i := 0; i < len(p); i++ {
		r.buf[r.writeIndex] = p[i]
		r.writeIndex++
	}
	return len(p), nil
}

// ReadableBytes returns the number of bytes that can be read from the ring buffer.
func (r *RingBuffer) ReadableBytes() int {
	if r.writeIndex >= r.readIndex {
		return r.writeIndex - r.readIndex
	}
	return len(r.buf) - r.readIndex + r.writeIndex
}

// WritableBytes returns the number of bytes that can be written to the ring buffer
func (r *RingBuffer) WritableBytes() int {
	if r.writeIndex >= r.readIndex {
		return len(r.buf) - r.writeIndex + r.readIndex
	}
	return r.readIndex - r.writeIndex
}

// growCapacity grows the capacity of the ring buffer
func (r *RingBuffer) growCapacity(newCap int) {
	oldCap := len(r.buf)
	newLen := oldCap << 1
	for newLen < newCap {
		newLen <<= 1
	}
	r.buf = append(r.buf, make([]byte, newLen-oldCap)...)
	r.mask = newLen - 1
}

// Clear resizes the index of the ring buffer, for performance purposes, the underlying buf is not cleared.
func (r *RingBuffer) Clear() {
	r.readIndex = 0
	r.writeIndex = 0
}

// New Create a new ring buffer with specified capacity
func New(size int) *RingBuffer {
	if size < 8 {
		size = 8
	} else if size > 8 {
		idealSize := 8
		for idealSize < size {
			idealSize <<= 1
		}
		size = idealSize
	}
	return &RingBuffer{
		buf:        make([]byte, size),
		mask:       size - 1,
		readIndex:  0,
		writeIndex: 0,
	}
}
