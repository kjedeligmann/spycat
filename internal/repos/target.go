package repos

import (
	"context"
	"database/sql"
	// "github.com/kjedeligmann/spycat/internal/models"
)

type TargetRepo struct {
	db *sql.DB
}

func NewTargetRepo(db *sql.DB) *TargetRepo {
	return &TargetRepo{db: db}
}

// ListByMission lists all targets for a given mission.
// func (r *TargetRepo) ListByMission(ctx context.Context, missionID int) ([]models.Target, error) {
// 	query := `
//         SELECT id, mission_id, name, country, notes, completed
//         FROM targets
//         WHERE mission_id = $1
//     `
// 	rows, err := r.db.QueryContext(ctx, query, missionID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()
//
// 	var targets []models.Target
// 	for rows.Next() {
// 		var t models.Target
// 		if err := rows.Scan(&t.ID, &t.MissionID, &t.Name, &t.Country, &t.Notes, &t.Completed); err != nil {
// 			return nil, err
// 		}
// 		targets = append(targets, t)
// 	}
// 	return targets, nil
// }

// MarkTargetCompleted marks the target as completed.
func (r *TargetRepo) MarkTargetCompleted(ctx context.Context, missionID, targetID int) error {
	query := `
        UPDATE targets
        SET completed = true
        WHERE mission_id = $1 AND id = $2
        RETURNING id
    `
	var id int
	return r.db.QueryRowContext(ctx, query, missionID, targetID).Scan(&id)
}

// UpdateTargetNotes updates the notes of a target if neither the target nor its mission is completed.
func (r *TargetRepo) UpdateTargetNotes(ctx context.Context, missionID, targetID int, notes string) error {
	query := `
        UPDATE targets
        SET notes = $1
        WHERE mission_id = $2
          AND id = $3
          AND completed = false
          AND (SELECT status FROM missions WHERE id = $2) <> 'completed'
        RETURNING id
    `
	var id int
	return r.db.QueryRowContext(ctx, query, notes, missionID, targetID).Scan(&id)
}

// DeleteTarget deletes the target if it is not completed and its mission is not completed.
func (r *TargetRepo) DeleteTarget(ctx context.Context, missionID, targetID int) error {
	query := `
        DELETE FROM targets
        WHERE mission_id = $1
          AND id = $2
          AND completed = false
          AND (SELECT status FROM missions WHERE id = $1) <> 'completed'
        RETURNING id
    `
	var id int
	return r.db.QueryRowContext(ctx, query, missionID, targetID).Scan(&id)
}
