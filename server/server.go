package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/File-share/config"
	"github.com/File-share/models"
)

const (
	filecontentreadingstate = "reading-file"
	readingmetadatastate    = "metadata-reading"
)

func StartServer() {
	server, err := net.Listen("tcp", ":8080")

	if err != nil {
		log.Fatal("Error in starting the server !!", err)
	}

	// start listening for the clients
	for {
		clientconn, err := server.Accept()

		if err != nil {
			log.Println("Error in accepting connection from the client!!", err)
		}
		log.Println("Accepted client connection with ip addr", clientconn.LocalAddr().String())
		go ReceiveFiles(clientconn)
	}
}

// so this will accept the connections from client
func ReceiveFiles(socket net.Conn) {
	readbuffer := make([]byte, config.GetConfig().Receiver.Tcpreadbuffersize)
	var fileptr *os.File
	var currentstate = readingmetadatastate
	var metadata []byte
	var contentcounter int64
	rcvprogress := 0.0

	var metadatastruct models.FileMetaData

	for {
		bytesread, err := socket.Read(readbuffer)
		if err != nil {
			log.Println("Error in reading the content from client", err)
			return
		}

		// try to consume entire buffer read now
		for i := 0; i < bytesread; i++ {
			if currentstate == readingmetadatastate {
				if int(readbuffer[i]) != 255 {
					metadata = append(metadata, readbuffer[i])

				} else {
					log.Println("Metadata : ", string(metadata))
					fileptr, metadatastruct, _ = HandleMetaDataOps(metadata)
					currentstate = filecontentreadingstate
					metadata = []byte{}
					continue

				}
			}

			if currentstate == filecontentreadingstate {
				// Calculate how many bytes remain in this chunk
				remainingBytes := metadatastruct.Size - contentcounter
				availableBytes := int64(bytesread - i)

				// Determine how many bytes to write in this round
				toWrite := remainingBytes
				if availableBytes < remainingBytes {
					toWrite = availableBytes
				}

				// Write the chunk
				fileptr.Write(readbuffer[i : i+int(toWrite)])
				contentcounter += toWrite
				i += int(toWrite) - 1 // -1 because the loop will increment i again
				rcvprogress = (float64(contentcounter) / float64(metadatastruct.Size)) * 100.0
				fmt.Printf("\rReceiving file : %s : %0.3f%%", metadatastruct.Filename, rcvprogress)

				if contentcounter >= metadatastruct.Size {
					contentcounter = 0
					fmt.Printf("\nReceived file : %s successfully\n", metadatastruct.Filename)
					currentstate = readingmetadatastate
					metadatastruct = models.FileMetaData{}
					fileptr.Close()
					fileptr = nil
				}
			}

		}

	}

}

func HandleMetaDataOps(data []byte) (file *os.File, metadata models.FileMetaData, err error) {

	err = json.Unmarshal(data, &metadata)
	if err != nil {
		log.Println("Error in parsing the metadata!", err)
		return
	}

	//directory should be already there
	filepath := config.GetConfig().Receiver.DownloadDirectory + metadata.Filename
	file, err = os.Create(filepath)

	if err != nil {
		log.Println("Error in creating file from metadata", err)
		return
	}

	return

}
