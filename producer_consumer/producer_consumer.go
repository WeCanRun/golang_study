package producer_consumer

import (
	"fmt"
	"time"
)

func Producer(factor int, out chan<- int) {
	for i := 0; ; i++ {
		time.Sleep(time.Duration(factor * int(time.Second)))
		num := i * factor
		out <- num
		fmt.Printf("producer the number: %d\n", num)
	}
}

func Consumer(factor int, in <-chan int) {
	for v := range in {
		time.Sleep(time.Duration(factor * int(time.Second)))
		fmt.Printf("consumer the num: %d\n", v)
	}
}
