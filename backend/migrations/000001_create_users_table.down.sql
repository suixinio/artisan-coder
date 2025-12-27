DROP TRIGGER update_users_updated_at ON users;
DROP FUNCTION update_updated_at_column();
DROP INDEX IF EXISTS idx_users_created_at;
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_username;
DROP TABLE IF EXISTS users;
