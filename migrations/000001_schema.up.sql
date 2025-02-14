-- IF NOT EXISTS ensures idempotent migrations
-- Create the spy_cats table first.
CREATE TABLE IF NOT EXISTS spy_cats (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    years_experience INT CHECK (years_experience >= 0),
    breed TEXT NOT NULL,
    salary DECIMAL(10,2) CHECK (salary >= 0)
);

-- Create the missions table.
CREATE TABLE IF NOT EXISTS missions (
    id SERIAL PRIMARY KEY,
    status TEXT CHECK (status IN ('ongoing', 'completed')) DEFAULT 'ongoing',
    cat_id INT UNIQUE
      REFERENCES spy_cats(id)
      ON DELETE SET NULL
);

-- One cat can have only one mission
ALTER TABLE missions ADD CONSTRAINT unique_cat_mission UNIQUE (cat_id);

-- Create the targets table.
CREATE TABLE IF NOT EXISTS targets (
    id SERIAL PRIMARY KEY,
    mission_id INT NOT NULL,
    name TEXT NOT NULL,
    country TEXT NOT NULL,
    notes TEXT DEFAULT '',
    completed BOOLEAN DEFAULT FALSE,
    CONSTRAINT fk_mission
      FOREIGN KEY (mission_id)
      REFERENCES missions(id)
      ON DELETE CASCADE,
    CONSTRAINT uq_target_mission_name UNIQUE (mission_id, name)
);

-- Trigger function to enforce max 3 targets per mission
CREATE OR REPLACE FUNCTION check_target_count() RETURNS TRIGGER AS $$
DECLARE target_count INT;
BEGIN
    SELECT COUNT(*) INTO target_count FROM targets WHERE mission_id = NEW.mission_id;

    IF target_count >= 3 THEN
        RAISE EXCEPTION 'A mission cannot have more than 3 targets';
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to enforce the target limit
CREATE TRIGGER enforce_target_limit
BEFORE INSERT ON targets
FOR EACH ROW
EXECUTE FUNCTION check_target_count();

