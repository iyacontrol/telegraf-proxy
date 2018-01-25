package config

// Cfg global config
var Cfg Config

// Config defines config for telegraf-proxy
type Config struct {
	Etcd      *Etcd
	HTTP      *HTTP
	Aggregate *Aggregate
	Center    *SettingCenter
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

// Aggregate metrcis aggregate
type Aggregate struct {
	Timeout int
}

// SettingCenter  config center
type SettingCenter struct {
	URL string
}
