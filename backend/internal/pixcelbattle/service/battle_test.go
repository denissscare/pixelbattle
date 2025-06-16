package service

import (
	"context"
	
	"pixelbattle/internal/config"
	"pixelbattle/pkg/logger"
	"pixelbattle/internal/pixcelbattle/metrics"
	"pixelbattle/internal/pixcelbattle/broker"
	"pixelbattle/internal/pixcelbattle/domain"
	
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestNewBattleService(t *testing.T) {
	ctx := context.Background()
	config := config.LoadConfig()
	
	rds, err := redis.NewClient(ctx, *config)
	require.NoError(t, err)
	defer rds.Close(ctx)
	
	var (
		nilBroker broker.NatsBroker
		nilLogger *logger.Logger
		nilMetrics metrics.Metrics
	)
	
	return NewBattleService(*rds, nilBroker, nilLogger, nilMetrics)
}

func TestInitCanvas(t *testing.T) {
	ctx := context.Background()
	
	bs := TestNewBattleService(t)
	_, err := bs.InitCanvas(ctx)
	require.NoError(t, err)
}

func TestUpdatePixel(t *testing.T) {
	ctx := context.Background()

	p := domain.Pixel {
	    X:         10,
	    Y:         20,
	    Color:     "#FF0000",
	    Author:    "user123",
	    Timestamp: time.Now(),
	}

	require.Nil(t, p.Validate())

	bs := TestNewBattleService(t)
	require.Nil(t, bs.UpdatePixel(ctx, p))
}
