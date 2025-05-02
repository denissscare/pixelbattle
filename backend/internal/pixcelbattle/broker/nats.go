package broker

import (
	"time"

	"github.com/nats-io/nats.go"
)

type Broker interface {
	Publish(subject string, data []byte) error
	Subscribe(subject string, handler func(msg *nats.Msg)) (*nats.Subscription, error)
	Close()
}

type NatsBroker struct {
	conn *nats.Conn
}

func NewBroker(url string) (*NatsBroker, error) {
	opts := []nats.Option{
		nats.Name("pixcelbattle-broker"),
		nats.MaxReconnects(-1),
		nats.ReconnectWait(5 * time.Second),
		nats.Timeout(5 * time.Second),
	}

	nc, err := nats.Connect(url, opts...)
	if err != nil {
		return nil, err
	}
	return &NatsBroker{conn: nc}, nil
}

func (b *NatsBroker) Subscribe(subject string, handler func(msg *nats.Msg)) (*nats.Subscription, error) {
	sub, err := b.conn.Subscribe(subject, handler)
	if err != nil {
		return nil, err
	}
	return sub, nil
}

func (b *NatsBroker) Publish(subject string, data []byte) error {
	return b.conn.Publish(subject, data)
}

func (b *NatsBroker) Close() {
	b.conn.Close()
}
