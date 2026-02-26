-- Drop indexes
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_incidents_created_at;
DROP INDEX IF EXISTS idx_incidents_created_by;
DROP INDEX IF EXISTS idx_incidents_assigned_to;
DROP INDEX IF EXISTS idx_incidents_severity;
DROP INDEX IF EXISTS idx_incidents_status;

-- Drop tables
DROP TABLE IF EXISTS incidents;
DROP TABLE IF EXISTS users;
