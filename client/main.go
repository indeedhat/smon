package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/crypto/ssh"
)

func main() {
	log.Println("booting")

	for i, c := range config.Clients {
		log.Println(c.Name)
		log.Println(c)

		sshConfig := &ssh.ClientConfig{
			User: c.Server.User,
			Auth: []ssh.AuthMethod{
				ssh.Password(c.Server.SSH.Password),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}
		log.Println("created config")

		sshClient, err := ssh.Dial(
			"tcp",
			fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port),
			sshConfig,
		)

		if err != nil {
			log.Fatalf("Failed to dial: %s", err)
		}
		defer sshClient.Close()
		log.Print("ssh connected")

		s, err := sshClient.Dial("unix", c.Server.Socket)
		if err != nil {
			log.Fatalf("unable to start netcat session: %s", err)
		}
		log.Println("socker connected")

		rpc := NewRpcClient(s, c, i)
		defer rpc.close()

		clients = append(clients, NewClient(c.Name))
		go rpc.loop()
		log.Println("started loop")
	}

	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP)

	<-signals
}
