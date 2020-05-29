package config

type Config struct {
	Http Http

	BasePath string `usage:"Path to directory to be served"`
}

type Http struct {
	Addr string `usage:"Server address to listen from"`
}
