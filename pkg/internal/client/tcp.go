package client

import (
	"encoding/json"
	"net"
	"tonger/pkg/internal/model"
	"tonger/pkg/log"
)

type TCPClients struct {
	Addrs       []string
	closeSignal chan bool
	clients     map[string]*TCPClient
}

func NewTCPClients(addrs []string) *TCPClients {
	var clients = new(TCPClients)
	clients.Addrs = addrs
	clients.clients = make(map[string]*TCPClient)
	for _, addr := range addrs {
		clients.clients[addr] = NewTCPClient(addr)
	}
	clients.closeSignal = make(chan bool, 1)
	return clients
}

func (clients *TCPClients) Stop() {
	clients.closeSignal <- true
}

func (clients *TCPClients) Send(message model.RPCMessage) {
	for _, client := range clients.clients {
		client.Send(message)
	}
}

func (clients *TCPClients) SendAsync(message model.RPCMessage, errChan chan error) {
	for _, client := range clients.clients {
		go client.Send(message)
	}
}

func (clients *TCPClients) Run() {
	for _, client := range clients.clients {
		go client.Conn()
	}
	for {
		select {
		case <-clients.closeSignal:
			for _, client := range clients.clients {
				client.Stop()
			}
			return
		}
	}
}

type TCPClient struct {
	Addr        string
	closeSignal chan bool
	messageChan chan model.RPCMessage
	data        []byte
	conn        net.Conn
}

func NewTCPClient(addr string) *TCPClient {
	return &TCPClient{
		Addr:        addr,
		closeSignal: make(chan bool, 1),
		messageChan: make(chan model.RPCMessage, 500),
	}
}

func (client *TCPClient) Send(message model.RPCMessage) {
	client.messageChan <- message
}

func (client *TCPClient) Stop() {
	client.closeSignal <- true
}

func (client *TCPClient) Conn() {
	conn, err := net.Dial("tcp", client.Addr)
	if err != nil {
		log.Logger.Error("client tcp conn failed. error : " + err.Error())
		return
	}
	client.conn = conn
	defer client.conn.Close()
	for {
		select {
		case <-client.closeSignal:
			log.Logger.Info("client tcp closed")
			return
		case message := <-client.messageChan:
			client.data, err = json.Marshal(message)
			if err != nil {
				log.Logger.Error(err.Error())
				continue
			}
			log.Logger.Info(string(client.data))
			if _, err = client.conn.Write(client.data); err != nil {
				log.Logger.Error(err.Error())
				return
			}
		}
	}
}
