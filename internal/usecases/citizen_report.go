package usecases

import (
	"context"
	"fmt"
	"log"
	"time"

	"hakathon-mvp/internal/domain/models"
	"hakathon-mvp/internal/domain/repositories"
	vld "hakathon-mvp/internal/pkg/validator"
)

type CitizenReportUseCase struct {
	repo          repositories.CitizenReportRepository
	cache         repositories.CitizenReportCache
	eventProducer repositories.EventProducer
	validator     *vld.CitizenReportValidator
}

func NewCitizenReportUseCase(
	repo repositories.CitizenReportRepository,
	cache repositories.CitizenReportCache,
	eventProducer repositories.EventProducer,
	validator *vld.CitizenReportValidator,
) *CitizenReportUseCase {
	return &CitizenReportUseCase{
		repo:          repo,
		cache:         cache,
		eventProducer: eventProducer,
		validator:     validator,
	}
}

func (uc *CitizenReportUseCase) CreateCitizenReport(ctx context.Context, citizenReport *models.CitizenReport) error {
	if err := uc.validator.Validate(citizenReport); err != nil {
		log.Printf("Validation failed for citizen report: %v, error: %v", citizenReport, err)
		return err
	}

	event := &models.CitizenReportEvent{
		EventID:    generateEventID(),
		EventType:  models.CitizenReportCreated,
		Timestamp:  time.Now(),
		ReportData: citizenReport,
		ProducerID: "citizen-report-api",
		Sequence:   time.Now().UnixNano(),
	}

	if err := uc.eventProducer.SendCitizenReportEvent(ctx, event); err != nil {
		return err
	}

	return nil
}

func (uc *CitizenReportUseCase) GetCitizenReport(ctx context.Context, id string) (*models.CitizenReport, error) {
	// Пытаемся получить из кеша
	cacheKey := fmt.Sprintf("citizen-report:%s", id)
	if cached, err := uc.cache.Get(ctx, cacheKey); err == nil && cached != nil {
		return cached, nil
	}

	// Если нет в кеше, идем в базу
	product, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// На случай, если реализация репозитория вернет (nil, nil)
	if product == nil {
		return nil, models.ErrCitizenReportNotFound
	}

	// Сохраняем в кеш
	if err := uc.cache.Set(ctx, cacheKey, product); err != nil {
		// Логируем ошибку, но не прерываем выполнение
		fmt.Printf("Failed to cache product: %v\n", err)
	}

	return product, nil
}

func (uc *CitizenReportUseCase) UpdateCitizenReport(ctx context.Context, citizenReport *models.CitizenReport) error {
	// Валидация
	if err := uc.validator.Validate(citizenReport); err != nil {
		return err
	}

	// Создаем событие для Kafka
	event := &models.CitizenReportEvent{
		EventID:    generateEventID(),
		EventType:  models.CitizenReportUpdated,
		Timestamp:  time.Now(),
		ReportID:   citizenReport.Id.String(),
		ReportData: citizenReport,
		ProducerID: "citizen-report-api",
		Sequence:   time.Now().UnixNano(),
	}

	if err := uc.eventProducer.SendCitizenReportEvent(ctx, event); err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("citizen-report:%s", citizenReport.Id.String())
	if err := uc.cache.Delete(ctx, cacheKey); err != nil {
		fmt.Printf("Failed to invalidate cache: %v\n", err)
	}

	return nil
}

func (uc *CitizenReportUseCase) DeleteCitizenReport(ctx context.Context, id string) error {
	event := &models.CitizenReportEvent{
		EventID:    generateEventID(),
		EventType:  models.CitizenReportDeleted,
		Timestamp:  time.Now(),
		ReportID:   id,
		ProducerID: "citizen-report-api",
		Sequence:   time.Now().UnixNano(),
	}

	if err := uc.eventProducer.SendCitizenReportEvent(ctx, event); err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("citizen-report:%s", id)
	if err := uc.cache.Delete(ctx, cacheKey); err != nil {
		fmt.Printf("Failed to invalidate cache: %v\n", err)
	}

	return nil
}

func (uc *CitizenReportUseCase) ListCitizenReports(ctx context.Context, filter models.CitizenReportFilter) ([]*models.CitizenReport, error) {
	// для списков также можно использовать кеш, но это сложнее из-за вариативности фильтров
	return uc.repo.List(ctx, filter)
}

func generateEventID() string {
	return fmt.Sprintf("event-%d", time.Now().UnixNano())
}
