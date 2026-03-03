package v1

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"hakathon-mvp/internal/domain/models"
	"hakathon-mvp/internal/usecases"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type CitizenReportHandler struct {
	citizenReportUC *usecases.CitizenReportUseCase
}

func NewCitizenReportHandler(citizenReportUC *usecases.CitizenReportUseCase) *CitizenReportHandler {
	return &CitizenReportHandler{
		citizenReportUC: citizenReportUC,
	}
}

func (h *CitizenReportHandler) CreateCitizenReport(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CreateCitizenReport
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrorResponse{
			Error:   err.Error(),
			Message: "Invalid JSON format",
		})
		return
	}

	citizenReport := &models.CitizenReport{
		Id:                 uuid.New(),
		RegionID:           req.RegionID,
		DistrictID:         req.DistrictID,
		InfrastructureName: req.InfrastructureName,
		SectorID:           req.SectorID,
		Description:        req.Description,
		PhotoPath:          req.PhotoPath,
	}

	if err := h.citizenReportUC.CreateCitizenReport(ctx, citizenReport); err != nil {
		switch {
		case strings.Contains(err.Error(), "validation failed"):
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, ErrorResponse{
				Error:   "validation_error",
				Message: err.Error(),
			})
		default:
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, ErrorResponse{
				Error:   err.Error(),
				Message: "Failed to create citizen report",
			})
		}
		return
	}

	render.Status(r, http.StatusAccepted) // 202 Accepted - операция принята в обработку
	render.JSON(w, r, CreateCitizenReportResponse{
		ID:      citizenReport.Id.String(),
		Message: "Citizen report creation accepted",
		Status:  "processing",
	})
}

func (h *CitizenReportHandler) GetCitizenReport(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := chi.URLParam(r, "id")
	citizenReport, err := h.citizenReportUC.GetCitizenReport(ctx, idStr)
	if err != nil {
		if err == models.ErrCitizenReportNotFound {
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, ErrorResponse{
				Error:   "not_found",
				Message: "Citizen report not found",
			})
			return
		}

		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to get product",
		})
		return
	}

	if citizenReport == nil {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, ErrorResponse{
			Error:   "not_found",
			Message: "Citizen report not found",
		})
		return
	}

	render.JSON(w, r, CitizenReportResponse{
		ID:                 citizenReport.Id.String(),
		RegionID:           citizenReport.RegionID,
		DistrictID:         citizenReport.DistrictID,
		InfrastructureName: citizenReport.InfrastructureName,
		SectorID:           citizenReport.SectorID,
		Description:        citizenReport.Description,
	})
}

func (h *CitizenReportHandler) ListCitizenReports(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	filter := models.CitizenReportFilter{
		Limit:  25, // дефолтный лимит
		Offset: 0,
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 && limit <= 100 {
			filter.Limit = limit
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			filter.Offset = offset
		}
	}

	citizenReports, err := h.citizenReportUC.ListCitizenReports(ctx, filter)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, ErrorResponse{
			Error:   err.Error(),
			Message: "Failed to list citizen reports",
		})
		return
	}

	response := ListCitizenReportsResponse{
		CitizenReports: make([]CitizenReportResponse, len(citizenReports)),
		Total:          len(citizenReports),
		Limit:          filter.Limit,
		Offset:         filter.Offset,
	}

	for i, citizenReport := range citizenReports {
		response.CitizenReports[i] = CitizenReportResponse{
			ID:                 citizenReport.Id.String(),
			RegionID:           citizenReport.RegionID,
			DistrictID:         citizenReport.DistrictID,
			InfrastructureName: citizenReport.InfrastructureName,
			SectorID:           citizenReport.SectorID,
			Description:        citizenReport.Description,
		}
	}

	render.JSON(w, r, response)
}

func (h *CitizenReportHandler) UpdateCitizenReport(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := chi.URLParam(r, "id")

	var req UpdateCitizenReportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid JSON format",
		})
		return
	}

	citizenReport := &models.CitizenReport{
		Id:                 uuid.MustParse(idStr),
		RegionID:           req.RegionID,
		DistrictID:         req.DistrictID,
		InfrastructureName: req.InfrastructureName,
		SectorID:           req.SectorID,
		Description:        req.Description,
		PhotoPath:          req.PhotoPath,
	}

	if err := h.citizenReportUC.UpdateCitizenReport(ctx, citizenReport); err != nil {
		switch {
		case err == models.ErrCitizenReportNotFound:
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, ErrorResponse{
				Error:   "not_found",
				Message: "Citizen report not found",
			})
		case strings.Contains(err.Error(), "validation failed"):
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, ErrorResponse{
				Error:   "validation_error",
				Message: err.Error(),
			})
		default:
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, ErrorResponse{
				Error:   "internal_error",
				Message: "Failed to update product",
			})
		}
		return
	}

	render.Status(r, http.StatusAccepted)
	render.JSON(w, r, CitizenReportResponseH{
		Message: "Citizen report update accepted",
		Status:  "processing",
	})
}

func (h *CitizenReportHandler) DeleteCitizenReport(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := chi.URLParam(r, "id")
	// var id uuid.UUID
	// if id, err := uuid.Parse(idStr); err != nil {
	// 	render.Status(r, http.StatusBadRequest)
	// 	render.JSON(w, r, ErrorResponse{
	// 		Error:   "invalid_id",
	// 		Message: "Invalid citizen report ID",
	// 	})
	// 	return
	// }

	if err := h.citizenReportUC.DeleteCitizenReport(ctx, idStr); err != nil {
		if err == models.ErrCitizenReportNotFound {
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, ErrorResponse{
				Error:   "not_found",
				Message: "Citizen report not found",
			})
			return
		}

		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to delete citizen report",
		})
		return
	}

	render.Status(r, http.StatusAccepted)
	render.JSON(w, r, CitizenReportResponseH{
		Message: "Citizen report deletion accepted",
		Status:  "processing",
	})
}
