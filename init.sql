-- Remove PostgreSQL specific extension creation
-- CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";

CREATE TABLE users (
    user_id VARCHAR(36) PRIMARY KEY,
    user_name VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    pass TEXT NOT NULL,
    role TEXT NOT NULL CHECK (role IN ('admin', 'manager', 'user')),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE data (
    data_id VARCHAR(36) PRIMARY KEY,
    company_data JSON NOT NULL
);
