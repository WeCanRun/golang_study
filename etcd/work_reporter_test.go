package etcd

import (
	"context"
	"math/rand"
	"testing"
)

func TestWorkReporter(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer client.Close()
	defer Clear(ctx)
	reporter := NewWorkReporter(ctx, "workReporter")
	go reporter.WatchNotify()
	for i := 0; i < 100; {
		i += rand.Intn(10)
		if i > 100 {
			i = 100
		}
		err := reporter.Report(i)
		if err != nil {
			t.Fatalf("err: %v\n", err)
		}
	}
}
