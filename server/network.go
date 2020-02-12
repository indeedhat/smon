package main

import (
	"io/ioutil"
	"strconv"
	"strings"
	"sync"
	"time"

	. "github.com/indeedhat/smon/message"
)

const (
	c_NetInterval = 3
)

var (
	netIntPrevious    map[string]NetInterface
	networkInterfaces map[string]NetInterface
	netMux            sync.Mutex
)

func init() {
	netIntPrevious = make(map[string]NetInterface)
	networkInterfaces = make(map[string]NetInterface)
}

func watchNetwork() {
	for {
		netMux.Lock()

		func() {
			device, err := ioutil.ReadFile("/proc/net/dev")
			if nil != err {
				return
			}

			for _, line := range strings.Split(string(device), "\n") {
				if !strings.Contains(line, ":") {
					continue
				}

				fields := strings.Fields(line)
				name := strings.Trim(strings.Split(line, ":")[0], " ")

				rx, _ := strconv.ParseUint(fields[1], 10, 64)
				tx, _ := strconv.ParseUint(fields[9], 10, 64)

				if intf, ok := networkInterfaces[name]; ok {
					prev, _ := netIntPrevious[name]

					prev.Rx += intf.Rx
					prev.Tx += intf.Tx
					intf.Rx = 0
					intf.Tx = 0

					if 0 != rx {
						intf.Rx = rx
					}

					if 0 != tx {
						intf.Tx = tx
					}

					netIntPrevious[name] = prev
					networkInterfaces[name] = intf
				} else {
					netIntPrevious[name] = NetInterface{name, rx, tx}
					networkInterfaces[name] = NetInterface{name, 0, 0}
				}
			}

			netMux.Unlock()
		}()

		time.Sleep(c_NetInterval * time.Second)
	}
}
