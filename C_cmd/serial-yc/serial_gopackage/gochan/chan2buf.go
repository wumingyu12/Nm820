package main

//chan 有缓冲试验
import (
	"fmt"
	"time"
)

var g_cmd = []byte{0x8a, 0x9b, 0x00, 0x01, 0x05, 0x02, 0x00, 0x09, 0x00}

func produce(p chan<- []byte) {
	for i := 0; i < 40; i++ {
		p <- g_cmd
		time.Sleep(1 * time.Second)
		fmt.Printf("send:%d,%x\n", i, g_cmd)
	}
}
func consumer(c <-chan []byte) {
	for {
		v := <-c
		fmt.Printf("receive:%x\n", v)
	}
}
func main() {
	ch := make(chan []byte, 1)
	var aa = []byte{0xff, 0xff, 0x44}
	ch <- aa
	go produce(ch)
	go consumer(ch)
	time.Sleep(100 * time.Second)
}
