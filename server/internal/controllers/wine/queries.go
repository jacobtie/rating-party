package wine

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jacobtie/rating-party/server/internal/platform/werrors"
)

func (c *Controller) GetAllWines(ctx context.Context, gameID string) ([]*Wine, error) {
	rows, err := c.db.QueryxContext(ctx, `
		SELECT BIN_TO_UUID(wine_id), wine_name, wine_code, wine_year FROM wine WHERE game_id = UUID_TO_BIN(?) ORDER BY wine_code ASC
	`, gameID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []*Wine{}, nil
		}
		return nil, fmt.Errorf("[wine.GetAllWines] failed to query: %w", err)
	}
	defer rows.Close()
	wines := make([]*Wine, 0)
	for rows.Next() {
		var wine Wine
		if err := rows.Scan(
			&wine.WineID,
			&wine.WineName,
			&wine.WineCode,
			&wine.WineYear,
		); err != nil {
			return nil, fmt.Errorf("[wine.GetAllWines] failed to scan row: %w", err)
		}
		wines = append(wines, &wine)
	}
	return wines, nil
}

func (c *Controller) GetSingleWine(ctx context.Context, wineID string) (*Wine, error) {
	row := c.db.QueryRowxContext(ctx, `
		SELECT BIN_TO_UUID(wine_id), wine_name, wine_code, wine_year FROM wine WHERE wine_id = UUID_TO_BIN(?)
	`, wineID)
	var wine Wine
	if err := row.Scan(
		&wine.WineID,
		&wine.WineName,
		&wine.WineCode,
		&wine.WineYear,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("[wine.GetSingleWine] no wine found: %w", werrors.ErrNotFound)
		}
		return nil, fmt.Errorf("[wine.GetSingleWine] failed to scan row: %w", err)
	}
	return &wine, nil
}

func (c *Controller) CreateWine(ctx context.Context, gameID, wineName, wineCode string, wineYear int) (*Wine, error) {
	wineID := uuid.New().String()
	if _, err := c.db.ExecContext(ctx, `
		INSERT INTO wine (wine_id, game_id, wine_name, wine_code, wine_year) VALUES (UUID_TO_BIN(?), UUID_TO_BIN(?), ?, ?, ?)
	`, wineID, gameID, wineName, wineCode, wineYear); err != nil {
		return nil, fmt.Errorf("[wine.CreateWine] failed to insert wine: %w", err)
	}
	wine, err := c.GetSingleWine(ctx, wineID)
	if err != nil {
		return nil, fmt.Errorf("[wine.CreateWine] failed to get wine: %w", err)
	}
	return wine, nil
}

func (c *Controller) UpdateWine(ctx context.Context, wineID, wineName, wineCode string, wineYear int) error {
	if _, err := c.GetSingleWine(ctx, wineID); err != nil {
		return fmt.Errorf("[wine.UpdateWine] failed to get wine: %w", err)
	}
	if _, err := c.db.ExecContext(ctx, `
		UPDATE wine SET wine_name = ?, wine_code = ?, wine_year = ? WHERE wine_id = UUID_TO_BIN(?)
	`, wineName, wineCode, wineYear, wineID); err != nil {
		return fmt.Errorf("[wine.UpdateWine] failed to update wine: %w", err)
	}
	return nil
}

func (c *Controller) DeleteWine(ctx context.Context, wineID string) (*Wine, error) {
	wine, err := c.GetSingleWine(ctx, wineID)
	if err != nil {
		return nil, fmt.Errorf("[wine.DeleteWine] failed to get wine: %w", err)
	}
	if _, err := c.db.ExecContext(ctx, `
		DELETE FROM wine WHERE wine_id = UUID_TO_BIN(?)
	`, wineID); err != nil {
		return nil, fmt.Errorf("[wine.DeleteWine] failed to delete wine: %w", err)
	}
	return wine, nil
}
