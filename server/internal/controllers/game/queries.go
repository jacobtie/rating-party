package game

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

func (c *Controller) GetAll(ctx context.Context) ([]*Game, error) {
	rows, err := c.db.DB.QueryxContext(ctx, `
		SELECT
			BIN_TO_UUID(game_id),
			game_name,
			game_code,
			created_at,
			updated_at
		FROM
			game
		;
	`)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("[game.GetAll] failed to query all games: %w", err)
	}
	games := make([]*Game, 0)
	for rows.Next() {
		var game Game
		if err := rows.Scan(
			&game.GameID,
			&game.GameName,
			&game.GameCode,
			&game.CreatedAt,
			&game.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("[game.GetAll] failed to scan row: %w", err)
		}
		games = append(games, &game)
	}
	return games, nil
}

func (c *Controller) GetSingle(ctx context.Context, gameID string) (*Game, error) {
	row := c.db.DB.QueryRowxContext(ctx, `
		SELECT
			BIN_TO_UUID(game_id),
			game_name,
			game_code,
			created_at,
			updated_at
		FROM
			game
		WHERE game_id = UUID_TO_BIN(?)
		;
	`, gameID)
	var game Game
	if err := row.Scan(
		&game.GameID,
		&game.GameName,
		&game.GameCode,
		&game.CreatedAt,
		&game.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("[game.GetSingle] failed to scan row: %w", err)
	}
	return &game, nil
}
