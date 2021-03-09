package main

import (
	"flag"
	"tonger/pkg"
)

var addr = flag.String("a", "127.0.0.1:10501", "server address")

func main() {
	pkg.Main(*addr)
}
