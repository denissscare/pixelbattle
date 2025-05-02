package main

import (
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	nc, _ := nats.Connect("nats://localhost:4222")
	ch := make(chan *nats.Msg, 1)
	nc.Subscribe("canvas.updates", func(m *nats.Msg) {
		fmt.Println("RECV:", string(m.Data))
		ch <- m
	})
	nc.Publish("canvas.updates", []byte(`{"x":9,"y":9,"color":"blue"}`))
	select {
	case <-ch:
	case <-time.After(time.Second):
		fmt.Println("no msg received")
	}
}
