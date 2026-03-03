package validator

import (
	"hakathon-mvp/internal/domain/models"
)

type CitizenReportValidator struct{}

func NewCitizenReportValidator() *CitizenReportValidator {
	return &CitizenReportValidator{}
}

func (v *CitizenReportValidator) Validate(citizenReport *models.CitizenReport) error {
	// Базовые проверки
	return nil
}
