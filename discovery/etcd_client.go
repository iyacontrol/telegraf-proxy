package discovery

import (
	"log"
	"time"

	"github.com/coreos/go-etcd/etcd"
)

type WrappedEtcdClient struct {
	etcdClient *etcd.Client
}

func (w *WrappedEtcdClient) WatchEtcd(dir string, ch chan *etcd.Response, stop chan bool, handle func(*etcd.Response)) {

	watcher := func() {
		for {
			watchCh := make(chan *etcd.Response, 10)
			var err error
			go func() {
				for {
					select {
					case rs := <-watchCh:
						if rs != nil {
							ch <- rs
						} else {
							//receives nil when etcd is not reachable
							time.Sleep(5 * time.Second)
						}
					case <-stop:
						break
					default:
						time.Sleep(time.Second)
					}
				}
			}()

			_, err = w.etcdClient.Watch(dir, 0, true, watchCh, stop)

			if err != nil {
				log.Printf("Error watching %s, %s, retry in 10s.", dir, err.Error())
				time.Sleep(10 * time.Second)
			}
		}
	}

	receiver := func() {
		for {
			res := <-ch
			if res != nil {
				handle(res)
			}
		}
		stop <- true
	}

	log.Printf("Watching %s.", dir)

	go watcher()
	go receiver()
}
