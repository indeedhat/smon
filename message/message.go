package message

import (
	"time"
)

type Message interface {
	MessageType() string
}

type PingPong struct {
	Message string
}

func (p PingPong) MessageType() string {
	return "message.PingPong"
}

type Cpu struct {
	Total uint64
	Idle  uint64
}

func (c Cpu) MessageType() string {
	return "message.Cpu"
}

type Memory struct {
	Total uint64
	Free  uint64
	Used  uint64
}

func (m Memory) MessageType() string {
	return "message.Memory"
}

type Mount struct {
	Mount string
	Total uint64
	Used  uint64
	Free  uint64
}

type Disk struct {
	Mounts []Mount
}

func (d Disk) MessageType() string {
	return "message.Disk"
}

type Date struct {
	DateTime time.Time
}

func (d Date) MessageType() string {
	return "message.Date"
}

type NetInterface struct {
	Name string
	Rx   uint64
	Tx   uint64
}

type Network struct {
	Interfaces []NetInterface
}

func (n Network) MessageType() string {
	return "message.Network"
}

type Uptime struct {
	Uptime uint64
}

func (u Uptime) MessageType() string {
	return "message.Uptime"
}
