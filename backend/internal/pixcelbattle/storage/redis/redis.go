package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"pixelbattle/internal/config"
	"pixelbattle/internal/pixcelbattle/domain"

	"github.com/redis/go-redis/v9"
)

type RedisRepo interface {
	GetCanvas(ctx context.Context) (map[string]domain.Pixel, error)
	SetPixcel(ctx context.Context, p domain.Pixel) error
}

type Redis struct {
	db *redis.Client
}

const (
	canvasKey = "canvas_state"
)

func NewClient(ctx context.Context, cfg config.Config) (*Redis, error) {
	db := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.DB,
		Username:     cfg.Redis.User,
		MaxRetries:   cfg.Redis.MaxRetries,
		DialTimeout:  cfg.Redis.DialTimeout,
		ReadTimeout:  cfg.Redis.Timeout,
		WriteTimeout: cfg.Redis.Timeout,
	})

	if err := db.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &Redis{db: db}, nil
}

func (r *Redis) GetCanvas(ctx context.Context) (map[string]domain.Pixel, error) {
	raw, err := r.db.HGetAll(ctx, canvasKey).Result()
	if err != nil {
		return nil, err
	}

	canvas := make(map[string]domain.Pixel, len(raw))

	for coord, blob := range raw {
		var p domain.Pixel
		if err := json.Unmarshal([]byte(blob), &p); err != nil {
			continue
		}
		canvas[coord] = p
	}
	return canvas, nil
}

func (r *Redis) SetPixcel(ctx context.Context, p domain.Pixel) error {
	blob, err := json.Marshal(p)
	if err != nil {
		return err
	}
	field := fmt.Sprintf("%d:%d", p.X, p.Y)
	if err := r.db.HSet(ctx, canvasKey, field, blob).Err(); err != nil {
		return err
	}
	return nil
}

func (r *Redis) Close(ctx context.Context) error {
	err := r.db.Close()
	if err != nil {
		return err
	}
	return nil
}
