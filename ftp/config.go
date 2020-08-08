package ftp

type Config struct {
	Host        string
	Port        int
	Credentials []Credential
}

type Credential struct {
	Username string
	Password string
	BasePath string
	ReadOnly bool
}
