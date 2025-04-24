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

CREATE TABLE account_data_map (
    account_id VARCHAR(36) NOT NULL,
    data_id VARCHAR(36) NOT NULL
);

CREATE TABLE company_data (
    id VARCHAR(255) PRIMARY KEY,
    company_name VARCHAR(255) NOT NULL,
    company_address TEXT,
    drive VARCHAR(255),
    type_of_drive VARCHAR(100),
    follow_up TEXT,
    is_contacted BOOLEAN DEFAULT FALSE,
    remarks TEXT,
    contact_details TEXT,
    hr_details TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
