package service

import (
	"context"
	"pixelbattle/internal/config"
	"pixelbattle/pkg/logger"
	"pixelbattle/internal/pixcelbattle/metrics"
	"github.com/stretchr/testify/require"
	"testing"
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