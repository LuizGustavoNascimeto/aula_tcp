package buffer

import (
	"bufio"
	"net"
	"sync"
)

type ReaderPool struct {
	readers sync.Map
}

var (
	pool *ReaderPool
	once sync.Once
)

func GetReaderPool() *ReaderPool {
	once.Do(func() {
		pool = &ReaderPool{}
	})
	return pool
}

func (p *ReaderPool) Reader(conn net.Conn) *bufio.Reader {
	if reader, ok := p.readers.Load(conn); ok {
		return reader.(*bufio.Reader)
	}

	reader := bufio.NewReader(conn)
	actual, _ := p.readers.LoadOrStore(conn, reader)
	return actual.(*bufio.Reader)
}

func (p *ReaderPool) Remove(conn net.Conn) {
	p.readers.Delete(conn)
}
