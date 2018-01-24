package cmd

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/iyacontrol/telegraf-proxy/api"
	"github.com/iyacontrol/telegraf-proxy/config"
	"github.com/iyacontrol/telegraf-proxy/discovery"
)

var (
	fConfig = flag.String("config", "", "configuration file to load")
	etcd    = flag.String("etcd", "http://127.0.0.1:2379", "etcd url")
)

var (
	nextVersion = "1.5.0"
	version     string
	commit      string
	branch      string
	center      *discovery.Center
)

func displayVersion() string {
	if version == "" {
		return fmt.Sprintf("v%s~%s", nextVersion, commit)
	}
	return "v" + version
}

func init() {
	// If commit or branch are not set, make that clear.
	if commit == "" {
		commit = "unknown"
	}
	if branch == "" {
		branch = "unknown"
	}
}

func main() {
	flag.Parse()

	// init config
	config.InitConfig(*fConfig)

	stop := make(chan struct{})
	signals := make(chan os.Signal)

	// init discovery

	endpoints := strings.Split(*etcd, ",")
	center = discovery.NewCenter(endpoints)

	// init api
	api.InitAPI(center)

	// wait for signals to stop or reload
	signal.Notify(signals, os.Interrupt, syscall.SIGHUP)
	go func() {
		select {
		case sig := <-signals:
			if sig == os.Interrupt {
				log.Printf("I! Closing Telegraf-proxy config\n")
				close(stop)
			}
			if sig == syscall.SIGHUP {
				log.Printf("I! Reloading Telegraf-proxy config\n")
			}
		case <-stop:
			return
		}
	}()

	log.Printf("I! Starting Telegraf-proxy %s\n", displayVersion())

	<-stop
	log.Printf("I! Stop Telegraf-proxy %s\n", displayVersion())
}
