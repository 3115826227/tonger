package server

import (
	"fmt"
	"net"
	"tonger/pkg/log"
)

type TCPServer struct {
	Addr        string
	closeSignal chan bool
	listener    net.Listener
}

func NewTCPServer(addr string) *TCPServer {
	return &TCPServer{
		Addr:        addr,
		closeSignal: make(chan bool, 1),
	}
}

func (server *TCPServer) Stop() {
	server.closeSignal <- true
}

func (server *TCPServer) Run() {
	listener, err := net.Listen("tcp", server.Addr)
	if err != nil {
		return
	}
	log.Logger.Info("tcp server start successful: " + server.Addr)
	for {
		select {
		case <-server.closeSignal:
			log.Logger.Info("tcp server stop: " + server.Addr)
			return
		default:
			conn, err := listener.Accept()
			if err != nil {
				continue
			}
			go server.handle(conn)
		}
	}
}

func (server *TCPServer) handle(conn net.Conn) {
	defer conn.Close()
	for {
		var buf = make([]byte, 10)
		n, err := conn.Read(buf)
		if err != nil {
			log.Logger.Error("conn read error" + err.Error())
			return
		}
		log.Logger.Info(fmt.Sprintf("read %d bytes, content is %v", n, string(buf[:n])))
	}
}
