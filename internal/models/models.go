package models

import "database/sql"

type SpyCat struct {
	ID              int     `json:"id"`
	Name            string  `json:"name"`
	YearsExperience int     `json:"years_experience"`
	Breed           string  `json:"breed"`
	Salary          float64 `json:"salary"`
	MissionID       int     `json:"mission_id"`
}

type Mission struct {
	ID     int           `json:"id"`
	CatID  sql.NullInt64 `json:"cat_id"`
	Status string        `json:"status"` // e.g., "ongoing" or "completed"

	Targets []Target `json:"targets"`
}

type Target struct {
	ID        int    `json:"id"`
	MissionID int    `json:"mission_id"`
	Name      string `json:"name"`
	Country   string `json:"country"`
	Notes     string `json:"notes"`
	Completed bool   `json:"completed"`
}
