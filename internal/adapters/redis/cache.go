package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"hakathon-mvp/internal/domain/models"

	"github.com/redis/go-redis/v9"
)

type CitizenReportCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewCitizenReportCache(client *redis.Client, ttl time.Duration) *CitizenReportCache {
	return &CitizenReportCache{
		client: client,
		ttl:    ttl,
	}
}

func (c *CitizenReportCache) Get(ctx context.Context, key string) (*models.CitizenReport, error) {
	data, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Ключ не найден - это не ошибка
		}
		return nil, fmt.Errorf("failed to get from cache: %w", err)
	}

	var res models.CitizenReport
	if err := json.Unmarshal([]byte(data), &res); err != nil {
		return nil, fmt.Errorf("failed to unmarshal citizen report: %w", err)
	}

	return &res, nil
}

func (c *CitizenReportCache) Set(ctx context.Context, key string, citizenReport *models.CitizenReport) error {
	data, err := json.Marshal(citizenReport)
	if err != nil {
		return fmt.Errorf("failed to marshal citizen report: %w", err)
	}

	err = c.client.Set(ctx, key, data, c.ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to set cache: %w", err)
	}

	return nil
}

func (c *CitizenReportCache) Delete(ctx context.Context, key string) error {
	err := c.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete from cache: %w", err)
	}
	return nil
}

func (c *CitizenReportCache) SetList(ctx context.Context, key string, citizenReports []*models.CitizenReport) error {
	data, err := json.Marshal(citizenReports)
	if err != nil {
		return fmt.Errorf("failed to marshal citizen reports list: %w", err)
	}

	err = c.client.Set(ctx, key, data, c.ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to set list cache: %w", err)
	}

	return nil
}

func (c *CitizenReportCache) GetList(ctx context.Context, key string) ([]*models.CitizenReport, error) {
	data, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Ключ не найден - это не ошибка
		}
		return nil, fmt.Errorf("failed to get list from cache: %w", err)
	}

	var products []*models.CitizenReport
	if err := json.Unmarshal([]byte(data), &products); err != nil {
		return nil, fmt.Errorf("failed to unmarshal products list: %w", err)
	}

	return products, nil
}
