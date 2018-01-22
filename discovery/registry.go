package discovery

import (
	"log"
	"sync"

	"github.com/coreos/go-etcd/etcd"
)

type Registry struct {
	sync.RWMutex
	data map[string]string
}

func (r *Registry) Register(node *etcd.Node) {
	r.Lock()
	log.Printf("Register (%s,%s).", node.Key, node.Value)
	r.data[node.Key] = node.Value
	r.Unlock()
}

func (r *Registry) Unregister(node *etcd.Node) {
	r.Lock()
	if _, found := r.data[node.Key]; found {
		log.Printf("Unregister (%s,%s).", node.Key, node.Value)
		delete(r.data, node.Key)
	}
	r.Unlock()
}

func (r *Registry) Translate(path string) []string {
	var targets []string
	r.Lock()
	for _, v := range r.data {
		targets = append(targets, v+"/"+path)
	}

	r.Unlock()
	return targets
}
