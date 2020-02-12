package main

import (
	"bytes"
	"fmt"
	"time"

	"github.com/indeedhat/smon/message"
)

type ClientData struct {
	Name    string
	Cpu     message.Cpu
	Memory  message.Memory
	Disk    message.Disk
	Date    message.Date
	Network message.Network
	Uptime  message.Uptime
}

func (c *ClientData) String() string {
	var buffer bytes.Buffer

	fmt.Fprintf(&buffer, "Cpu:\n  %s\n\n", formatPercent(float64(c.Cpu.Total), float64(c.Cpu.Total-c.Cpu.Idle), 2))

	// mod := uint64(1024 * 1024 * 1024)
	fmt.Fprintf(
		&buffer,
		"Memory:\n  Total: %s\n  Used: %s\n  Free: %s\n\n",
		formatBytes(c.Memory.Total, 2),
		formatBytes(c.Memory.Used, 2),
		formatBytes(c.Memory.Free, 2),
	)

	fmt.Fprintf(&buffer, "Date:\n  %s\n\n", c.Date.DateTime.Format("02/01/2006 15:04:05"))

	fmt.Fprintf(&buffer, "Uptime:\n %s\n\n", formatDuration(time.Duration(c.Uptime.Uptime)*time.Second))

	fmt.Fprintln(&buffer, "Disk:")

	for _, disk := range c.Disk.Mounts {
		fmt.Fprintf(
			&buffer,
			"  %s\n    Total: %s\n    Used: %s\n    Free: %s\n",
			disk.Mount,
			formatBytes(disk.Total, 2),
			formatBytes(disk.Used, 2),
			formatBytes(disk.Free, 2),
		)
	}

	fmt.Fprintln(&buffer, "")

	fmt.Fprintln(&buffer, "Network:")

	for _, iface := range c.Network.Interfaces {
		fmt.Fprintf(
			&buffer,
			"  %s\n    Rx: %s\n    Tx: %s\n",
			iface.Name,
			formatBits(iface.Rx, 2),
			formatBits(iface.Tx, 2),
		)
	}

	return buffer.String()
}

var clients []*ClientData

func NewClient(name string) *ClientData {
	return &ClientData{
		Name: name,
	}
}
