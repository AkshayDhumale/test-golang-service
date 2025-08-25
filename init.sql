-- init.sql - Database initialization script

-- Create the users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert some sample data
INSERT INTO users (name, email) VALUES 
    ('John Doe', 'john.doe@example.com'),
    ('Jane Smith', 'jane.smith@example.com'),
    ('Bob Johnson', 'bob.johnson@example.com')
ON CONFLICT (email) DO NOTHING;

-- Create an index on email for better performance
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- Create an index on created_at for time-based queries
CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at);
