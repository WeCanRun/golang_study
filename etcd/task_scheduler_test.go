package etcd

import (
	"context"
	"testing"
	"time"
)

func TestTaskScheduler(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer client.Close()
	defer Clear(ctx)

	scheduler := NewTaskScheduler(ctx, "taskName")
	go scheduler.WatchNotify()
	scheduler.Notify(Scheduled)
	time.Sleep(time.Second)
	scheduler.Notify(Starting)
	time.Sleep(time.Second)
	scheduler.Notify(Running)
	time.Sleep(time.Second)
	scheduler.Notify(Finished)
}
