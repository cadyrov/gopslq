package gopsql

type Config struct {
	Host           string `json:"host"`
	Port           int    `json:"port"`
	UserName       string `json:"userName"`
	DbName         string `json:"dbName"`
	Password       string `json:"password"`
	SslMode        string `json:"sslMode"`
	Binary         bool   `json:"binary"`
	MaxConnections int    `json:"maxConnections"`
	ConnectionIdle int    `json:"connectionIdle"`
}
