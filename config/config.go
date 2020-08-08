package config

import (
	"bitumen/ftp"
)

type Config struct {
	Http Http

	BasePath string `usage:"Path to directory to be served"`

	Ftp ftp.Config
}

type Http struct {
	Addr string `usage:"Server address to listen from"`
}
