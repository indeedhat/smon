package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"syscall"
	"time"

	. "github.com/indeedhat/smon/message"
)

// rpc host struct thing
type Smon struct {
	last time.Time
}

// our remotely invocable function
func (s *Smon) Ping(req RpcRequest, res *RpcResponse) (err error) {
	s.begin()

	fmt.Println("ping")
	fmt.Println(req)
	res.Message = PingPong{
		Message: "pong",
	}

	return
}

// Cpu returns the total and idle usage on the cpu as a whole
func (s *Smon) Cpu(req RpcRequest, res *RpcResponse) error {
	s.begin()

	contents, err := ioutil.ReadFile("/proc/stat")
	if nil != err {
		return err
	}

	response := Cpu{}
	lines := strings.Split(string(contents), "\n")

	for _, line := range lines {
		fields := strings.Fields(line)

		if "cpu" != fields[0] {
			continue
		}

		for i, field := range fields {
			val, _ := strconv.ParseUint(field, 10, 64)

			response.Total += val

			if 4 == i {
				response.Idle = val
			}
		}

		break
	}

	res.Message = response

	return nil
}

// Memory returns the stats of used, free and total memory
func (s *Smon) Memory(req RpcRequest, res *RpcResponse) error {
	s.begin()

	contents, err := ioutil.ReadFile("/proc/meminfo")
	if nil != err {
		return err
	}

	response := Memory{}
	lines := strings.Split(string(contents), "\n")

	for _, line := range lines {
		if strings.HasPrefix(line, "MemTotal") {
			fields := strings.Fields(line)
			response.Total, _ = strconv.ParseUint(fields[1], 10, 64)
			response.Total *= 1024
		} else if strings.HasPrefix(line, "MemAvailable") {
			fields := strings.Fields(line)
			response.Free, _ = strconv.ParseUint(fields[1], 10, 64)
			response.Free *= 1024
		}
	}

	response.Used = response.Total - response.Free
	res.Message = response

	return nil
}

// Date returns the date and time on the machine
func (s *Smon) Date(req RpcRequest, res *RpcResponse) error {
	s.begin()

	res.Message = Date{
		DateTime: time.Now(),
	}

	return nil
}

// Uptime returns the number of secods the machine has been up
func (s *Smon) Uptime(req RpcRequest, res *RpcResponse) error {
	s.begin()

	line, err := ioutil.ReadFile("/proc/uptime")
	if nil != err {
		return err
	}

	parts := strings.Split(string(line), ".")

	response := Uptime{}
	response.Uptime, err = strconv.ParseUint(parts[0], 10, 64)
	if nil != err {
		return err
	}

	res.Message = response
	return nil
}

// Disk returns the usage statistics for the given mount points
func (s *Smon) Disk(req RpcRequest, res *RpcResponse) error {
	s.begin()

	request, ok := req.Message.(Disk)
	if !ok {
		return fmt.Errorf("invalid request")
	}

	if 1 > len(request.Mounts) {
		return fmt.Errorf("no mount points given")
	}

	disk := Disk{}
	for _, mount := range request.Mounts {
		fs := syscall.Statfs_t{}
		err := syscall.Statfs(mount.Mount, &fs)
		if nil != err {
			return err
		}

		mnt := Mount{
			Mount: mount.Mount,
			Total: fs.Blocks * uint64(fs.Bsize),
			Free:  fs.Bfree * uint64(fs.Bsize),
		}

		mnt.Used = mnt.Total - mnt.Free
		disk.Mounts = append(disk.Mounts, mnt)
	}

	res.Message = disk

	return nil
}

// Network returns bandwidth stats for the machine
func (s *Smon) Network(req RpcRequest, res *RpcResponse) error {
	s.begin()

	netMux.Lock()
	network := Network{}

	for _, intf := range networkInterfaces {
		network.Interfaces = append(network.Interfaces, NetInterface{
			Name: intf.Name,
			Rx:   intf.Rx / c_NetInterval,
			Tx:   intf.Tx / c_NetInterval,
		})
	}

	netMux.Unlock()

	res.Message = network
	return nil
}

// begin a new command
func (s *Smon) begin() {
	s.last = time.Now()
}
