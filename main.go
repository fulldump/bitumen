package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/fulldump/goconfig"

	"bitumen/config"
	"bitumen/fileserver"
	"bitumen/ftp"
)

func main() {
	fmt.Println("Bitumen Server")

	c := &config.Config{
		BasePath: "/tmp",
		Http: config.Http{
			Addr: ":8080",
		},
		Ftp: ftp.Config{
			Host:     "0.0.0.0",
			Port:     2121,
			Username: "test",
			Password: "test",
			BasePath: "/tmp/scanner",
		},
	}

	goconfig.Read(&c)

	fs := fileserver.FileServer(fileserver.Dir(c.BasePath))

	s := &http.Server{
		Addr:    c.Http.Addr,
		Handler: fs,
	}

	//fmt.Println("Starting sftp")
	//go startSftp()

	go func() {
		// TODO: Recover handler
		for {

			fmt.Println("Starting ftp")
			err := ftp.ListenAndServe(&c.Ftp)
			if err == nil {
				return
			}
			fmt.Println("FTP error:", err)

			nextRetry := 10 * time.Second
			fmt.Printf("FTP next retry in %s\n", nextRetry)
			time.Sleep(nextRetry)
		}
	}()

	fmt.Println("Serving on ", s.Addr)
	if err := s.ListenAndServe(); err != nil {
		fmt.Println("ERROR:", err)
	}

}
