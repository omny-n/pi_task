package main

import (
	"context"
	"flag"

	"github.com/golang/glog"
	"github.com/omny-n/pi_task/server"
)

var (
	addr    = flag.String("addr", ":9090", "endpoint of the gRPC service")
	network = flag.String("network", "tcp", "a valid network type which is consistent to -addr")
)

func main() {
	flag.Parse()
	defer glog.Flush()

	ctx := context.Background()
	if err := server.Run(ctx, *network, *addr); err != nil {
		glog.Fatal(err)
	}
}
