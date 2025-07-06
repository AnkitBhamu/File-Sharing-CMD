package models

type Config struct {
	Sender   SenderConfig   `yaml:"sender"`
	Receiver ReceiverConfig `yaml:"receiver"`
}

type SenderConfig struct {
	FileReadChunksize int64 `yaml:"filereadchunksize"`
}

type ReceiverConfig struct {
	DownloadDirectory string `yaml:"downloadDir"`
	Tcpreadbuffersize int64  `yaml:"tcpreadbuffersize"`
}
