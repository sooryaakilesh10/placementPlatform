-- Remove PostgreSQL specific extension creation
-- CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";

CREATE TABLE users (
    user_id VARCHAR(36) PRIMARY KEY,
    user_name VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    pass TEXT NOT NULL,
    role TEXT NOT NULL CHECK (role IN ('admin', 'manager', 'placement_officer')),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE system_settings (
    id INT PRIMARY KEY AUTO_INCREMENT,
    approval_mode VARCHAR(10) NOT NULL DEFAULT 'manual' CHECK (approval_mode IN ('auto', 'manual')),
    last_updated DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_by VARCHAR(36),
    FOREIGN KEY (updated_by) REFERENCES users(user_id)
);

CREATE TABLE companies (
    company_id VARCHAR(36) PRIMARY KEY,
    company_name VARCHAR(255) NOT NULL,
    last_contacted BOOLEAN DEFAULT FALSE,
    follow_up DATETIME,
    packages JSON, -- Array of package values
    remarks TEXT,
    target_branch VARCHAR(100),
    is_validation BOOLEAN DEFAULT FALSE,
    approved BOOLEAN DEFAULT FALSE,
    location VARCHAR(255),
    hr_name VARCHAR(100),
    hr_email VARCHAR(255),
    hr_phone VARCHAR(20),
    hr_position VARCHAR(100),
    hr_linkedin_url VARCHAR(255),
    website VARCHAR(255),
    industry VARCHAR(100),
    founded_year INT,
    company_size VARCHAR(50),
    description TEXT,
    logo_url VARCHAR(255),
    assigned_to VARCHAR(36), -- ID of the placement officer
    created_by VARCHAR(36) NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    is_data_validated BOOLEAN DEFAULT FALSE,
    approval_status VARCHAR(20) DEFAULT 'PENDING' CHECK (approval_status IN ('PENDING', 'APPROVED', 'REJECTED')),
    approval_notes TEXT,
    FOREIGN KEY (assigned_to) REFERENCES users(user_id),
    FOREIGN KEY (created_by) REFERENCES users(user_id)
);

CREATE TABLE recruitment_drives (
    drive_id VARCHAR(36) PRIMARY KEY,
    company_id VARCHAR(36) NOT NULL,
    status VARCHAR(20) NOT NULL CHECK (status IN ('SCHEDULED', 'COMPLETED', 'NO_HIRING')),
    scheduled_date DATETIME,
    number_of_offers INT DEFAULT 0,
    number_hired INT DEFAULT 0,
    roles_offered JSON, -- Array of roles
    min_cgpa FLOAT,
    eligible_branches JSON, -- Array of eligible branches
    notes TEXT,
    created_by VARCHAR(36) NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (company_id) REFERENCES companies(company_id),
    FOREIGN KEY (created_by) REFERENCES users(user_id)
);

CREATE TABLE company_assignments (
    assignment_id VARCHAR(36) PRIMARY KEY,
    company_id VARCHAR(36) NOT NULL,
    officer_id VARCHAR(36) NOT NULL,
    assigned_by VARCHAR(36) NOT NULL,
    assigned_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    FOREIGN KEY (company_id) REFERENCES companies(company_id),
    FOREIGN KEY (officer_id) REFERENCES users(user_id),
    FOREIGN KEY (assigned_by) REFERENCES users(user_id)
);

-- Initialize system settings
INSERT INTO system_settings (approval_mode) VALUES ('manual');

-- Create initial admin user (password: admin123)
INSERT INTO users (user_id, user_name, email, pass, role) 
VALUES ('9c0f7891-56c8-4621-aec8-da80b16fed6e', 'Admin', 'admin@system.com', 'admin123', 'admin');
