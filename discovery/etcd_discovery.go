package discovery

import (
	"encoding/json"
	"log"
	"time"

	"github.com/coreos/etcd/client"
	"github.com/iyacontrol/telegraf-proxy/config"
	"golang.org/x/net/context"
)

// Center registry center
type Center struct {
	members map[string]*Member
	KeysAPI client.KeysAPI
}

// Member is a client machine
type Member struct {
	InGroup bool
	IP      string
	Name    string
}

// WorkerInfo is the service register information to etcd
type WorkerInfo struct {
	Name string
	IP   string
}

// NewCenter ...
func NewCenter(endpoints []string) *Center {
	cfg := client.Config{
		Endpoints:               endpoints,
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}

	etcdClient, err := client.New(cfg)
	if err != nil {
		log.Fatal("Error: cannot connec to etcd:", err)
	}

	master := &Center{
		members: make(map[string]*Member),
		KeysAPI: client.NewKeysAPI(etcdClient),
	}
	go master.WatchWorkers()
	return master
}

func (m *Center) AddWorker(info *WorkerInfo) {
	member := &Member{
		InGroup: true,
		IP:      info.IP,
		Name:    info.Name,
	}
	m.members[member.Name] = member
}

func (m *Center) UpdateWorker(info *WorkerInfo) {
	member := m.members[info.Name]
	member.InGroup = true
}

func NodeToWorkerInfo(node *client.Node) *WorkerInfo {
	log.Println(node.Value)
	info := &WorkerInfo{}
	err := json.Unmarshal([]byte(node.Value), info)
	if err != nil {
		log.Print(err)
	}
	return info
}

func (m *Center) WatchWorkers() {
	api := m.KeysAPI
	watcher := api.Watcher(config.Cfg.Etcd.Dir, &client.WatcherOptions{
		Recursive: true,
	})
	for {
		res, err := watcher.Next(context.Background())
		if err != nil {
			log.Println("Error watch workers:", err)
			break
		}
		if res.Action == "expire" {
			info := NodeToWorkerInfo(res.PrevNode)
			log.Println("Expire worker ", info.Name)
			member, ok := m.members[info.Name]
			if ok {
				member.InGroup = false
			}
		} else if res.Action == "set" {
			info := NodeToWorkerInfo(res.Node)
			if _, ok := m.members[info.Name]; ok {
				log.Println("Update worker ", info.Name)
				m.UpdateWorker(info)
			} else {
				log.Println("Add worker ", info.Name)
				m.AddWorker(info)
			}
		} else if res.Action == "delete" {
			info := NodeToWorkerInfo(res.Node)
			log.Println("Delete worker ", info.Name)
			delete(m.members, info.Name)
		}
	}
}

func (m *Center) Translate(path string) []string {
	var urls []string
	for _, v := range m.members {
		urls = append(urls, "http://"+v.IP+path)
	}

	return urls
}
