package pkg

import (
	"tonger/pkg/client"
	"tonger/pkg/config"
	"tonger/pkg/constant"
	"tonger/pkg/log"
	"tonger/pkg/model"
	"tonger/pkg/server"
	"strings"
	"sync"
	"time"
)

type Broker struct {
	wg           sync.WaitGroup
	closedSignal chan bool
	s            *server.TCPServer
	cs           *client.TCPClients
}

func NewBroker(addr string, cluster string) *Broker {
	return &Broker{
		wg:           sync.WaitGroup{},
		closedSignal: make(chan bool, 1),
		s:            server.NewTCPServer(addr),
		cs:           client.NewTCPClients(strings.Split(cluster, ";")),
	}
}

func (broker *Broker) Stop() {
	broker.closedSignal <- true
}

func (broker *Broker) goFunc(f func()) {
	broker.wg.Add(1)
	go func() {
		defer broker.wg.Done()
		f()
	}()
}

func (broker *Broker) runServer() {
	broker.s.Run()
}

func (broker *Broker) stopServer() {
	broker.s.Stop()
}

func (broker *Broker) runClients() {
	broker.cs.Run()
}

func (broker *Broker) stopClients() {
	broker.cs.Stop()
}

func (broker *Broker) Run() {
	broker.goFunc(broker.runServer)
	broker.goFunc(broker.runClients)
	d := time.Duration(config.Conf.HeartbeatTime * 1e6)
	ticker := time.NewTicker(d)
	go func() {
		for {
			select {
			case <-broker.closedSignal:
				broker.stopServer()
				broker.stopClients()
				return
			case <-ticker.C:
				log.Logger.Info("send to heartbeat")
				var message = model.RPCMessage{MessageType: constant.HeartBeatSignal}
				errChan := make(chan error, 10)
				broker.cs.SendAsync(message, errChan)
			}
		}
	}()
	broker.wg.Wait()
}
