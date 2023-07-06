package participant

import (
	"context"
	"fmt"
)

func (c *Controller) GetAllParticipantsByGameID(ctx context.Context, gameID string) ([]*Participant, error) {
	rows, err := c.db.DB.QueryxContext(ctx, `
		SELECT participant_id, game_id, username FROM participant WHERE game_id = UUID_TO_BIN(?)
	`, gameID)
	if err != nil {
		return nil, fmt.Errorf("[participant.GetAllParticipantsByGameID] failed to query participants: %w", err)
	}
	defer rows.Close()
	participants := make([]*Participant, 0)
	for rows.Next() {
		var p Participant
		if err := rows.Scan(
			&p.ParticipantID,
			&p.GameID,
			&p.Username,
		); err != nil {
			return nil, fmt.Errorf("[participant.GetAllParticipantsByGameID] failed to scan participant: %w", err)
		}
	}
	return participants, nil
}
