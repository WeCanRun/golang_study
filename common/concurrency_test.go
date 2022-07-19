package common

import (
	"errors"
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"

	"golang.org/x/sync/errgroup"
)

var mu sync.Mutex
var chain string

func TestMutex(t *testing.T) {
	chain = "main"
	AA()
	fmt.Println(chain)
}
func AA() {
	mu.Lock()
	defer mu.Unlock()
	chain = chain + " --> A"
	B()
}
func B() {
	chain = chain + " --> B"
	C()
}
func C() {
	mu.Lock()
	defer mu.Unlock()
	chain = chain + " --> C"
}

func TestWaitGroup(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		time.Sleep(time.Second)
		wg.Done()
		wg.Add(1)
	}()
	wg.Wait()
}

func TestOnce(t *testing.T) {
	once := sync.Once{}
	once.Do(func() {
		t.Log(a)
	})

	once.Do(func() {
		t.Log(b)
	})
}

type MyMutex struct {
	count int
	sync.Mutex
}

func TestMyMutex(t *testing.T) {
	var count int
	var mu, mu2 sync.Mutex
	t.Logf("%#v %#v", mu, mu2)

	func() {
		mu.Lock()
		mu2 = mu
		t.Logf("%#v %#v", mu, mu2)
		count++
		mu.Unlock()
		//mu2.Unlock()
		t.Logf("%#v %#v", mu, mu2)

	}()

	mu2.Lock()
	count++
	mu2.Unlock()
	fmt.Println(count, count)
}

func TestErrGroup(t *testing.T) {
	var g errgroup.Group
	for i := 0; i <= 2; i++ {
		j := i
		g.Go(func() error {
			if j%2 == 0 {
				return nil
			}
			return errors.New("单数")
		})
	}
	t.Log(runtime.NumGoroutine())
	if err := g.Wait(); err != nil {
		t.Error(err)
	}
}
