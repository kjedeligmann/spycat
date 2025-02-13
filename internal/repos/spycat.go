package repos

import (
	"context"
	"database/sql"
	"github.com/kjedeligmann/spycat/internal/models"
)

type SpyCatRepo struct {
	db *sql.DB
}

func NewSpyCatRepo(db *sql.DB) *SpyCatRepo {
	return &SpyCatRepo{db: db}
}

// Create inserts a new SpyCat record into the DB
func (r *SpyCatRepo) Create(ctx context.Context, cat *models.SpyCat) error {
	query := `
        INSERT INTO spy_cats (name, years_experience, breed, salary)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `
	return r.db.QueryRowContext(
		ctx,
		query,
		cat.Name,
		cat.YearsExperience,
		cat.Breed,
		cat.Salary,
	).Scan(&cat.ID)
}

// Read retrieves a SpyCat by ID
func (r *SpyCatRepo) Read(ctx context.Context, id int) (*models.SpyCat, error) {
	query := `
        SELECT id, name, years_experience, breed, salary
        FROM spy_cats
        WHERE id = $1
    `
	row := r.db.QueryRowContext(ctx, query, id)

	var cat models.SpyCat
	if err := row.Scan(&cat.ID, &cat.Name, &cat.YearsExperience, &cat.Breed, &cat.Salary); err != nil {
		return nil, err
	}
	return &cat, nil
}

// List retrieves all SpyCat records from the database.
func (r *SpyCatRepo) List(ctx context.Context) ([]models.SpyCat, error) {
	query := `SELECT id, name, years_experience, breed, salary FROM spy_cats`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cats []models.SpyCat
	for rows.Next() {
		var cat models.SpyCat
		if err := rows.Scan(&cat.ID, &cat.Name, &cat.YearsExperience, &cat.Breed, &cat.Salary); err != nil {
			return nil, err
		}
		cats = append(cats, cat)
	}
	return cats, nil
}

// UpdateSalary modifies an existing SpyCats salary
func (r *SpyCatRepo) UpdateSalary(ctx context.Context, catID int, newSalary float64) error {
	query := `
        UPDATE spy_cats
        SET salary = $1
        WHERE id = $2
    `
	_, err := r.db.ExecContext(ctx, query, newSalary, catID)
	return err
}

// Delete removes a SpyCat by ID
func (r *SpyCatRepo) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM spy_cats WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
