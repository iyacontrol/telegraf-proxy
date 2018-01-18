package discovery

import (
	"log"
	"time"

	"github.com/coreos/go-etcd/etcd"
	"github.com/iyacontrol/telegraf-proxy/config"
)

func InitDiscovery(stop chan struct{}, reg *Registry) {

	wec := &WrappedEtcdClient{
		etcdClient: etcd.NewClient(config.Cfg.Etcd.Endpoints),
	}
	dir := "/telegraf"
	if config.Cfg.Etcd.Dir != "" {
		dir = config.Cfg.Etcd.Dir
	}

	attempts := 3
	for attempts > 0 {
		resp, err := wec.etcdClient.Get(dir, false, true)
		if err == nil {
			reg.Register(resp.Node)
			break
		} else {
			log.Printf("Error getting %s, %s, retry in 1s.", dir, err.Error())
			attempts--
			time.Sleep(1000 * time.Millisecond)
		}
	}

	if attempts == 0 {
		log.Fatal("Exiting")
	}

	ch := make(chan *etcd.Response, 10)

	handle := func(resp *etcd.Response) {
		setOp := []string{"create", "set", "update", "compareAndSwap"}
		deleteOp := []string{"delete", "expire", "compareAndDelete"}

		for _, op := range setOp {
			if op == resp.Action {
				if resp.PrevNode != nil {
					reg.Unregister(resp.PrevNode)
				}
				reg.Register(resp.Node)
			}
		}

		for _, op := range deleteOp {
			if op == resp.Action {
				reg.Unregister(resp.PrevNode)
			}
		}
	}

	wec.WatchEtcd(dir, ch, stop, handle)

}
