package ftp

import (
	"fmt"
	"github.com/koofr/graval"
)

type AuthHook func(username, password string, driver *LocalDriver) bool

type LocalDriverFactory struct {
	BasePath string
	ReadOnly bool
	AuthHook AuthHook
}

func (f *LocalDriverFactory) NewDriver() (d graval.FTPDriver, err error) {

	driver := &LocalDriver{
		BasePath: f.BasePath,
		ReadOnly: f.ReadOnly,
		AuthHook: f.AuthHook,
	}

	if driver.AuthHook == nil {
		fmt.Println("WARNING: Default credentials")
		driver.AuthHook = DefaultCredentials
	}

	return driver, nil
}

func DefaultCredentials(username, password string, driver *LocalDriver) bool {
	return username == "admin" && password == "admin"
}
