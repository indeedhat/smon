package main

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"net/rpc"
	"time"

	"github.com/indeedhat/smon/message"
)

func init() {
	gob.Register(message.RpcRequest{})
	gob.Register(message.Cpu{})
	gob.Register(message.Memory{})
	gob.Register(message.Disk{})
	gob.Register(message.Date{})
	gob.Register(message.Network{})
	gob.Register(message.Uptime{})
	gob.Register(message.PingPong{})
	gob.Register(message.RpcResponse{})
}

type RpcClient struct {
	client *rpc.Client
	index  int
	config ClientConfig
}

func (r RpcClient) String() string {
	var buffer bytes.Buffer
	writer := bufio.NewWriter(&buffer)

	cpu := clients[r.index].Cpu
	fmt.Fprintf(writer, "Cpu:\n  Total: %d\n  Idle: %d\n\n", cpu.Total, cpu.Idle)

	mem := clients[r.index].Memory
	mod := uint64(1024 * 1024 * 1024)
	fmt.Fprintf(
		writer,
		"Memory:\n  Total: %2d\n  Used: %2d\n  Free: %2d\n\n",
		mem.Total/mod,
		mem.Used/mod,
		mem.Free/mod,
	)

	return buffer.String()
}

func (r RpcClient) loop() {
	for {
		r.cpu()
		r.memory()
		r.date()
		r.uptime()
		r.network()
		r.disk()

		log.Println(clients[r.index])
		time.Sleep(time.Duration(r.config.Interval) * time.Millisecond)
	}
}

func (r RpcClient) cpu() {
	if !r.config.Modules.Cpu {
		return
	}

	cpu, ok := (*r.request("Smon.Cpu", message.Cpu{})).(message.Cpu)
	if !ok {
		log.Println("rpc.cpu: bad response")
		return
	}

	clients[r.index].Cpu = cpu
}

func (r RpcClient) memory() {
	if !r.config.Modules.Memory {
		return
	}

	mem, ok := (*r.request("Smon.Memory", message.Memory{})).(message.Memory)
	if !ok {
		log.Println("rpc.memory: bad response")
		return
	}

	clients[r.index].Memory = mem
}

func (r RpcClient) date() {
	if !r.config.Modules.Date {
		return
	}

	date, ok := (*r.request("Smon.Date", message.Date{})).(message.Date)
	if !ok {
		log.Println("rpc.date: bad response")
		return
	}

	clients[r.index].Date = date
}

func (r RpcClient) network() {
	if !r.config.Modules.Network.Enable {
		return
	}

	network, ok := (*r.request("Smon.Network", message.Network{})).(message.Network)
	if !ok {
		log.Println("rpc.network: bad response")
		return
	}

	clients[r.index].Network = network
}

func (r RpcClient) uptime() {
	if !r.config.Modules.Uptime {
		return
	}

	uptime, ok := (*r.request("Smon.Uptime", message.Uptime{})).(message.Uptime)
	if !ok {
		log.Println("rpc.uptime: bad response")
		return
	}

	clients[r.index].Uptime = uptime
}

func (r RpcClient) disk() {
	if 0 == len(r.config.Modules.Disk) {
		return
	}

	mes := message.Disk{}
	for i := range r.config.Modules.Disk {
		mes.Mounts = append(mes.Mounts, message.Mount{
			Mount: r.config.Modules.Disk[i],
		})
	}

	disk, ok := (*r.request("Smon.Disk", mes)).(message.Disk)
	if !ok {
		log.Println("rpc.disk: bad response")
		return
	}

	clients[r.index].Disk = disk
}

func (r RpcClient) request(key string, mes message.Message) *message.Message {
	req := message.RpcRequest{mes}
	var res message.RpcResponse

	if err := r.client.Call(key, req, &res); nil != err {
		log.Printf("%s: %s\n", key, err)
		return nil
	}

	return &res.Message
}

func (r RpcClient) close() {
	r.client.Close()
}

func NewRpcClient(conn io.ReadWriteCloser, config ClientConfig, index int) RpcClient {
	return RpcClient{
		client: rpc.NewClient(conn),
		index:  index,
		config: config,
	}
}
