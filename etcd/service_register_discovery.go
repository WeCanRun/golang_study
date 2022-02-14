package etcd

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"log"
	"sync"
)

var DiscoveryRoot = "/discovery"

type serviceRegisterDiscovery struct {
	mux              sync.Mutex
	ctx              context.Context
	leaseIdMap       map[string]clientv3.LeaseID
	serviceName      string
	keepAliveChanMap map[string]<-chan *clientv3.LeaseKeepAliveResponse
}

func NewServiceRegisterDiscovery(ctx context.Context, serviceName string) *serviceRegisterDiscovery {
	leaseIdMap, keepAliveChanMap := make(map[string]clientv3.LeaseID), make(map[string]<-chan *clientv3.LeaseKeepAliveResponse)
	return &serviceRegisterDiscovery{
		leaseIdMap:       leaseIdMap,
		serviceName:      serviceName,
		keepAliveChanMap: keepAliveChanMap,
		ctx:              ctx,
	}
}

func (s *serviceRegisterDiscovery) key(endpoint string) string {
	return fmt.Sprintf("%s/%s/%s", DiscoveryRoot, s.serviceName, endpoint)
}

func (s *serviceRegisterDiscovery) Register(endpoint string, lease int64) error {
	s.mux.Lock()
	defer s.mux.Unlock()
	key := s.key(endpoint)
	leaseId, err := PutWithLease(s.ctx, key, endpoint, lease)
	if err != nil {
		return err
	}

	// 自动续约
	keepAliveChan, err := client.KeepAlive(s.ctx, leaseId)
	if err != nil {
		return err
	}

	s.leaseIdMap[key], s.keepAliveChanMap[key] = leaseId, keepAliveChan
	return nil
}

func (s *serviceRegisterDiscovery) ListenLeaseChan() {
	for len(s.keepAliveChanMap) == 0 {
	}

	var length = len(s.keepAliveChanMap)
	for length != 0 {
		var i = 0
		var wg sync.WaitGroup
		wg.Add(length)
		for key, ch := range s.keepAliveChanMap {
			i++
			go func() {
				for c := range ch {
					log.Println(key, "keepalive successful", c.Revision)
				}
			}()
			if i >= length {
				break
			}
			wg.Done()
		}
		wg.Wait()
		length = len(s.keepAliveChanMap)
	}
	log.Printf("all lease of %s is canceled", s.serviceName)
}

func (s *serviceRegisterDiscovery) CloseService() error {
	var err error
	for _, leasId := range s.leaseIdMap {
		_, err = client.Revoke(s.ctx, leasId)
	}
	s.leaseIdMap = nil
	return err
}

func (s *serviceRegisterDiscovery) CancelEndpoint(endpoint string) bool {
	leaseId, ok := s.leaseIdMap[s.key(endpoint)]
	if !ok {
		return false
	}
	_, err := client.Revoke(s.ctx, leaseId)
	delete(s.leaseIdMap, s.key(endpoint))
	delete(s.keepAliveChanMap, s.key(endpoint))
	if err != nil {
		return false
	}
	return true
}

func (s *serviceRegisterDiscovery) Discovery() []string {
	var services []string
	prefix := GetWithPrefix(s.ctx, s.key(""))
	for _, p := range prefix {
		for _, v := range p {
			services = append(services, v)
		}
	}
	return services
}
