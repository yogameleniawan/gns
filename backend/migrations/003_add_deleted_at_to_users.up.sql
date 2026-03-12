-- Add deleted_at column to users table for soft delete support
ALTER TABLE users ADD COLUMN deleted_at TIMESTAMP;

-- Create index on deleted_at for better query performance
CREATE INDEX idx_users_deleted_at ON users(deleted_at);
