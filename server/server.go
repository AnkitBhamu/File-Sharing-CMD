package server

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"github.com/File-share/config"
	"github.com/File-share/flags"
	"github.com/File-share/models"
)

const (
	filecontentreadingstate = "reading-file"
	readingmetadatastate    = "metadata-reading"
)

func StartServer() {
	server, err := net.Listen("tcp", ":"+flags.Port())

	fmt.Println("Wating for senders on ", flags.Port(), " port...")

	if err != nil {
		fmt.Println("Error in starting the server !!", err)
	}

	// start listening for the clients
	for {
		clientconn, err := server.Accept()

		if err != nil {
			fmt.Println("Error in accepting connection from the sender...", err)
		}
		fmt.Println("Connected to sender...")
		go ReceiveFiles(clientconn)
	}
}

// so this will accept the connections from client
func ReceiveFiles(socket net.Conn) {
	readbuffer := make([]byte, config.GetConfig().Receiver.Tcpreadbuffersize)
	var fileptr *os.File
	var currentstate = readingmetadatastate
	var metadata []byte
	// var contentcounter int64

	var metadatastruct models.FileMetaData

	for {
		bytesread, err := socket.Read(readbuffer)
		if err != nil {
			// fmt.Println("Error in reading the content from client", err)
			continue
		}

		// try to consume entire buffer read now
		for i := 0; i < bytesread; i++ {
			if currentstate == readingmetadatastate {
				if int(readbuffer[i]) != 255 {
					metadata = append(metadata, readbuffer[i])

				} else {
					fileptr, metadatastruct, _ = HandleMetaDataOps(metadata)
					index, maxbyteread, err := HandleFileDownload(socket, metadatastruct, fileptr, readbuffer, i+1, bytesread)

					if err != nil {
						// fmt.Println("Error in downloading the file", err)
						break
					}
					metadata = []byte{}
					i = index
					bytesread = maxbyteread

				}
			}

		}

	}

}

func HandleFileDownload(socket net.Conn, metadatastruct models.FileMetaData, fileptr *os.File, readbuffer []byte, start_index, bytesread int) (int, int, error) {

	//first read the remaining buffer first
	var contentcounter int64
	var wg sync.WaitGroup
	var downloadcompleted bool
	var maxbyteread int
	var index int
	contentcounter = 0
	downloadcompleted = false
	maxbyteread = bytesread
	wg.Add(1)

	// start the speed meter
	go SpeedProgressMeter(&wg, &contentcounter, metadatastruct.Size, metadatastruct.Filename)

	for i := start_index; i < bytesread; i++ {
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

		if contentcounter >= metadatastruct.Size {
			fileptr.Close()
			downloadcompleted = true
			index = i
			break
		}
	}

	// again reading the buffer of data from client
	for {
		bytesread, err := socket.Read(readbuffer)
		if err != nil {
			// fmt.Println("Error in reading the content from client", err)
			// delete the file also from os
			fileptr.Close()
			os.Remove(flags.DownloadDirectory() + "/" + metadatastruct.Filename)
			contentcounter = -1
			return -1, -1, err
		}
		maxbyteread = bytesread

		for i := 0; i < bytesread; i++ {
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

			if contentcounter >= metadatastruct.Size {
				fileptr.Close()
				downloadcompleted = true
				index = i
				break
			}
		}

		if downloadcompleted {
			break
		}
	}

	wg.Wait()
	return index, maxbyteread, nil

}

func HandleMetaDataOps(data []byte) (file *os.File, metadata models.FileMetaData, err error) {

	err = json.Unmarshal(data, &metadata)
	if err != nil {
		fmt.Println("Error in parsing the metadata!", err)
		return
	}

	//directory should be already there
	filepath := flags.DownloadDirectory() + "/" + metadata.Filename
	file, err = os.Create(filepath)

	if err != nil {
		fmt.Println("Error in creating file from metadata", err)
		return
	}

	return

}

func SpeedProgressMeter(wg *sync.WaitGroup, totaldatarecvd *int64, totalbytestorecv int64, filename string) {
	// iterate every 1 second till all databeing sent
	var totalrecved int64
	totalrecved = 0

	for {

		time.Sleep(1 * time.Second)
		// means download is failed
		if *totaldatarecvd == -1 {
			fmt.Printf("\033[2K\rReceiving file :%s(failed!) ", filename)
			wg.Done()
			return
		}

		if *totaldatarecvd >= totalbytestorecv {
			fmt.Printf("\033[2K\r%-80s\rReceiving file :%s(done)\n", "", filename)
			wg.Done()
			return
		}

		speed_divider := 1024 * 1024
		speed := float64(*totaldatarecvd-totalrecved) / float64(speed_divider)
		totalrecved = *totaldatarecvd
		progress := (float64(*totaldatarecvd) / float64(totalbytestorecv)) * 100.0

		fmt.Printf("\033[2K\rReceiving file:%s(%0.3f%%) speed: %0.3f MB/s ", filename, progress, speed)

	}
}
