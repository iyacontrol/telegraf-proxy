package cmd

import (
	"flag"
	"log"

	"github.com/iyacontrol/telegraf-proxy/config"
	"github.com/iyacontrol/telegraf-proxy/discovery"
)

var fConfig = flag.String("config", "", "configuration file to load")

var (
	registry *discovery.Registry
)

func main() {
	flag.Parse()

	// init config
	config.InitConfig(*fConfig)

	stop := make(chan bool, 1)

	// init discovery
	discovery.InitDiscovery(registry, stop)

	<-stop
	log.Println("Stopped.")

}

func init() {
	registry := &discovery.Registry{
		data: make(map[string]string),
	}
}
