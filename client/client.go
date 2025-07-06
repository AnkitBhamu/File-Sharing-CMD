package client

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/File-share/config"
	"github.com/File-share/models"
)

func ConnectToServer() {
	log.Println("Please provide address the server!")
	var ipadrr string
	_, err := fmt.Scanf("%s\n", &ipadrr)

	if err != nil {
		log.Println("Error in getting port and ip of the server!", err)
		return
	}

	serversocket, err := net.Dial("tcp", ipadrr)

	if err != nil {
		log.Println("Error in connecting to server", err)
		return
	}

	log.Println("Successfully connected to server!!")

	SendFiles(serversocket)

}

func GetFilePaths() []string {
	filepaths := make([]string, 0)
	var filepath string

	for {
		_, err := fmt.Scanf("%s", &filepath)

		// this error will come at the next line
		if err != nil {
			log.Println("Error in reading filepath!!", err)
			break
		}

		filepaths = append(filepaths, filepath)

	}

	return filepaths
}

func SendFiles(socket net.Conn) {
	log.Println("Give file paths separated by space line need to send to server(reciever)")

	for {
		filepaths := GetFilePaths()
		log.Println("Files selected are : ", filepaths)

		for _, filename := range filepaths {
			SendFile(filename, socket)
		}
	}

}

func SendFile(filename string, socket net.Conn) {
	log.Println("Sending the file : ", filename)

	//first read the file
	file, err := os.Open(filename)

	if err != nil {
		log.Println("Error in opening the file!", err)
		return
	}

	//getting the stat of the file
	fileinfo, err := file.Stat()

	if err != nil {
		log.Println("Error in getting the fileinfo!!")
		return
	}

	// first send the metadata of the file to server
	metadata := models.FileMetaData{
		Filename: fileinfo.Name(),
		Size:     fileinfo.Size(),
	}

	marshalled_data, err := json.Marshal(metadata)

	if err != nil {
		log.Println("Error in marshalling the metadata for the file")
		return
	}

	// here -1 we are using because this is never used in json
	marshalled_data = append(marshalled_data, []byte{0xFF}...)

	// send metadata to server
	byteswritten, err := socket.Write(marshalled_data)

	if err != nil || byteswritten != len(marshalled_data) {
		log.Println("Error in sending the metadata to client!")
		return
	}

	// now start sending the content
	// read the file in chunks
	progress := 0.0
	totaldatasent := 0
	filebuffer := make([]byte, config.GetConfig().Sender.FileReadChunksize)
	fmt.Printf("\rSending file : %s : %0.3f%%", fileinfo.Name(), progress)

	for {
		bytesread, err := file.Read(filebuffer)

		if err != nil && err != io.EOF {
			log.Println("Error in reading the file!!")
			return
		}

		// file is completed already
		if err == io.EOF {
			break
		}

		bytessent, err := socket.Write(filebuffer[:bytesread])

		if (bytessent != bytesread) || err != nil {
			log.Println("Error in sending this chunk to server", err)
			return
		}
		totaldatasent += bytessent
		progress = (float64(totaldatasent) / float64(fileinfo.Size())) * 100.0
		fmt.Printf("\rSending file : %s : %0.3f%%", fileinfo.Name(), progress)

	}
	fmt.Println("\nFile : ", fileinfo.Name(), "sent successfully!")

}
