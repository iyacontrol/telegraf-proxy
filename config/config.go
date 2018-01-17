package config

// Cfg global config
var Cfg Config

// Config defines config for telegraf-proxy
type Config struct {
	Etcd *Etcd
	HTTP *HTTP
}

// Etcd defines etcd config
type Etcd struct {
	Endpoints []string //etcd urls
	Dir       string   // registry path
}

// HTTP http addr and port
type HTTP struct {
	Address string
	Port    string
}
