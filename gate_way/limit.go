package gate_way

import (
	"fmt"
	"log"
	"net"
	"sync"
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
	MaxConn     int `default value => 2000`
	MaxBuffFlow int `max flow buff regin ,default 2kb`
}

func (l *limiter) wait() bool {
	select {
	case <-l.close:
		log.Print(l.Addr(), " is close")
		return false
	case l.accept <- struct{}{}:
		//log.Print(l.Addr()," is accept")
		return true
	}

}
func Limiter(config Config, listener net.Listener) net.Listener {
	if config.MaxConn == 0 {
		config.MaxConn = 2000
	}
	if config.MaxBuffFlow == 0 {
		config.MaxBuffFlow = 2048
	}
	FLOW = make(chan int, config.MaxBuffFlow)

	return &limiter{
		Listener: listener,
		accept:   make(chan struct{}, config.MaxConn), //TODO 通过信道缓冲容量,来限制访问数量
		close:    make(chan struct{}),
	}
}
func (l *limiter) Accept() (net.Conn, error) {
	t := l.wait()
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
			fmt.Println(l.Flow)
		}
		<-l.accept
	}}, nil
}
func (l *limiter) Close() error { //这是用来关闭信道,仅仅关闭一次 , 请注意,这是继承的 net.Listener接口
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
func (l *limitListenerConn) Close() error { //取出accept中的一个值,这里是继承的 net.Conn 接口, 不一样!!!
	err := l.Conn.Close()
	l.shutdownOnce.Do(l.shutdown)
	return err
}
