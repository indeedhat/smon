package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/rpc"
	"os"

	. "github.com/indeedhat/smon/message"
)

// ./local_client /tmp/foo Brian
func main() {
	client, err := rpc.Dial("unix", os.Args[1])
	if err != nil {
		log.Fatalf("failed: %s", err)
	}

	gob.Register(RpcRequest{})
	gob.Register(Cpu{})
	gob.Register(Memory{})
	gob.Register(Disk{})
	gob.Register(Date{})
	gob.Register(Network{})
	gob.Register(Uptime{})
	gob.Register(PingPong{})

	req := &RpcRequest{Message: PingPong{Message: "ping"}}
	var res RpcResponse

	err = client.Call("Smon.Network", req, &res)
	if err != nil {
		log.Fatalf("error in rpc: %s", err)
	}

	fmt.Println(res.Message)
}
