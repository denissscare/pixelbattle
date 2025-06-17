package service

import (
	"context"
	"encoding/json"
	"errors"
	"pixelbattle/internal/pixcelbattle/broker"
	"pixelbattle/internal/pixcelbattle/domain"
	"pixelbattle/internal/pixcelbattle/metrics"
	"pixelbattle/internal/pixcelbattle/storage/postgres"
	pixelredis "pixelbattle/internal/pixcelbattle/storage/redis"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/require"
)

type mockRedis struct {
	canvas          map[string]domain.Pixel
	getCanvasErr    error
	setPixelErr     error
	setPixelLastVal domain.Pixel
}

func (m *mockRedis) GetCanvas(ctx context.Context) (map[string]domain.Pixel, error) {
	return m.canvas, m.getCanvasErr
}
func (m *mockRedis) SetPixcel(ctx context.Context, p domain.Pixel) error {
	m.setPixelLastVal = p
	return m.setPixelErr
}

type mockPGStorage struct {
}

func (pg *mockPGStorage) SavePixelHistory(ctx context.Context, p domain.Pixel) error {
	return nil
}

type mockBroker struct {
	publishedTopic string
	publishedData  []byte
	publishErr     error

	subscribeTopic string
	subscribeFn    func(msg *nats.Msg)
	subscribeErr   error
}

func (m *mockBroker) Publish(subject string, data []byte) error {
	m.publishedTopic = subject
	m.publishedData = data
	return m.publishErr
}
func (m *mockBroker) Subscribe(subject string, handler func(msg *nats.Msg)) (*nats.Subscription, error) {
	m.subscribeTopic = subject
	m.subscribeFn = handler
	return &nats.Subscription{}, m.subscribeErr
}
func (m *mockBroker) Close() {}

var _ broker.Broker = (*mockBroker)(nil)

type mockMetrics struct {
	incPixels, incConns, decConns int
}

func (m *mockMetrics) IncPixelsPlaced()                                          { m.incPixels++ }
func (m *mockMetrics) IncPixelErrors()                                           {}
func (m *mockMetrics) IncHTTPRequest(method, path, status string)                {}
func (m *mockMetrics) ObserveHTTPDuration(method, path string, duration float64) {}
func (m *mockMetrics) IncActiveConnections()                                     { m.incConns++ }
func (m *mockMetrics) DecActiveConnections()                                     { m.decConns++ }

func newTestBattleService(redis pixelredis.RedisRepo, pgstorage postgres.PGRepo, broker broker.Broker, metrics metrics.Metrics) *BattleService {
	return NewBattleService(redis, pgstorage, broker, nil, metrics)
}

func TestInitCanvasM(t *testing.T) {
	ctx := context.Background()
	testPixels := map[string]domain.Pixel{
		"p1": {
			X:         1,
			Y:         1,
			Color:     "#FF0000",
			Author:    "test",
			Timestamp: time.Now(),
		},
	}
	mredis := &mockRedis{canvas: testPixels}
	bs := newTestBattleService(mredis, &mockPGStorage{}, &mockBroker{}, &mockMetrics{})
	canvas, err := bs.InitCanvas(ctx)
	require.NoError(t, err)
	require.Equal(t, testPixels, canvas)

	mredis.getCanvasErr = errors.New("err!")
	_, err = bs.InitCanvas(ctx)
	require.Error(t, err)
}

func TestUpdatePixelSuccess(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	p := domain.Pixel{
		X:         2,
		Y:         2,
		Color:     "#FFFFFF",
		Author:    "user",
		Timestamp: now,
	}
	mredis := &mockRedis{}
	mbroker := &mockBroker{}
	mmetrics := &mockMetrics{}

	bs := newTestBattleService(mredis, &mockPGStorage{}, mbroker, mmetrics)
	err := bs.UpdatePixel(ctx, p)
	require.NoError(t, err)

	require.Equal(t, p.X, mredis.setPixelLastVal.X)
	require.Equal(t, p.Y, mredis.setPixelLastVal.Y)
	require.Equal(t, p.Color, mredis.setPixelLastVal.Color)
	require.Equal(t, p.Author, mredis.setPixelLastVal.Author)
	require.WithinDuration(t, p.Timestamp, mredis.setPixelLastVal.Timestamp, time.Millisecond)

	var published domain.Pixel
	require.NoError(t, json.Unmarshal(mbroker.publishedData, &published))
	require.Equal(t, p.X, published.X)
	require.Equal(t, p.Y, published.Y)
	require.Equal(t, p.Color, published.Color)
	require.Equal(t, p.Author, published.Author)
	require.WithinDuration(t, p.Timestamp, published.Timestamp, time.Millisecond)
}
func TestUpdatePixelValidationError(t *testing.T) {
	ctx := context.Background()

	bs := newTestBattleService(&mockRedis{}, &mockPGStorage{}, &mockBroker{}, &mockMetrics{})

	badPixel := domain.Pixel{X: -1, Y: -1, Color: "zzzzz"}
	err := bs.UpdatePixel(ctx, badPixel)
	require.Error(t, err)
}
func TestUpdatePixelRedisError(t *testing.T) {
	ctx := context.Background()
	mredis := &mockRedis{setPixelErr: errors.New("redis error")}
	bs := newTestBattleService(mredis, &mockPGStorage{}, &mockBroker{}, &mockMetrics{})

	p := domain.Pixel{
		X:         1,
		Y:         2,
		Color:     "#000000",
		Author:    "u",
		Timestamp: time.Now(),
	}
	require.Error(t, bs.UpdatePixel(ctx, p))
}
func TestUpdatePixelBrokerError(t *testing.T) {
	ctx := context.Background()
	mredis := &mockRedis{}
	mbroker := &mockBroker{publishErr: errors.New("publish error")}
	bs := newTestBattleService(mredis, &mockPGStorage{}, mbroker, &mockMetrics{})

	p := domain.Pixel{
		X:         1,
		Y:         2,
		Color:     "#AAAAAA",
		Author:    "u",
		Timestamp: time.Now(),
	}
	require.Error(t, bs.UpdatePixel(ctx, p))
}
