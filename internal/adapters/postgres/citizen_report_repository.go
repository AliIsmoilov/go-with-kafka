package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"hakathon-mvp/internal/domain/models"

	_ "github.com/lib/pq"
)

type CitizenReportRepository struct {
	db *sql.DB
}

func NewCitizenReportRepository(db *sql.DB) *CitizenReportRepository {
	return &CitizenReportRepository{
		db: db,
	}
}

func (r *CitizenReportRepository) Create(ctx context.Context, citizenReport *models.CitizenReport) error {
	// we expect the caller to have generated a UUID for the report and included it
	// in the event payload. store the provided identifier to keep records in sync.
	query := `
		INSERT INTO citizen_reports (id, region_id, district_id, infrastructure_name, sector_id, description, photo_path)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING created_at, updated_at
	`

	err := r.db.QueryRowContext(ctx, query,
		citizenReport.Id,
		citizenReport.RegionID,
		citizenReport.DistrictID,
		citizenReport.InfrastructureName,
		citizenReport.SectorID,
		citizenReport.Description,
		citizenReport.PhotoPath,
	).Scan(&citizenReport.CreatedAt, &citizenReport.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create citizen report: %w", err)
	}

	return nil
}

func (r *CitizenReportRepository) GetByID(ctx context.Context, id string) (*models.CitizenReport, error) {
	query := `
		SELECT id, region_id, district_id, infrastructure_name, sector_id, description, photo_path, created_at, updated_at
		FROM citizen_reports
		WHERE id = $1
	`

	var res models.CitizenReport
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&res.Id,
		&res.RegionID,
		&res.DistrictID,
		&res.InfrastructureName,
		&res.SectorID,
		&res.Description,
		&res.PhotoPath,
		&res.CreatedAt,
		&res.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrCitizenReportNotFound
		}
		return nil, fmt.Errorf("failed to get citizen report: %w", err)
	}

	return &res, nil
}

func (r *CitizenReportRepository) Update(ctx context.Context, citizenReport *models.CitizenReport) error {
	query := `
		UPDATE citizen_reports 
		SET region_id = $1, district_id = $2, infrastructure_name = $3, sector_id = $4, description = $5, photo_path = $6, updated_at = NOW()
		WHERE id = $7
		RETURNING updated_at
	`

	err := r.db.QueryRowContext(ctx, query,
		citizenReport.RegionID,
		citizenReport.DistrictID,
		citizenReport.InfrastructureName,
		citizenReport.SectorID,
		citizenReport.Description,
		citizenReport.PhotoPath,
		citizenReport.Id,
	).Scan(&citizenReport.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.ErrCitizenReportNotFound
		}
		return fmt.Errorf("failed to update citizen report: %w", err)
	}

	return nil
}

func (r *CitizenReportRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM citizen_reports WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete citizen report: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return models.ErrCitizenReportNotFound
	}

	return nil
}

func (r *CitizenReportRepository) List(ctx context.Context, filter models.CitizenReportFilter) ([]*models.CitizenReport, error) {
	query := `
		SELECT id, region_id, district_id, infrastructure_name, sector_id, description, photo_path, created_at, updated_at
		FROM citizen_reports
		WHERE 1=1
	`
	args := []interface{}{}
	argCounter := 1

	// Добавляем условия фильтрации
	if filter.RegionID != nil {
		query += fmt.Sprintf(" AND region_id = $%d", argCounter)
		args = append(args, *filter.RegionID)
		argCounter++
	}

	if filter.DistrictID != nil {
		query += fmt.Sprintf(" AND district_id = $%d", argCounter)
		args = append(args, *filter.DistrictID)
		argCounter++
	}

	if filter.SectorID != nil {
		query += fmt.Sprintf(" AND sector_id = $%d", argCounter)
		args = append(args, *filter.SectorID)
		argCounter++
	}

	query += " ORDER BY created_at DESC"
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argCounter, argCounter+1)
	args = append(args, filter.Limit, filter.Offset)
	argCounter += 2

	// Выполняем запрос
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list citizen reports: %w", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			panic(err)
		}
	}(rows)

	var citizenReports []*models.CitizenReport
	for rows.Next() {
		var citizenReport models.CitizenReport

		err := rows.Scan(
			&citizenReport.Id,
			&citizenReport.RegionID,
			&citizenReport.DistrictID,
			&citizenReport.InfrastructureName,
			&citizenReport.SectorID,
			&citizenReport.Description,
			&citizenReport.PhotoPath,
			&citizenReport.CreatedAt,
			&citizenReport.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan citizen report: %w", err)
		}

		citizenReports = append(citizenReports, &citizenReport)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return citizenReports, nil
}
