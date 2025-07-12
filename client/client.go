package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/File-share/constants"
	"github.com/File-share/flags"
	"github.com/File-share/models"
)

func ConnectToServer() {
	ipaddr := flags.ReceiverIP()
	serversocket, err := net.Dial("tcp", ipaddr)

	if err != nil {
		fmt.Println("Error in connecting to server", err)
		return
	}

	fmt.Println("Successfully connected to receiver...")
	SendFiles(serversocket)

}

func GetFilePaths() []string {

	filetosentdir := flags.FilestosentDir()
	filedata, err := os.ReadFile(filetosentdir)

	if err != nil {
		fmt.Println("error in reading the content of filestosent file", err)
		return nil
	}
	// fmt.Println("Filecontent is : ", string(filedata))
	filedatastring := strings.ReplaceAll(string(filedata), "\r\n", "\n")
	filepaths := strings.Split(string(filedatastring), "\n")

	return filepaths
}

func SendFiles(socket net.Conn) {
	fmt.Println("Reading filestosent file....")
	filepaths := GetFilePaths()

	for _, filename := range filepaths {
		if filename == "\n" {
			continue
		}
		SendFile(filename, socket)
	}

}

func SendFile(filename string, socket net.Conn) {
	//first read the file
	file, err := os.Open(filename)

	if err != nil {
		fmt.Println("Error in opening the file!", filename, err)
		return
	}

	//getting the stat of the file
	fileinfo, err := file.Stat()

	if err != nil {
		fmt.Println("Error in getting the fileinfo!!")
		return
	}

	// first send the metadata of the file to server
	metadata := models.FileMetaData{
		Filename: fileinfo.Name(),
		Size:     fileinfo.Size(),
	}

	marshalled_data, err := json.Marshal(metadata)

	if err != nil {
		fmt.Println("Error in marshalling the metadata for the file")
		return
	}

	// here -1 we are using because this is never used in json
	marshalled_data = append(marshalled_data, []byte{0xFF}...)

	// send metadata to server
	byteswritten, err := socket.Write(marshalled_data)

	if err != nil || byteswritten != len(marshalled_data) {
		// fmt.Println("Error in sending the metadata to client!")
		fmt.Printf("\rSending file :%s(failed!) \n", metadata.Filename)
		return
	}

	// now start sending the content
	// read the file in chunks
	var totaldatasent int64
	totaldatasent = 0
	filebuffer := make([]byte, constants.FileReadChunksize)

	// wait till the  speedmeter ends
	var wg sync.WaitGroup
	// wait for one thread to work
	wg.Add(1)

	// just to track the progress
	go SpeedProgressMeter(&wg, &totaldatasent, fileinfo.Size(), fileinfo.Name())

	for {
		bytesread, err := file.Read(filebuffer)
		if err != nil && err != io.EOF {
			// fmt.Println("Error in reading the file!!")
			totaldatasent = -1
			return
		}

		// file is completed already
		if err == io.EOF {
			break
		}

		bytessent, err := socket.Write(filebuffer[:bytesread])

		if (bytessent != bytesread) || err != nil {
			// fmt.Println("Error in sending this chunk to server", err)
			totaldatasent = -1
		}
		totaldatasent += int64(bytessent)

	}

	wg.Wait()

}

func SpeedProgressMeter(wg *sync.WaitGroup, totaldatasent *int64, totalbytesneedtosend int64, filename string) {
	// iterate every 1 second till all databeing sent
	var totalsent int64
	totalsent = 0
	for {
		time.Sleep(1 * time.Second)

		//means send failed
		if *totaldatasent == -1 {
			fmt.Printf("\033[2K\rSending file :%s(failed!) \n", filename)
			wg.Done()
			return
		}

		speed_divider := 1024 * 1024
		speed := float64(*totaldatasent-totalsent) / float64(speed_divider)
		totalsent = *totaldatasent
		progress := (float64(*totaldatasent) / float64(totalbytesneedtosend)) * 100.0

		if *totaldatasent >= totalbytesneedtosend {
			fmt.Printf("\033[2K\rSending file:%s(done)\n ", filename)
			wg.Done()
			return
		}

		fmt.Printf("\033[2K\rSending file:%s(%0.3f%%) speed: %0.3f MB/s ", filename, progress, speed)

	}
}
