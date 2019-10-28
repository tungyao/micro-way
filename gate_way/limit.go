package gate_way

import (
	"log"
	"net"
	"sync"
)

type limiter struct {
	net.Listener
	accept    chan struct{}
	closeOnce sync.Once
	close     chan struct{}
}
type Config struct {
	MaxConn int `默认为2000`
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

func (l *limitListenerConn) Close() error { //取出accept中的一个值,这里是继承的 net.Conn 接口, 不一样!!!
	err := l.Conn.Close()
	l.shutdownOnce.Do(l.shutdown)
	return err
}
