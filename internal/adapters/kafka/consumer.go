package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"hakathon-mvp/internal/domain/models"
	"hakathon-mvp/internal/domain/repositories"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
)

type Consumer struct {
	reader       *kafka.Reader
	citizenRepo  repositories.CitizenReportRepository
	cache        repositories.CitizenReportCache
	batchSize    int
	batchTimeout time.Duration
}

func NewConsumer(brokers []string, topic, groupID string, batchSize int, batchTimeout time.Duration) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		Topic:          topic,
		GroupID:        groupID,
		MinBytes:       10e3, // 10KB
		MaxBytes:       10e6, // 10MB
		CommitInterval: time.Second,
		StartOffset:    kafka.LastOffset,
		MaxWait:        batchTimeout,
		QueueCapacity:  batchSize,
	})

	return &Consumer{
		reader:       reader,
		batchSize:    batchSize,
		batchTimeout: batchTimeout,
	}
}

func NewConsumerWithAuth(brokers []string, topic, groupID, username, password string, batchSize int, batchTimeout time.Duration) *Consumer {
	dialer := &kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
		SASLMechanism: plain.Mechanism{
			Username: username,
			Password: password,
		},
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		Topic:          topic,
		GroupID:        groupID,
		MinBytes:       10e3,
		MaxBytes:       10e6,
		CommitInterval: time.Second,
		StartOffset:    kafka.LastOffset,
		MaxWait:        batchTimeout,
		QueueCapacity:  batchSize,
		Dialer:         dialer,
	})

	return &Consumer{
		reader:       reader,
		batchSize:    batchSize,
		batchTimeout: batchTimeout,
	}
}

func (c *Consumer) SetCitizenReportRepo(repo repositories.CitizenReportRepository) {
	c.citizenRepo = repo
}

func (c *Consumer) SetCache(cache repositories.CitizenReportCache) {
	c.cache = cache
}

func (c *Consumer) Start(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			msg, err := c.reader.FetchMessage(ctx)
			if err != nil {
				return fmt.Errorf("failed to fetch message: %w", err)
			}

			if err := c.processMessage(ctx, msg); err != nil {
				fmt.Printf("Failed to process message: %v\n", err)
				continue
			}

			if err := c.reader.CommitMessages(ctx, msg); err != nil {
				return fmt.Errorf("failed to commit message: %w", err)
			}
		}
	}
}

func (c *Consumer) processMessage(ctx context.Context, msg kafka.Message) error {
	var event models.CitizenReportEvent
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		return fmt.Errorf("failed to unmarshal event: %w", err)
	}

	switch event.EventType {
	case models.CitizenReportCreated:
		return c.handleCitizenReportCreated(ctx, &event)
	case models.CitizenReportUpdated:
		return c.handleCitizenReportUpdated(ctx, &event)
	case models.CitizenReportDeleted:
		return c.handleCitizenReportDeleted(ctx, &event)
	default:
		return fmt.Errorf("unknown event type: %s", event.EventType)
	}
}

func (c *Consumer) handleCitizenReportCreated(ctx context.Context, event *models.CitizenReportEvent) error {
	if event.ReportData == nil {
		return fmt.Errorf("citizen report data is nil for create event")
	}

	if err := c.citizenRepo.Create(ctx, event.ReportData); err != nil {
		return fmt.Errorf("failed to create citizen report: %w", err)
	}

	cacheKey := fmt.Sprintf("citizen-report:%s", event.ReportData.Id.String())
	if err := c.cache.Set(ctx, cacheKey, event.ReportData); err != nil {
		fmt.Printf("Failed to cache citizen report %s: %v\n", event.ReportData.Id.String(), err)
	}

	fmt.Printf("Successfully created citizen report: %s\n", event.ReportData.Id.String())
	return nil
}

func (c *Consumer) handleCitizenReportUpdated(ctx context.Context, event *models.CitizenReportEvent) error {
	if event.ReportData == nil {
		return fmt.Errorf("citizen report data is nil for update event")
	}

	if err := c.citizenRepo.Update(ctx, event.ReportData); err != nil {
		return fmt.Errorf("failed to update citizen report: %w", err)
	}

	cacheKey := fmt.Sprintf("citizen-report:%s", event.ReportData.Id.String())
	if err := c.cache.Set(ctx, cacheKey, event.ReportData); err != nil {
		fmt.Printf("Failed to cache citizen report %s: %v\n", event.ReportData.Id.String(), err)
	}

	fmt.Printf("Successfully updated citizen report: %s\n", event.ReportData.Id.String())
	return nil
}

func (c *Consumer) handleCitizenReportDeleted(ctx context.Context, event *models.CitizenReportEvent) error {
	if err := c.citizenRepo.Delete(ctx, event.ReportID); err != nil {
		return fmt.Errorf("failed to delete citizen report: %w", err)
	}

	cacheKey := fmt.Sprintf("citizen-report:%s", event.ReportData.Id.String())
	if err := c.cache.Delete(ctx, cacheKey); err != nil {
		fmt.Printf("Failed to delete citizen report %s from cache: %v\n", event.ReportData.Id.String(), err)
	}

	fmt.Printf("Successfully deleted citizen report: %s\n", event.ReportData.Id.String())
	return nil
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
