-- Create indexes
CREATE INDEX idx_user_id ON users(id);
CREATE INDEX idx_user_email ON users(email);

CREATE INDEX idx_activity_id ON activities(id);
CREATE INDEX idx_activity_user_id ON activities(user_id);
CREATE INDEX idx_activities_user_time ON activities(user_id, done_at);
CREATE INDEX idx_activities_user_calories ON activities(user_id, calories_burned);