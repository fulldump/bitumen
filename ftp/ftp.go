package ftp

import (
	"github.com/koofr/graval"
	"log"
)

type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	BasePath string
}

func ListenAndServe(c *Config) error {

	factory := &LocalDriverFactory{c.Username, c.Password, c.BasePath}

	server := graval.NewFTPServer(&graval.FTPServerOpts{
		ServerName: "Example FTP server",
		Factory:    factory,
		Hostname:   c.Host,
		Port:       c.Port,
		PassiveOpts: &graval.PassiveOpts{
			ListenAddress: c.Host,
			NatAddress:    c.Host,
			PassivePorts: &graval.PassivePorts{
				Low:  42000,
				High: 45000,
			},
		},
	})

	log.Printf("Example FTP server listening on %s:%d", c.Host, c.Port)
	log.Printf("Access: ftp://%s:%s@%s:%d/", c.Username, c.Password, c.Host, c.Port)

	return server.ListenAndServe()
}
