package gate_way

import (
	"net"
	"sync"

	"./util"
)

var FLOW chan int

type limiter struct {
	net.Listener
	Flow      int64
	accept    chan struct{}
	closeOnce sync.Once
	close     chan struct{}
}
type Config struct {
	MaxConn     int `default value => 4096`
	MaxBuffFlow int `max flow buff regin ,default 4kb ,but you can use bigest`
}

func (l *limiter) wait() bool {
	select {
	case <-l.close:
		return false
	case l.accept <- struct{}{}:
		return true
	}
}
func Limiter(config *Config, listener net.Listener) net.Listener {
	util.CheckConfig(config, Config{
		MaxConn:     4096,
		MaxBuffFlow: 4096,
	})
	FLOW = make(chan int, config.MaxBuffFlow)
	return &limiter{
		Listener: listener,
		accept:   make(chan struct{}, config.MaxConn), // Limit the number of accesses by channel buffer capacity
		close:    make(chan struct{}),
	}
}
func (l *limiter) Accept() (net.Conn, error) {
	t := l.wait()
	// fmt.Println(l.Flow)
	a, err := l.Listener.Accept()
	if err != nil {
		if t {
			<-l.accept
		}

		return nil, err
	}
	return &limitListenerConn{Conn: a, shutdown: func() {
		n := <-FLOW
		if n != 0 {
			l.Flow += int64(n)
		}
		<-l.accept
	}}, nil
}
func (l *limiter) Close() error { // This is used to close the channel, just close once, please note that this is the inherited net.listener interface
	err := l.Listener.Close()
	l.closeOnce.Do(func() {
		close(l.close)
	})
	return err
}

type limitListenerConn struct {
	net.Conn
	shutdownOnce sync.Once
	shutdown     func()
}

func (l *limitListenerConn) Read(b []byte) (n int, err error) {
	n, err = l.Conn.Read(b)
	FLOW <- n
	return n, err
}
func (l *limitListenerConn) Write(b []byte) (n int, err error) {
	FLOW <- len(b)
	return l.Conn.Write(b)
}
func (l *limitListenerConn) Close() error { // Take out a value in accept, here is the inherited net.conn interface, not the same!!!
	err := l.Conn.Close()
	l.shutdownOnce.Do(l.shutdown)
	return err
}
