package models

import (
	"time"

	"github.com/google/uuid"
)

type CitizenReport struct {
	Id                 uuid.UUID `json:"id"`
	RegionID           int64     `json:"region_id" validate:"required"`
	DistrictID         int64     `json:"district_id" validate:"required"`
	InfrastructureName string    `json:"infrastructure_name" validate:"required,max=255"`
	SectorID           int64     `json:"sector_id" validate:"required"`
	Description        *string   `json:"description"`
	PhotoPath          *string   `json:"photo_path"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type CitizenReportFilter struct {
	RegionID   *int64 `json:"region_id,omitempty"`
	DistrictID *int64 `json:"district_id,omitempty"`
	SectorID   *int64 `json:"sector_id,omitempty"`
	Limit      int    `json:"limit"`
	Offset     int    `json:"offset"`
}
