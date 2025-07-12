package main

import (
	"fmt"

	"github.com/File-share/client"
	"github.com/File-share/config"
	"github.com/File-share/constants"
	"github.com/File-share/flags"
	"github.com/File-share/server"
)

func main() {
	config.Init()
	flags.Init()
	var mode = flags.Mode()

	if mode == constants.Receiver {
		server.StartServer()

	} else if mode == constants.Sender {
		client.ConnectToServer()
	} else {
		fmt.Println("Invalid mode selected please use sender or receiver")
	}

}
