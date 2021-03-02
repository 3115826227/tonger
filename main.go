package main

import (
	"flag"
	"tonger/pkg"
	"tonger/pkg/config"
)

var addr = flag.String("a", "127.0.0.1:10501", "server address")

func main() {
	flag.Parse()
	broker := pkg.NewBroker(*addr, config.Conf.Cluster)
	broker.Run()
}
