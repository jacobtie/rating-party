package rating

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jacobtie/rating-party/server/internal/platform/werrors"
)

func (c *Controller) GetAllByGameID(ctx context.Context, gameID string) ([]*Rating, error) {
	rows, err := c.db.DB.QueryxContext(ctx, `
		SELECT
			rating_id,
			r.game_id,
			r.participant_id,
			wine_id,
			p.username,
			sight_rating,
			aroma_rating,
			taste_rating,
			overall_rating,
			(sight_rating + aroma_rating + taste_rating + overall_rating) AS total_rating,
			comments
		FROM
			rating r
			INNER JOIN participant p ON r.participant_id = p.participant_id
		WHERE
			r.game_id = $1
		;
	`, gameID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []*Rating{}, nil
		}
		return nil, fmt.Errorf("[rating.GetAllByGameID] failed to query all ratings: %w", err)
	}
	defer rows.Close()
	ratings := make([]*Rating, 0)
	for rows.Next() {
		var rating Rating
		if err := rows.Scan(
			&rating.RatingID,
			&rating.GameID,
			&rating.ParticipantID,
			&rating.WineID,
			&rating.Username,
			&rating.SightRating,
			&rating.AromaRating,
			&rating.TasteRating,
			&rating.OverallRating,
			&rating.TotalRating,
			&rating.Comments,
		); err != nil {
			return nil, fmt.Errorf("[rating.GetAllByGameID] failed to scan row: %w", err)
		}
		ratings = append(ratings, &rating)
	}
	return ratings, nil
}

func (c *Controller) GetAllByGameIDAndParticipantID(ctx context.Context, gameID, participantID string) ([]*Rating, error) {
	rows, err := c.db.DB.QueryxContext(ctx, `
		SELECT
			rating_id,
			game_id,
			participant_id,
			wine_id,
			sight_rating,
			aroma_rating,
			taste_rating,
			overall_rating,
			(sight_rating + aroma_rating + taste_rating + overall_rating) AS total_rating,
			comments
		FROM
			rating
		WHERE
			game_id = $1
			AND participant_id = $2
		;
	`, gameID, participantID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []*Rating{}, nil
		}
		return nil, fmt.Errorf("[rating.GetAllByGameIDAndParticipantID] failed to query all ratings: %w", err)
	}
	defer rows.Close()
	ratings := make([]*Rating, 0)
	for rows.Next() {
		var rating Rating
		if err := rows.Scan(
			&rating.RatingID,
			&rating.GameID,
			&rating.ParticipantID,
			&rating.WineID,
			&rating.SightRating,
			&rating.AromaRating,
			&rating.TasteRating,
			&rating.OverallRating,
			&rating.TotalRating,
			&rating.Comments,
		); err != nil {
			return nil, fmt.Errorf("[rating.GetAllByGameIDAndParticipantID] failed to scan row: %w", err)
		}
		ratings = append(ratings, &rating)
	}
	return ratings, nil
}

func (c *Controller) UpsertRating(ctx context.Context, rating *Rating) (*Rating, error) {
	existingRating, err := c.GetRatingByParticipantIDAndWineID(ctx, rating.ParticipantID, rating.WineID)
	if err != nil {
		if errors.Is(err, werrors.ErrNotFound) {
			return c.CreateRating(ctx, rating)
		}
		return nil, fmt.Errorf("[rating.UpsertRating] failed to get existing rating: %w", err)
	}
	updatedRating, err := c.UpdateRating(ctx, existingRating.RatingID, rating)
	if err != nil {
		return nil, fmt.Errorf("[rating.UpsertRating] failed to update rating: %w", err)
	}
	return updatedRating, nil
}

func (c *Controller) CreateRating(ctx context.Context, rating *Rating) (*Rating, error) {
	ratingID := uuid.New().String()
	if _, err := c.db.DB.ExecContext(ctx, `
		INSERT INTO rating (
			rating_id,
			game_id,
			participant_id,
			wine_id,
			sight_rating,
			aroma_rating,
			taste_rating,
			overall_rating,
			comments
		) VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7,
			$8,
			$9
		)
		;
	`, ratingID, rating.GameID, rating.ParticipantID, rating.WineID, rating.SightRating, rating.AromaRating, rating.TasteRating, rating.OverallRating, rating.Comments); err != nil {
		return nil, fmt.Errorf("[rating.CreateRating] failed to create rating: %w", err)
	}
	createdRating, err := c.GetRatingByParticipantIDAndWineID(ctx, rating.ParticipantID, rating.WineID)
	if err != nil {
		return nil, fmt.Errorf("[rating.CreateRating] failed to get created rating: %w", err)
	}
	return createdRating, nil
}

func (c *Controller) UpdateRating(ctx context.Context, ratingID string, rating *Rating) (*Rating, error) {
	if _, err := c.db.DB.ExecContext(ctx, `
		UPDATE
			rating
		SET
			sight_rating = $1,
			aroma_rating = $2,
			taste_rating = $3,
			overall_rating = $4,
			comments = $5
		WHERE
			rating_id = $6
		;
	`, rating.SightRating, rating.AromaRating, rating.TasteRating, rating.OverallRating, rating.Comments, ratingID); err != nil {
		return nil, fmt.Errorf("[rating.UpdateRating] failed to update rating: %w", err)
	}
	updatedRating, err := c.GetRatingByParticipantIDAndWineID(ctx, rating.ParticipantID, rating.WineID)
	if err != nil {
		return nil, fmt.Errorf("[rating.UpdateRating] failed to get updated rating: %w", err)
	}
	return updatedRating, nil
}

func (c *Controller) GetRatingByParticipantIDAndWineID(ctx context.Context, participantID, wineID string) (*Rating, error) {
	row := c.db.DB.QueryRowxContext(ctx, `
		SELECT
			rating_id,
			game_id,
			participant_id,
			wine_id,
			sight_rating,
			aroma_rating,
			taste_rating,
			overall_rating,
			comments
		FROM
			rating
		WHERE
			participant_id = $1
			AND wine_id = $2
		;
	`, participantID, wineID)
	var rating Rating
	if err := row.Scan(
		&rating.RatingID,
		&rating.GameID,
		&rating.ParticipantID,
		&rating.WineID,
		&rating.SightRating,
		&rating.AromaRating,
		&rating.TasteRating,
		&rating.OverallRating,
		&rating.Comments,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, werrors.ErrNotFound
		}
		return nil, fmt.Errorf("[rating.GetRating] failed to scan row: %w", err)
	}
	return &rating, nil
}
