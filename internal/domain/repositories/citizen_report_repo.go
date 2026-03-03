package repositories

import (
	"context"

	"hakathon-mvp/internal/domain/models"
)

type CitizenReportRepository interface {
	Create(ctx context.Context, report *models.CitizenReport) error
	GetByID(ctx context.Context, id string) (*models.CitizenReport, error)
	Update(ctx context.Context, report *models.CitizenReport) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filter models.CitizenReportFilter) ([]*models.CitizenReport, error)
}

// CitizenReportCache определяет контракт для кеширования отчетов граждан
type CitizenReportCache interface {
	Get(ctx context.Context, key string) (*models.CitizenReport, error)
	Set(ctx context.Context, key string, report *models.CitizenReport) error
	Delete(ctx context.Context, key string) error
	SetList(ctx context.Context, key string, reports []*models.CitizenReport) error
	GetList(ctx context.Context, key string) ([]*models.CitizenReport, error)
}

// EventProducer определяет контракт для отправки событий в Kafka
type EventProducer interface {
	SendCitizenReportEvent(ctx context.Context, event *models.CitizenReportEvent) error
	Close() error
}

// EventConsumer определяет контракт для потребления событий из Kafka
type EventConsumer interface {
	Start(ctx context.Context) error
	Close() error
}
