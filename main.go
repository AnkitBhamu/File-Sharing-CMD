package main

import (
	"fmt"
	"log"

	"github.com/File-share/client"
	"github.com/File-share/config"
	"github.com/File-share/server"
)

func main() {
	config.Init()
	fmt.Println("Select the mode:\n1.Receiver\n2.Sender")
	var mode int
	fmt.Scanf("%d\n", &mode)

	if mode == 1 {
		log.Println("We are receiver!!")
		server.StartServer()

	}

	if mode == 2 {
		log.Println("We are sender!!")
		client.ConnectToServer()
	}
}
