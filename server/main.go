package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// catchError will cause a log and fatal error if one is present
func catchError(err error, message string) {
	fmt.Println("catch", err, message)
	if nil == err {
		return
	}

	log.Fatalln(message, err)
}

func main() {
	path := os.Args[1]

	application := serveRpc(path)
	defer os.Remove(path)

	go watchNetwork()
	loop(application)
}

func loop(s *Smon) {
	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP)

	for {
		select {
		case <-signals:
			fmt.Println("got close")
			return

		case <-time.After(20 * time.Second):
			if false && time.Now().After(s.last.Add(30*time.Second)) {
				fmt.Println("timeout")
				return
			} else {
				fmt.Println("no timeout")
			}
		}
	}
}
