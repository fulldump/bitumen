package ftp

import (
	"fmt"
	"log"

	"github.com/koofr/graval"
)

func Serve(c *Config) error {

	factory := &LocalDriverFactory{}

	factory.AuthHook = func(username, password string, driver *LocalDriver) bool {

		fmt.Println("DEBUG user:", username)

		for _, credential := range c.Credentials {
			if credential.Username != username {
				fmt.Println("DEBUG Skipping user", credential.Username)
				continue // skip user
			}

			if credential.Password != password {
				fmt.Println("DEBUG Login attemp error for user", credential.Username)
				return false
			}

			fmt.Println("DEBUG Success login")
			driver.BasePath = credential.BasePath
			driver.ReadOnly = credential.ReadOnly
			return true
		}

		fmt.Println("DEBUG User does not exist")
		return false
	}

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

	return server.ListenAndServe()
}
