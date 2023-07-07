package game

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jacobtie/rating-party/server/internal/platform/werrors"
)

func (c *Controller) GetAll(ctx context.Context) ([]*Game, error) {
	rows, err := c.db.DB.QueryxContext(ctx, `
		SELECT
			game_id,
			game_name,
			game_code,
			is_running,
			are_results_shared
		FROM
			game
		;
	`)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []*Game{}, nil
		}
		return nil, fmt.Errorf("[game.GetAll] failed to query all games: %w", err)
	}
	defer rows.Close()
	games := make([]*Game, 0)
	for rows.Next() {
		var game Game
		if err := rows.Scan(
			&game.GameID,
			&game.GameName,
			&game.GameCode,
			&game.IsRunning,
			&game.AreResultsShared,
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
			game_id,
			game_name,
			game_code,
			is_running,
			are_results_shared
		FROM
			game
		WHERE game_id = $1
		;
	`, gameID)
	var game Game
	if err := row.Scan(
		&game.GameID,
		&game.GameName,
		&game.GameCode,
		&game.IsRunning,
		&game.AreResultsShared,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("[game.GetSingle] no game found: %w", werrors.ErrNotFound)
		}
		return nil, fmt.Errorf("[game.GetSingle] failed to scan row: %w", err)
	}
	return &game, nil
}

func (c *Controller) Create(ctx context.Context, gameName string) (*Game, error) {
	gameID := uuid.New().String()
	gameCode := c.GenerateGameCode()
	if _, err := c.db.DB.ExecContext(ctx, `
		INSERT INTO game (game_id, game_name, game_code) VALUES ($1, $2, $3);
	`, gameID, gameName, gameCode); err != nil {
		return nil, fmt.Errorf("[game.Create] failed to create game: %w", err)
	}
	game, err := c.GetSingle(ctx, gameID)
	if err != nil {
		return nil, fmt.Errorf("[game.Create] failed to get created game: %w", err)
	}
	return game, nil
}

func (c *Controller) Update(ctx context.Context, gameID, gameName string, isRunning, areResultsShared bool) error {
	if _, err := c.GetSingle(ctx, gameID); err != nil {
		return fmt.Errorf("[game.Update] failed to get game: %w", err)
	}
	if _, err := c.db.DB.ExecContext(ctx, `
		UPDATE game SET game_name = $1, is_running = $2, are_results_shared = $3 WHERE game_id = $4;
	`, gameName, isRunning, areResultsShared, gameID); err != nil {
		return fmt.Errorf("[game.Update] failed to update game: %w", err)
	}
	return nil
}

func (c *Controller) Delete(ctx context.Context, gameID string) (*Game, error) {
	game, err := c.GetSingle(ctx, gameID)
	if err != nil {
		return nil, fmt.Errorf("[game.Delete] failed to get game: %w", err)
	}
	if _, err := c.db.DB.ExecContext(ctx, `
		DELETE FROM game WHERE game_id = $1;
	`, gameID); err != nil {
		return nil, fmt.Errorf("[game.Delete] failed to delete game: %w", err)
	}
	return game, nil
}
