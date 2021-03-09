package pkg

import "tonger/pkg/config"

func Main(addr string) {
	broker := NewBroker(addr, config.Conf.Cluster)
	broker.Run()
}
