package etcd

import (
	"strconv"
	"testing"
)

var (
	ctx   = Context()
	key   = "test"
	value = "testValue"
)

func TestPut(t *testing.T) {
	Put(ctx, key, value)
}

func TestGet(t *testing.T) {
	Get(ctx, key)
}

func TestDel(t *testing.T) {
	Del(ctx, key)
}

func TestWatch(t *testing.T) {
	defer GetClient().Close()
	ctx = Context()
	go Watch(ctx, "/")
	for i := 0; i < 5; i++ {
		strI := strconv.Itoa(i)
		Put(ctx, key+strI, value)
	}
}

func TestPutWithLease(t *testing.T) {
	PutWithLease(ctx, key+"/timeout", value, 1)
	get := Get(ctx, key+"/timeout")
	t.Log(get)
}

func TestGetWithPrefix(t *testing.T) {
	values := GetWithPrefix(ctx, key)
	t.Log(values)
}

func TestMutex(t *testing.T) {
	Mutex(ctx)
}

func TestLock(t *testing.T) {
	Lock()
}

func TestCampaign(t *testing.T) {
	Campaign(ctx)
}

func TestClear(t *testing.T) {
	Clear(ctx)
}
