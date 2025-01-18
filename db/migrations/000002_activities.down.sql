-- DROP trigger
DROP TRIGGER IF EXISTS set_timestamp_activities ON activities CASCADE;
DROP FUNCTION IF EXISTS trigger_set_timestamp CASCADE;

-- DROP users
DROP TABLE IF EXISTS activities CASCADE;

-- DROP enum
DROP TYPE IF EXISTS enum_activity_types CASCADE;
