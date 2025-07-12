package constants

const (
	Sender            = "sender"
	Receiver          = "receiver"
	FileReadChunksize = 1048576
	Tcpreadbuffersize = 1638400

	ConfigFile = "./config.yml"
	Help       = "-help"
	HelpString = `Usage: fsgoclient -mode=<sender|receiver> [options]

Modes:
  -mode=sender           Start the tool in sender mode
  -mode=receiver         Start the tool in receiver mode

Sender Options:
  -rcvIp=<ip:port>       IP address and port of the receiver (e.g., localhost:8080)
  -sfdr=<path>           Path to a file containing files path that needs to send (one path per line)

Receiver Options:
  -downloadDir=<path>    Directory where received files will be saved
  -port <port>           Port on which receiver needs to listen (default:8080)

Examples:
  Sender:
    fsgoclient  -mode=sender -rcvIp=localhost:8080 -sfdr="C:\path\to\filestosend.txt"

  Receiver:
    fsgoclient  -port=8080  -mode=receiver -downloadDir="C:\path\to\downloadDir"

Notes:
  - Ensure both sender and receiver are on the same network or accessible via the given IP.
  - In sender mode, the file provided with -sfdr should contain valid paths to the files you want to send.
  - Use absolute paths to avoid file resolution issues.
`
)
