package etcd

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"log"
)

type config struct {
	server
	database
}

type server struct {
	Addr  string `json:"addr"`
	Port  string `json:"port"`
	Name  string `json:"app"`
	Token string `json:"token"`
}

type database struct {
	Type     string `json:"type"`
	Url      string `json:"url"`
	User     string `json:"user"`
	Password string `json:"password"`
}

var (
	// /config/appName/profile
	configKeyFmt = "/config/%s/%s"
	Config       config
	Env          = flag.String("profile", "test", "use profile for config")
)

func configPath(appName string) string {
	return fmt.Sprintf(configKeyFmt, appName, *Env)
}

func InitConfig(appName string) {
	path := configPath(appName)
	resp, err := client.Get(context.Background(), path)
	if err != nil {
		log.Fatalf("get config from etcd fail, err: %v", err)
	}

	if err := json.Unmarshal(resp.Kvs[0].Value, &Config); err != nil {
		log.Fatalf("unmarshal config info fail, err: %v", err)
	}

	log.Println("init config info success...")

	go watchConfig(path)
}

func watchConfig(path string) {
	watch := client.Watch(context.Background(), path)
	for w := range watch {
		log.Printf("w: %v\n", w)
		for _, e := range w.Events {
			switch e.Type {
			case mvccpb.PUT:
				if err := json.Unmarshal(e.Kv.Value, &Config); err != nil {
					log.Panicf("format err of %s, err: %v", string(e.Kv.Value), err)
				}
				log.Println("update config success...")
			case mvccpb.DELETE:
				log.Panicf("config info is deleted")
			}
		}
	}
}
