-- Create enum
CREATE TYPE enum_activity_types as ENUM (
    'Walking',
    'Yoga',
    'Stretching',
    'Cycling',
    'Swimming',
    'Dancing',
    'Hiking',
    'Running',
    'HIIT',
    'JumpRope'
);

-- Create table activities
CREATE TABLE activities (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    activity_type enum_activity_types NOT NULL,
    done_at TIMESTAMPTZ NOT NULL,
    duration_in_minutes INT NOT NULL,                   
    calories_burned INT NOT NULL, 
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create updated_at trigger function
CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create triggers
CREATE TRIGGER set_timestamp_activities
    BEFORE UPDATE ON activities
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_timestamp();