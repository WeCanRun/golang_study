package producer_consumer

import "testing"

func TestProducerAndConsumer(t *testing.T) {
	ch := make(chan int)
	factor := 1
	go Producer(factor, ch)
	Consumer(factor, ch)
}
