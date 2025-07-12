package flags

import (
	"flag"

	"github.com/File-share/constants"
)

var mode *string
var receiverIp *string
var filestosentdir *string
var port *string
var downloaddir *string

func Init() {
	mode = flag.String("mode","", "flag to set the mode")
	receiverIp = flag.String("rcvIp", "", "flag to set receiver ip")
	filestosentdir = flag.String("sfdr", "", "flag to set filestosentdir")
	port = flag.String("port", "8080", "flag to set the mode")
	downloaddir = flag.String("downloadDir", "", "flag to set download folder path")
	flag.Parse()
}

func Mode() string {
	return *mode
}

func ReceiverIP() string {
	return *receiverIp
}
func FilestosentDir() string {
	return *filestosentdir
}
func Port() string {
	return *port
}

func DownloadDirectory() string {
	return *downloaddir
}
