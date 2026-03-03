package v1

type CreateCitizenReport struct {
	RegionID           int64   `json:"region_id" validate:"required"`
	DistrictID         int64   `json:"district_id" validate:"required"`
	InfrastructureName string  `json:"infrastructure_name" validate:"required,max=255"`
	SectorID           int64   `json:"sector_id" validate:"required"`
	Description        *string `json:"description"`
	PhotoPath          *string `json:"photo_path"`
}

type UpdateCitizenReportRequest struct {
	RegionID           int64   `json:"region_id" validate:"required"`
	DistrictID         int64   `json:"district_id" validate:"required"`
	InfrastructureName string  `json:"infrastructure_name" validate:"required,max=255"`
	SectorID           int64   `json:"sector_id" validate:"required"`
	Description        *string `json:"description"`
	PhotoPath          *string `json:"photo_path"`
}

type CitizenReportResponse struct {
	ID                 string  `json:"id"`
	RegionID           int64   `json:"region_id"`
	DistrictID         int64   `json:"district_id"`
	InfrastructureName string  `json:"infrastructure_name"`
	SectorID           int64   `json:"sector_id"`
	Description        *string `json:"description,omitempty"`
}

type CreateCitizenReportResponse struct {
	ID      string `json:"id"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

type CitizenReportResponseH struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

type DeleteCitizenReportResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

type ListCitizenReportsResponse struct {
	CitizenReports []CitizenReportResponse `json:"citizen_reports"`
	Total          int                     `json:"total"`
	Limit          int                     `json:"limit"`
	Offset         int                     `json:"offset"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
