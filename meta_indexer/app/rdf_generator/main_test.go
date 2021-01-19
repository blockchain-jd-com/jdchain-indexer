package main

import (
	"testing"
	"time"
)

func TestChan(t *testing.T) {
	receiver := make(chan int, 3)
	go func() {
		for {
			i := <-receiver
			t.Log("received...", i)
			if i >= 3 {
				t.Log("quit receiving at ", i)
				break
			}
		}
	}()
	time.Sleep(1 * time.Second)
	go func() {
		for i := 0; i < 10; i++ {
			select {
			case receiver <- i:
				t.Log("send ...", i)
			default:
				t.Log("failed to send ", i)
			}
		}
	}()
	time.Sleep(3 * time.Second)
}
