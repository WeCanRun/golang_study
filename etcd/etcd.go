package etcd

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var client *clientv3.Client

func init() {
	var err error
	client, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})

	if err != nil {
		log.Fatal("init client fail")
	}
	log.Printf("connect success, %s", client.Username)
}

func GetClient() *clientv3.Client {
	return client
}

func Context() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cancel()
	}()
	return ctx
}

func Put(ctx context.Context, key, value string) (err error) {
	_, err = client.Put(ctx, key, value)
	if err != nil {
		log.Printf("put to etcd failed, err:%v\n", err)
	}
	return
}

// 使用租约
func PutWithLease(ctx context.Context, key, value string, ttl int64) (clientv3.LeaseID, error) {
	grant, err := client.Grant(ctx, ttl)
	if err != nil {
		log.Printf("PutWithLease to etcd failed, err:%v\n", err)
		return 0, err
	}
	_, err = client.Put(ctx, key, value, clientv3.WithLease(grant.ID))
	if err != nil {
		log.Printf("PutWithLease to etcd failed, err:%v\n", err)
		return 0, err
	}
	return grant.ID, nil
}

func Del(ctx context.Context, key string) (err error) {
	_, err = client.Delete(ctx, key)
	if err != nil {
		log.Printf("del to etcd failed, err:%v\n", err)
	}
	return
}

func Get(ctx context.Context, key string) string {
	resp, err := client.Get(ctx, key)
	if err != nil {
		log.Printf("get from etcd failed, err:%v\n", err)
		return ""
	}

	return string(resp.Kvs[0].Value)
}

// 获取以key为前缀的kv
func GetWithPrefix(ctx context.Context, key string) (resp []map[string]string) {
	get, err := client.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		log.Printf("GetWithPrefix from etcd failed, err:%v\n", err)
		return resp
	}

	for _, v := range get.Kvs {
		//log.Printf("index: %d, keyValue: %s", k, v)
		m := make(map[string]string)
		m[string(v.Key)] = string(v.Value)
		resp = append(resp, m)
	}

	return resp
}

// 监控机制
func Watch(ctx context.Context, key string) {
	watch := client.Watch(ctx, key, clientv3.WithPrefix())
	for w := range watch {
		for _, ev := range w.Events {
			log.Printf("watch type: %s, k: %s, v: %s\n", ev.Type, string(ev.Kv.Key), string(ev.Kv.Value))
		}
	}

}

func Clear(ctx context.Context) {
	_, err := client.Delete(ctx, "/", clientv3.WithPrefix())
	if err != nil {
		log.Printf("clear etcd fail, %v\n", err)
	}
}

// 基于etcd实现分布式锁 github.com/coreos/etcd/clientv3/concurrency
func Mutex(ctx context.Context) {
	// 获取两个session模拟锁竞争
	session1, err := concurrency.NewSession(client)
	if err != nil {
		log.Printf("get session fail, %v\n", err)
	}

	session2, err := concurrency.NewSession(client)
	if err != nil {
		log.Printf("get session fail, %v\n", err)
	}

	lock := "/test/mutex"
	mutex1 := concurrency.NewMutex(session1, lock)
	mutex2 := concurrency.NewMutex(session2, lock)

	if err := mutex1.Lock(ctx); err != nil {
		log.Printf("mutex1 lock fail, %v", err)
		return
	}

	log.Println("get lock for session1")

	mutex2Locked := make(chan struct{})

	go func() {
		defer close(mutex2Locked)
		// 阻塞知道等待session1释放锁
		if err := mutex2.Lock(ctx); err != nil {
			log.Printf("mutex2 lock fail, %v\n", err)
		}
	}()

	time.Sleep(time.Second)
	// session1 释放锁
	if err := mutex1.Unlock(ctx); err != nil {
		log.Printf("mutex1 unlock fail, %v\n", err)
	}
	log.Println("session1 unlocked")

	<-mutex2Locked
	log.Println("session2 get locked successful")

	if err := mutex2.Unlock(ctx); err != nil {
		log.Printf("mutex2 lock fail, %v\n", err)
	}
}

func Lock() {
	// 获取两个session模拟锁竞争
	session1, err := concurrency.NewSession(client)
	if err != nil {
		log.Printf("get session fail, %v\n", err)
	}

	session2, err := concurrency.NewSession(client)
	if err != nil {
		log.Printf("get session fail, %v\n", err)
	}

	lock := "/test/lock"
	locker1 := concurrency.NewLocker(session1, lock)
	locker2 := concurrency.NewLocker(session2, lock)

	locker1.Lock()
	log.Println("get lock for session1")

	locked := make(chan struct{})

	go func() {
		defer close(locked)
		log.Println("before locker2 lock")
		// 阻塞知道等待session1释放锁
		locker2.Lock()
		log.Println("after locker2 lock")
	}()

	time.Sleep(time.Second)
	// session1 释放锁
	locker1.Unlock()
	log.Println("session1 unlocked")

	<-locked
	log.Println("get lock for session2")

	locker2.Unlock()
	log.Println("session2 unlocked")

}

var (
	LeaderShouldDo = func() {}
	QuitLeader     = func() {}
	Leader         = ""
)

func Campaign(ctx context.Context) {
	for {
		// 为client创建session，并绑定租约，租约的ttl为5s, 如传入的ttl<0, 将使用默认的ttl（60s）
		session, err := concurrency.NewSession(client, concurrency.WithTTL(5))
		if err != nil {
			continue
		}

		// campaignPrefix 是需要监听的目录前缀，发起选举会自动在末尾补 `/`
		campaignPrefix := "/test/campaign"
		election := concurrency.NewElection(session, campaignPrefix)
		hostname, _ := os.Hostname()
		campaignId := fmt.Sprintf("%s-%d", hostname, os.Getpid())
		// 1、以 keyPrefix + leaseIdMap 为 Key, campaignId 为 value创建键值对
		// 2、一直阻塞除非成功选举为 Leader 或者 context 过期
		// 3、成功选举为 Leader 的条件是当前 Key 比 keyPrefix 目录下的其他 Key 的 CreateRevision 小
		// 4、每个 Key 都监听（等待）比自己 CreateRevision 小但 CreateRevision 最大的 Key 被删除
		if err := election.Campaign(ctx, campaignId); err != nil {
			select {
			case <-ctx.Done():
				return
			default:
			}
			continue
		}

		Leader = campaignId
		log.Printf("leader: %s\n", Leader)
		if LeaderShouldDo != nil {
			LeaderShouldDo()
		}

		select {
		case <-session.Done():
			// 放弃 Leader 的身份，重新发起新一轮选举
			election.Resign(ctx)
			session.Close()
			if QuitLeader != nil {
				QuitLeader()
			}
			log.Printf("%s quit leader\n", Leader)
			Leader = ""
			continue
		case <-ctx.Done():
			session.Close()
			if QuitLeader != nil {
				QuitLeader()
			}
			log.Printf("%s quit leader\n", Leader)
			Leader = ""
			return
		}
	}
}
