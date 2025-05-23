package service

import (
	"context"
	"encoding/json"
	"pixelbattle/internal/pixcelbattle/broker"
	"pixelbattle/internal/pixcelbattle/domain"
	"pixelbattle/internal/pixcelbattle/storage/redis"
	"pixelbattle/pkg/logger"

	"github.com/nats-io/nats.go"
)

type BattleService struct {
	repo   redis.Redis
	broker broker.NatsBroker
	log    *logger.Logger
}

func NewBattleService(repo redis.Redis, broker broker.NatsBroker, log *logger.Logger) *BattleService {
	return &BattleService{repo: repo, broker: broker, log: log}
}
func (s *BattleService) InitCanvas(ctx context.Context) (map[string]domain.Pixel, error) {
	return s.repo.GetCanvas(ctx)
}

func (s *BattleService) UpdatePixel(ctx context.Context, p domain.Pixel) error {
	if err := p.Validate(); err != nil {
		return err
	}

	if err := s.repo.SetPixcel(ctx, p); err != nil {
		return err
	}
	data, err := json.Marshal(p)
	if err != nil {
		return err
	}
	if err := s.broker.Publish("canvas.updates", data); err != nil {
		return err
	}
	return nil
}

func (s *BattleService) Stream(ctx context.Context) (<-chan domain.Pixel, error) {
	ch := make(chan domain.Pixel, 100) // небольшой буфер на всякий случай

	handler := func(msg *nats.Msg) {
		var p domain.Pixel
		if err := json.Unmarshal(msg.Data, &p); err != nil {
			return
		}
		select {
		case ch <- p:
		case <-ctx.Done():
		}
	}

	sub, err := s.broker.Subscribe("canvas.updates", handler)
	if err != nil {
		return nil, err
	}

	go func() {
		<-ctx.Done()
		sub.Unsubscribe()
		close(ch)
	}()

	return ch, nil
}
