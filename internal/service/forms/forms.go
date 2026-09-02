package forms

import (
	"context"
	"fmt"

	"github.com/hardal7/chrono/internal/db"
	query "github.com/hardal7/chrono/internal/db/sqlc"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/jackc/pgx/v5/pgtype"
)

func CreateFeatureRequest(ctx context.Context, r dto.CreateFeatureRequest) error {
	err := db.Queries.CreateFeatureRequest(ctx, query.CreateFeatureRequestParams{
		Name:     pgtype.Text{String: r.Name, Valid: r.Name != ""},
		Email:    pgtype.Text{String: r.Email, Valid: r.Email != ""},
		Title:    r.Title,
		Problem:  r.Problem,
		Feature:  r.Feature,
		Priority: r.Priority,
	})

	return fmt.Errorf("Failed to create feature request: %w: %w", db.ErrRunQuery, err)
}

func CreateBugReport(ctx context.Context, r dto.CreateBugReport) error {
	err := db.Queries.CreateBugReport(ctx, query.CreateBugReportParams{
		Name:        pgtype.Text{String: r.Name, Valid: r.Name != ""},
		Email:       pgtype.Text{String: r.Email, Valid: r.Email != ""},
		Title:       r.Title,
		Description: r.Description,
		Steps:       pgtype.Text{String: r.Steps, Valid: r.Steps != ""},
		Environment: pgtype.Text{String: r.Environment, Valid: r.Environment != ""},
		Additional:  pgtype.Text{String: r.Additional, Valid: r.Additional != ""},
	})

	return fmt.Errorf("Failed to create bug report: %w: %w", db.ErrRunQuery, err)
}
