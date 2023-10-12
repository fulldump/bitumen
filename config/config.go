package config

import (
	"bitumen/ftp"
	"bitumen/sftpinterface"
)

type Config struct {
	Http Http

	BasePath string `usage:"Path to directory to be served"`

	Ftp ftp.Config

	Sftp sftpinterface.Config
}

type Http struct {
	Addr string `usage:"Server address to listen from"`
}
