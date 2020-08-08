package ftp

import (
	"github.com/koofr/graval"
)

type LocalDriverFactory struct {
	Username string
	Password string
	BasePath string
}

func (f *LocalDriverFactory) NewDriver() (d graval.FTPDriver, err error) {
	return &LocalDriver{
		Username: f.Username,
		Password: f.Password,
		BasePath: f.BasePath,
	}, nil
}
