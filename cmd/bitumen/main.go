package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/fulldump/goconfig"

	"bitumen/config"
	"bitumen/fileserver"
	"bitumen/ftp"
	"bitumen/sftpinterface"
)

func main() {
	fmt.Println("Bitumen Server")

	c := &config.Config{
		BasePath: "/tmp",
		Http: config.Http{
			Addr: ":8080",
		},
		Ftp: ftp.Config{
			Host: "0.0.0.0",
			Port: 2121,
			Credentials: []ftp.Credential{
				{
					Username: "alice",
					Password: "111",
					BasePath: "/tmp/alice",
					ReadOnly: false,
				},
				{
					Username: "bob",
					Password: "222",
					BasePath: "/tmp/bob",
					ReadOnly: false,
				},
			},
		},
		Sftp: sftpinterface.Config{
			Addr:     ":2222",
			User:     "user",
			Password: "pass",
		},
	}

	goconfig.Read(&c)

	// *
	// FTP server
	go func() {
		// TODO: handle panic
		for {

			fmt.Println("Starting ftp")
			err := ftp.Serve(&c.Ftp)
			if err == nil {
				return
			}
			fmt.Println("FTP error:", err)

			nextRetry := 10 * time.Second
			fmt.Printf("FTP next retry in %s\n", nextRetry)
			time.Sleep(nextRetry)
		}
	}()

	// SFTP server
	go func() {
		// TODO: handle panic
		for {
			sftp_server, err := sftpinterface.NewSftpInterface(c.Sftp)
			if err != nil {
				fmt.Println("new sftp:", err.Error())
			}
			err = sftp_server.Serve()
			if err != nil {
				log.Panicln("serve sftp:", err.Error())
				time.Sleep(10 * time.Second)
			}
		}
	}()
	// */

	// HTTP Server
	fs := fileserver.FileServer(fileserver.Dir(c.BasePath))
	http_server := &http.Server{
		Addr:    c.Http.Addr,
		Handler: fs,
	}
	fmt.Println("Serving on ", http_server.Addr)
	if err := http_server.ListenAndServe(); err != nil {
		fmt.Println("ERROR:", err)
	}

}
