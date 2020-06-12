package gopsql

type Config struct {
	Host           string `json:"host" yaml:"host"`
	Port           int    `json:"port" yaml:"port"`
	UserName       string `json:"userName" yaml:"userName"`
	DbName         string `json:"dbName" yaml:"dbName"`
	Password       string `json:"password" yaml:"password"`
	SslMode        string `json:"sslMode" yaml:"sslMode"`
	Binary         bool   `json:"binary" yaml:"binary"`
	MaxConnections int    `json:"maxConnections" yaml:"maxConnections"`
	ConnectionIdle int    `json:"connectionIdle" yaml:"connectionIdle"`
}
