-- Create users table for authentication
-- Migration: 001_create_users_table.sql
-- Description: Initial users table with email/password authentication

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create index for faster email lookups
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- Insert a test user with password "testpassword123" (bcrypt hashed)
-- Password hash for "testpassword123": $2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewdBPj5B5wuG7S.O
INSERT INTO users (email, password) VALUES 
('test@vestroll.com', '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewdBPj5B5wuG7S.O')
ON CONFLICT (email) DO NOTHING;