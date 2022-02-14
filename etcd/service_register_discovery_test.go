package etcd

import (
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"
)

func TestNewServiceRegisterDiscovery(t *testing.T) {
	svcRegisterDiscovery := NewServiceRegisterDiscovery(ctx, "web")
	defer svcRegisterDiscovery.CloseService()
	go func() {
		for {
			discovery := svcRegisterDiscovery.Discovery()
			log.Println(discovery)
			time.Sleep(time.Second * 3)
		}
	}()

	go func() {
		for i := 0; ; i++ {
			endpoint := fmt.Sprintf("127.0.0.1:909%d", i)
			err := svcRegisterDiscovery.Register(endpoint, 5)
			if err != nil {
				log.Printf("err: %v", err)
			}
			if i > 5 {
				delEndpoint := fmt.Sprintf("127.0.0.1:909%d", rand.Intn(i))
				for !svcRegisterDiscovery.CancelEndpoint(delEndpoint) {
					delEndpoint = fmt.Sprintf("127.0.0.1:909%d", rand.Intn(i))
				}
			}
			time.Sleep(time.Second * 3)
		}
	}()

	go svcRegisterDiscovery.ListenLeaseChan()
	select {}

}
