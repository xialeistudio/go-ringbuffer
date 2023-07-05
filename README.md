# RingBuffer

A [RingBuffer](https://en.wikipedia.org/wiki/Circular_buffer) implementation in Go.

## Feature

+ High Performance
+ Automatic grow
+ Not thread-safe
+ Compatible with built-in io interface

## Get Started

```bash
go get -v github.com/xialeistudio/go-ringbuffer
```

```go
rb := ringbuffer.New(8)
// write data
n,err := rb.Write([]byte{'a', 'b', 'c', 'd'})
// read data
data := make([]byte, 4)
n, err = rb.Read(data)
```

## Benchmark

Hardware: Apple M1 Pro

```text
Write 64B chunks BenchmarkWrite-10    	21153045	        55.69 ns/op
```

```text
Write+Read 64B chunks BenchmarkRead-10    	 6038793	       209.2 ns/op
```