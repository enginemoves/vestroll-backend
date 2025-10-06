-- VestRoll Database Setup Script
-- Run this script to create the necessary database and tables

-- Create database (run as postgres superuser)
-- CREATE DATABASE vestroll;

-- Connect to vestroll database and run the following:

-- Create users table for user authentication
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at);

-- Insert a test user (optional - for testing)
-- Password is 'TestPassword123!' hashed with bcrypt
-- INSERT INTO users (email, password, full_name) VALUES 
-- ('test@vestroll.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'Test User');

-- Verify the table was created
SELECT 'Users table created successfully' as status;
SELECT COUNT(*) as user_count FROM users;
