-- ============================================
-- Sales Module Database Schema
-- Database: salesdb (PostgreSQL)
-- ============================================

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ============================================
-- USERS TABLE (for JWT authentication)
-- ============================================
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) DEFAULT 'admin',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ============================================
-- CUSTOMERS TABLE
-- ============================================
CREATE TABLE IF NOT EXISTS customers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100),
    phone VARCHAR(20),
    address TEXT,
    company VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ============================================
-- SALES TABLE
-- ============================================
CREATE TABLE IF NOT EXISTS sales (
    id SERIAL PRIMARY KEY,
    customer_id INTEGER REFERENCES customers(id) ON DELETE SET NULL,
    order_id VARCHAR(50) UNIQUE,
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'completed', 'cancelled')),
    subtotal DECIMAL(12,2) DEFAULT 0.00,
    discount_percent DECIMAL(5,2) DEFAULT 0.00,
    discount_amount DECIMAL(12,2) DEFAULT 0.00,
    tax_amount DECIMAL(12,2) DEFAULT 0.00,
    total_amount DECIMAL(12,2) DEFAULT 0.00,
    notes TEXT,
    created_by INTEGER REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ============================================
-- SALE ITEMS TABLE
-- ============================================
CREATE TABLE IF NOT EXISTS sale_items (
    id SERIAL PRIMARY KEY,
    sale_id INTEGER NOT NULL REFERENCES sales(id) ON DELETE CASCADE,
    quantity INTEGER NOT NULL DEFAULT 1,
    unit_price DECIMAL(12,2) NOT NULL DEFAULT 0.00,
    line_total DECIMAL(12,2) NOT NULL DEFAULT 0.00
);

-- ============================================
-- INDEXES
-- ============================================
CREATE INDEX IF NOT EXISTS idx_sales_customer_id ON sales(customer_id);
CREATE INDEX IF NOT EXISTS idx_sales_status ON sales(status);
CREATE INDEX IF NOT EXISTS idx_sales_created_at ON sales(created_at);
CREATE INDEX IF NOT EXISTS idx_sales_order_id ON sales(order_id);
CREATE INDEX IF NOT EXISTS idx_sale_items_sale_id ON sale_items(sale_id);

-- ============================================
-- SEED DATA
-- ============================================

-- Admin user (password: admin123)
-- bcrypt hash for "admin123"
INSERT INTO users (username, email, password_hash, role) VALUES
('admin', 'admin@salesmodule.com', '$2a$10$voPywEBD.A5WThPii2FuM.ng49rNhJqgn4BYlke8hSclsB87RreIa', 'admin'),
('manager', 'manager@salesmodule.com', '$2a$10$voPywEBD.A5WThPii2FuM.ng49rNhJqgn4BYlke8hSclsB87RreIa', 'admin')
ON CONFLICT (username) DO NOTHING;
