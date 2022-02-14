package etcd

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"log"
	"time"
)

var pathTemplate = "/test/probe/%s"

type probe struct {
	ctx      context.Context
	workName string
	ttl      int64
	leaseId  clientv3.LeaseID
}

func NewProbe(ctx context.Context, name string, ttl int64) *probe {
	p := &probe{
		ctx:      ctx,
		workName: name,
		ttl:      ttl,
	}
	leaseId, _ := PutWithLease(ctx, p.probePath(), p.probePath(), ttl)
	p.leaseId = leaseId
	return p
}

func (p *probe) probePath() string {
	return fmt.Sprintf(pathTemplate, p.workName)
}

func (p *probe) postHealth() error {
	_, err := client.KeepAliveOnce(p.ctx, p.leaseId)
	if err != nil {
		log.Printf("keepalive fail, err: %v\n", err)
	}
	return err
}

func (p *probe) watchLifeProbe() {
	for Get(p.ctx, p.probePath()) != "" {
		log.Printf("%s is healthy\n", p.workName)
		time.Sleep(time.Duration(time.Second.Nanoseconds() * (p.ttl - 1)))
	}
	log.Printf("%s is unhealthy\n", p.workName)
}
