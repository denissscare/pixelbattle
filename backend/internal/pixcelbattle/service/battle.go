package service

import (
	"context"
	"encoding/json"
	"pixelbattle/internal/pixcelbattle/broker"
	"pixelbattle/internal/pixcelbattle/domain"
	"pixelbattle/internal/pixcelbattle/metrics"
	pgrepo "pixelbattle/internal/pixcelbattle/storage/postgres"
	"pixelbattle/internal/pixcelbattle/storage/redis"
	"pixelbattle/pkg/logger"

	"github.com/nats-io/nats.go"
)

type BattleService struct {
	repo    redis.RedisRepo
	pgRepo  pgrepo.PGRepo
	broker  broker.Broker
	log     *logger.Logger
	metrics metrics.Metrics
}

func NewBattleService(repo redis.RedisRepo, pgRepo pgrepo.PGRepo, broker broker.Broker, log *logger.Logger, metrics metrics.Metrics) *BattleService {
	return &BattleService{repo: repo, pgRepo: pgRepo, broker: broker, log: log, metrics: metrics}
}
func (s *BattleService) InitCanvas(ctx context.Context) (map[string]domain.Pixel, error) {
	return s.repo.GetCanvas(ctx)
}

func (s *BattleService) GetAllPixelHistory(ctx context.Context) ([]domain.Pixel, error) {
	return s.pgRepo.GetAllPixelHistory(ctx)
}

func (s *BattleService) GetLastPixelByAuthor(ctx context.Context, author string) (*domain.Pixel, error) {
	return s.pgRepo.GetLastPixelByAuthor(ctx, author)
}

func (s *BattleService) UpdatePixel(ctx context.Context, p domain.Pixel) error {
	if err := p.Validate(); err != nil {
		return err
	}

	if err := s.repo.SetPixcel(ctx, p); err != nil {
		return err
	}
	s.metrics.IncPixelsPlaced()

	if err := s.pgRepo.SavePixelHistory(ctx, p); err != nil {
		s.log.Warnf("Failed to save pixel history: %v", err)
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
	ch := make(chan domain.Pixel, 100)

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
	s.metrics.IncActiveConnections()

	go func() {
		<-ctx.Done()
		sub.Unsubscribe()
		s.metrics.DecActiveConnections()
		close(ch)
	}()

	return ch, nil
}
