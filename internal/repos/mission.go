package repos

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/kjedeligmann/spycat/internal/models"
)

type MissionRepo struct {
	db *sql.DB
}

func NewMissionRepo(db *sql.DB) *MissionRepo {
	return &MissionRepo{db: db}
}

// Create a mission along with its targets.
func (r *MissionRepo) Create(ctx context.Context, mission *models.Mission) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var query string
	query = `INSERT INTO missions (cat_id, status) VALUES ($1, $2) RETURNING id`
	err = tx.QueryRowContext(ctx, query, &mission.CatID, mission.Status).Scan(&mission.ID)

	if err != nil {
		return err
	}

	// Insert targets if provided
	for _, target := range mission.Targets {
		query := `INSERT INTO targets (mission_id, name, country, notes, completed) VALUES ($1, $2, $3, $4, $5)`
		_, err := tx.ExecContext(ctx, query, mission.ID, target.Name, target.Country, target.Notes, target.Completed)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// Get retrieves a Mission by its ID.
func (r *MissionRepo) Get(ctx context.Context, id int) (*models.Mission, error) {
	var mission models.Mission

	query := `
        SELECT id, cat_id, status
        FROM missions
        WHERE id = $1
    `
	row := r.db.QueryRowContext(ctx, query, id)
	err := row.Scan(&mission.ID, &mission.CatID, &mission.Status)
	fmt.Println(mission)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	// Get targets
	query = `
        SELECT id, name, country, notes, completed
        FROM targets
        WHERE mission_id = $1
    `
	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var target models.Target
		if err := rows.Scan(&target.ID, &target.Name, &target.Country, &target.Notes, &target.Completed); err != nil {
			return nil, err
		}
		mission.Targets = append(mission.Targets, target)
	}

	return &mission, nil
}

// Update modifies an existing Mission (for example, to update its status).
func (r *MissionRepo) Update(ctx context.Context, mission *models.Mission) error {
	query := `
        UPDATE missions
        SET cat_id = $1, status = $2
        WHERE id = $3
    `
	_, err := r.db.ExecContext(ctx, query, mission.CatID, mission.Status, mission.ID)
	return err
}

// List retrieves all missions from the database.
func (r *MissionRepo) List(ctx context.Context) ([]models.Mission, error) {
	query := `
        SELECT id, cat_id, status
        FROM missions
    `
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var missions []models.Mission
	for rows.Next() {
		var m models.Mission
		if err := rows.Scan(&m.ID, &m.CatID, &m.Status); err != nil {
			return nil, err
		}
		missions = append(missions, m)
	}
	return missions, nil
}

// Assign spy cat to a mission
func (r *MissionRepo) AssignSpyCat(ctx context.Context, missionID, catID int) error {
	query := `UPDATE missions SET cat_id = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, catID, missionID)
	return err
}

// Mark mission as completed
func (r *MissionRepo) MarkMissionCompleted(ctx context.Context, id int) error {
	query := `UPDATE missions SET status = 'completed' WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// Delete a mission (only if not assigned to a cat)
func (r *MissionRepo) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM missions WHERE id = $1 AND cat_id IS NULL`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("mission cannot be deleted")
	}
	return nil
}

// Add a new target to a mission if the mission is not completed and target count is less than 3
func (r *MissionRepo) AddTarget(ctx context.Context, missionID int, target *models.Target) error {
	// Check mission status
	var status string
	if err := r.db.QueryRowContext(ctx, `SELECT status FROM missions WHERE id = $1`, missionID).Scan(&status); err != nil {
		return err
	}
	if status == "completed" {
		return errors.New("cannot add target to a completed mission")
	}

	// Check the number of targets already in the mission
	var count int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM targets WHERE mission_id = $1`, missionID).Scan(&count); err != nil {
		return err
	}
	if count >= 3 {
		return errors.New("mission already has the maximum number of targets")
	}

	query := `
        INSERT INTO targets (mission_id, name, country, notes, completed)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id
    `
	if err := r.db.QueryRowContext(ctx, query,
		missionID,
		target.Name,
		target.Country,
		target.Notes,
		target.Completed,
	).Scan(&target.ID); err != nil {
		return err
	}
	target.MissionID = missionID
	return nil
}
