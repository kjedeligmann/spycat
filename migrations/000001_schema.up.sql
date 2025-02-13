-- IF NOT EXISTS ensures idempotent migrations
CREATE TABLE IF NOT EXISTS missions (
    id SERIAL PRIMARY KEY,
    status TEXT CHECK (status IN ('ongoing', 'completed')) DEFAULT 'ongoing'
);

CREATE TABLE IF NOT EXISTS spy_cats (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    years_experience INT CHECK (years_experience >= 0),
    breed TEXT NOT NULL,
    salary DECIMAL(10,2) CHECK (salary >= 0),
    mission_id INT REFERENCES missions(id) ON DELETE SET NULL
);

-- Add cat_id column to missions and then create the foreign key constraint
ALTER TABLE missions ADD COLUMN cat_id INT UNIQUE;
ALTER TABLE missions ADD CONSTRAINT fk_cat_id FOREIGN KEY (cat_id) REFERENCES spy_cats(id) ON DELETE SET NULL;

CREATE TABLE IF NOT EXISTS targets (
    id SERIAL PRIMARY KEY,
    mission_id INT REFERENCES missions(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    country TEXT NOT NULL,
    notes TEXT DEFAULT '',
    completed BOOLEAN DEFAULT FALSE,
    UNIQUE (mission_id, name)  -- Ensures target names are unique within a mission
);

