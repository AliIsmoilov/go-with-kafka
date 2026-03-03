package services

import (
	"hakathon-mvp/internal/domain/models"
)

type CitizenReportValidator struct{}

func NewCitizenReportValidator() *CitizenReportValidator {
	return &CitizenReportValidator{}
}

func (v *CitizenReportValidator) Validate(report *models.CitizenReport) error {
	// Базовые проверки
	return nil
}
