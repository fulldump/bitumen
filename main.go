package main

import (
	"fmt"
	"net/http"

	"github.com/fulldump/goconfig"

	"bitumen/config"
)

func main() {
	fmt.Println("Bitumen Server")

	c := &config.Config{
		BasePath: "/tmp",
		Http: config.Http{
			Addr: ":8080",
		},
	}

	goconfig.Read(&c)

	fs := http.FileServer(http.Dir(c.BasePath))

	s := &http.Server{
		Addr:    c.Http.Addr,
		Handler: fs,
	}

	fmt.Println("Serving on ", s.Addr)
	if err := s.ListenAndServe(); err != nil {
		fmt.Println("ERROR:", err)
	}

}
