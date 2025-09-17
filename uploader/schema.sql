-- UC Wages Database Schema

-- Create table for wage records
CREATE TABLE IF NOT EXISTS uc_wages (
    id SERIAL PRIMARY KEY,
    location VARCHAR(50) NOT NULL,
    year INTEGER NOT NULL,
    employee_id INTEGER,
    firstname VARCHAR(100),
    lastname VARCHAR(100),
    title VARCHAR(200),
    basepay DECIMAL(12,2) DEFAULT 0.00,
    overtimepay DECIMAL(12,2) DEFAULT 0.00,
    adjustpay DECIMAL(12,2) DEFAULT 0.00,
    grosspay DECIMAL(12,2) DEFAULT 0.00,
    scraped_at TIMESTAMP,
    uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(location, year, employee_id)
);

-- Create indexes for faster queries
CREATE INDEX IF NOT EXISTS idx_uc_wages_location ON uc_wages(location);
CREATE INDEX IF NOT EXISTS idx_uc_wages_year ON uc_wages(year);
CREATE INDEX IF NOT EXISTS idx_uc_wages_location_year ON uc_wages(location, year);
CREATE INDEX IF NOT EXISTS idx_uc_wages_grosspay ON uc_wages(grosspay);
CREATE INDEX IF NOT EXISTS idx_uc_wages_title ON uc_wages(title);

-- Create table for upload tracking
CREATE TABLE IF NOT EXISTS upload_progress (
    id SERIAL PRIMARY KEY,
    location VARCHAR(50) NOT NULL,
    year INTEGER NOT NULL,
    total_records INTEGER,
    uploaded_records INTEGER DEFAULT 0,
    status VARCHAR(20) DEFAULT 'pending', -- pending, processing, completed, failed
    started_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP,
    error_message TEXT,
    UNIQUE(location, year)
);