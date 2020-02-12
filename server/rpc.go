package main

import (
	"encoding/gob"
	"net"
	"net/rpc"
	"time"

	. "github.com/indeedhat/smon/message"
)

// serveRpc
// Setup the rpc chanel to listen for commands
func serveRpc(path string) (application *Smon) {
	application = &Smon{time.Now()}
	catchError(rpc.Register(application), "Failed to register application")

	listener, err := net.Listen("unix", path)
	catchError(err, "Failed to open socket")

	gob.Register(RpcRequest{})
	gob.Register(Cpu{})
	gob.Register(Memory{})
	gob.Register(Disk{})
	gob.Register(Date{})
	gob.Register(Network{})
	gob.Register(Uptime{})
	gob.Register(PingPong{})
	gob.Register(RpcResponse{})

	go rpc.Accept(listener)
	return
}
