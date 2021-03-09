package pkg

import (
	"strings"
	"sync"
	"time"
	"tonger/pkg/config"
	"tonger/pkg/constant"
	"tonger/pkg/internal/client"
	"tonger/pkg/internal/model"
	"tonger/pkg/internal/server"
	"tonger/pkg/log"
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
	// run server
	broker.goFunc(broker.runServer)
	// run client
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
				// send to heartbeat to other broker
				log.Logger.Info("send to heartbeat")
				var message = model.RPCMessage{MessageType: constant.HeartBeatSignal}
				broker.cs.Send(message)
			}
		}
	}()
	broker.wg.Wait()
}
