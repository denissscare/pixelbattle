package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"pixelbattle/internal/config"
	"pixelbattle/internal/pixcelbattle/domain"
	"pixelbattle/pkg/logger"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	db  *redis.Client
	log logger.Logger
}

const (
	canvasKey = "canvas_state"
)

func NewClient(ctx context.Context, log logger.Logger, cfg config.Config) (*Redis, error) {
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
		log.Errorf("failed to connect to redis server: %s", err.Error())
		return nil, err
	}

	return &Redis{db: db, log: log}, nil
}

func (r *Redis) GetCanvas(ctx context.Context) (map[string]domain.Pixel, error) {
	raw, err := r.db.HGetAll(ctx, canvasKey).Result()
	if err != nil {
		r.log.Errorf("failed to get canvas: %s", err)
		return nil, err
	}

	canvas := make(map[string]domain.Pixel, len(raw))

	for coord, blob := range raw {
		var p domain.Pixel
		if err := json.Unmarshal([]byte(blob), &p); err != nil {
			r.log.Errorf("unmarshal pixcel %s: %v", coord, err)
			continue
		}
		canvas[coord] = p
	}
	return canvas, nil
}

func (r *Redis) SetPixcel(ctx context.Context, p domain.Pixel) error {
	blob, err := json.Marshal(p)
	if err != nil {
		r.log.Errorf("failed marshal pixel: %s", err)
		return err
	}
	field := fmt.Sprintf("%d:%d", p.X, p.Y)
	if err := r.db.HSet(ctx, canvasKey, field, blob).Err(); err != nil {
		r.log.Errorf("failed to set pixel: %s", err)
		return err
	}
	return nil
}
