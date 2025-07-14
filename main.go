package main

import (
	"fmt"
	"os"

	"github.com/File-share/client"
	"github.com/File-share/constants"
	"github.com/File-share/flags"
	"github.com/File-share/server"
)

func main() {
	if len(os.Args) >= 2 && os.Args[1] == constants.Help {
		fmt.Print(constants.HelpString)
		return
	}

	flags.Init()
	var mode = flags.Mode()

	switch mode {
	case constants.Sender:
		client.ConnectToServer()

	case constants.Receiver:
		server.StartServer()

	default:
		fmt.Println("Invalid mode selected please use sender or receiver")
	}

}
