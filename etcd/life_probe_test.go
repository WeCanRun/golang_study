package etcd

import (
	"context"
	"log"
	"testing"
	"time"
)

func TestNewProbe(t *testing.T) {
	probe := NewProbe(context.Background(), "test", 5)
	go probe.watchLifeProbe()
	go func() {
		for {
			if err := probe.postHealth(); err != nil {
				log.Println("post health fail, err: ", err)
				return
			}
			log.Println("post health success ")
			time.Sleep(time.Second * 4)
		}
	}()
	select {}
}
