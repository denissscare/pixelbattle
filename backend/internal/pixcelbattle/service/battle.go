package service

import (
	"context"
	"encoding/json"
	"pixelbattle/internal/pixcelbattle/broker"
	"pixelbattle/internal/pixcelbattle/domain"
	"pixelbattle/internal/pixcelbattle/storage/redis"
	"pixelbattle/pkg/logger"

	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
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
	canvas, err := s.repo.GetCanvas(ctx)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"action":    "GetCanvas",
			"component": "pixcelbattle.service.InitCanvas",
			"success":   false,
		}).Errorf("failed to load canvas: %v", err)
		return nil, err
	}
	return canvas, err
}
func (s *BattleService) UpdatePixel(ctx context.Context, p domain.Pixel) error {
	if err := s.repo.SetPixcel(ctx, p); err != nil {
		s.log.WithFields(logrus.Fields{
			"action":    "SetPixcel",
			"component": "pixcelbatlle.service.UpdatePixel",
			"success":   false,
		}).Errorf("failed to update pixcel: %s", err)
		return err
	}
	data, err := json.Marshal(p)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"action":    "json.Marshal",
			"component": "pixcelbatlle.service.UpdatePixel",
			"success":   false,
		}).Errorf("failed to marshal pixcel: %s", err)
		return err
	}
	if err := s.broker.Publish("canvas.updates", data); err != nil {
		s.log.WithFields(logrus.Fields{
			"action":    "broker.Publish",
			"component": "pixcelbatlle.service.UpdatePixel",
			"success":   false,
		}).Errorf("failed to publish pixcel: %s", err)
		return err
	}
	return nil
}

func (s *BattleService) Stream(ctx context.Context) (<-chan domain.Pixel, error) {
	ch := make(chan domain.Pixel)

	var msgHandler func(msg *nats.Msg) = func(msg *nats.Msg) {
		var p domain.Pixel
		if err := json.Unmarshal(msg.Data, &p); err != nil {
			s.log.WithFields(logrus.Fields{
				"action":    "json.Unmarshal",
				"component": "pixcelbatlle.service.Stream",
				"success":   false,
			}).Errorf("failed to unmarshal pixcel: %s", err)
			return
		}
		select {
		case ch <- p:
		case <-ctx.Done():
		}
	}

	sub, err := s.broker.Subscribe("canvas.updates", msgHandler)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"action":    "broker.Subscribe",
			"component": "pixcelbatlle.service.Stream",
			"success":   false,
		}).Errorf("failed to subscribe to broker: %s", err)
		return nil, err
	}

	go func() {
		<-ctx.Done()
		if err := sub.Unsubscribe(); err != nil {
			s.log.WithFields(logrus.Fields{
				"action":    "sub.Unsubscribe",
				"component": "pixcelbatlle.service.Stream",
				"success":   false,
			}).Errorf("failed to unsubscribe from NATS: %s", err)
		}
		close(ch)
	}()

	return ch, nil
}
