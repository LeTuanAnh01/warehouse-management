-- ============================================
-- WAREHOUSE MANAGEMENT SYSTEM - DATABASE SCHEMA
-- PostgreSQL Database
-- ============================================

-- Drop tables if exists (for clean setup)
DROP TABLE IF EXISTS transaction_details CASCADE;
DROP TABLE IF EXISTS transactions CASCADE;
DROP TABLE IF EXISTS inventory CASCADE;
DROP TABLE IF EXISTS user_warehouse_permissions CASCADE;
DROP TABLE IF EXISTS products CASCADE;
DROP TABLE IF EXISTS warehouses CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS roles CASCADE;

-- ============================================
-- 1. ROLES TABLE - Vai trò hệ thống
-- ============================================
CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ============================================
-- 2. USERS TABLE - Tài khoản người dùng
-- ============================================
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(100),
    phone VARCHAR(20),
    role_id INTEGER REFERENCES roles(id),
    is_active BOOLEAN DEFAULT true,
    last_login TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ============================================
-- 3. WAREHOUSES TABLE - Kho hàng
-- ============================================
CREATE TABLE warehouses (
    id SERIAL PRIMARY KEY,
    code VARCHAR(20) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    address TEXT,
    manager_id INTEGER REFERENCES users(id),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ============================================
-- 4. USER_WAREHOUSE_PERMISSIONS - Phân quyền theo kho
-- ============================================
CREATE TABLE user_warehouse_permissions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    warehouse_id INTEGER REFERENCES warehouses(id) ON DELETE CASCADE,
    can_view BOOLEAN DEFAULT true,
    can_import BOOLEAN DEFAULT false,
    can_export BOOLEAN DEFAULT false,
    can_manage BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, warehouse_id)
);

-- ============================================
-- 5. PRODUCTS TABLE - Sản phẩm
-- ============================================
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    unit VARCHAR(20),
    category VARCHAR(50),
    min_stock INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ============================================
-- 6. INVENTORY TABLE - Tồn kho
-- ============================================
CREATE TABLE inventory (
    id SERIAL PRIMARY KEY,
    warehouse_id INTEGER REFERENCES warehouses(id) ON DELETE CASCADE,
    product_id INTEGER REFERENCES products(id) ON DELETE CASCADE,
    quantity INTEGER DEFAULT 0 CHECK (quantity >= 0),
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(warehouse_id, product_id)
);

-- ============================================
-- 7. TRANSACTIONS TABLE - Giao dịch nhập/xuất
-- ============================================
CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    transaction_code VARCHAR(50) UNIQUE NOT NULL,
    warehouse_id INTEGER REFERENCES warehouses(id),
    type VARCHAR(20) NOT NULL CHECK (type IN ('IMPORT', 'EXPORT')),
    status VARCHAR(20) DEFAULT 'COMPLETED' CHECK (status IN ('PENDING', 'COMPLETED', 'CANCELLED')),
    user_id INTEGER REFERENCES users(id),
    note TEXT,
    transaction_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ============================================
-- 8. TRANSACTION_DETAILS TABLE - Chi tiết giao dịch
-- ============================================
CREATE TABLE transaction_details (
    id SERIAL PRIMARY KEY,
    transaction_id INTEGER REFERENCES transactions(id) ON DELETE CASCADE,
    product_id INTEGER REFERENCES products(id),
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    unit_price DECIMAL(15,2) DEFAULT 0,
    note TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ============================================
-- INDEXES for better performance
-- ============================================
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role_id ON users(role_id);
CREATE INDEX idx_warehouses_code ON warehouses(code);
CREATE INDEX idx_products_code ON products(code);
CREATE INDEX idx_inventory_warehouse ON inventory(warehouse_id);
CREATE INDEX idx_inventory_product ON inventory(product_id);
CREATE INDEX idx_transactions_warehouse ON transactions(warehouse_id);
CREATE INDEX idx_transactions_user ON transactions(user_id);
CREATE INDEX idx_transactions_type ON transactions(type);
CREATE INDEX idx_transactions_date ON transactions(transaction_date);
CREATE INDEX idx_transaction_details_transaction ON transaction_details(transaction_id);
CREATE INDEX idx_transaction_details_product ON transaction_details(product_id);

-- ============================================
-- TRIGGERS for auto-update timestamp
-- ============================================
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_warehouses_updated_at BEFORE UPDATE ON warehouses
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_products_updated_at BEFORE UPDATE ON products
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_transactions_updated_at BEFORE UPDATE ON transactions
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_roles_updated_at BEFORE UPDATE ON roles
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ============================================
-- TRIGGER: Auto update inventory on transaction
-- ============================================
CREATE OR REPLACE FUNCTION update_inventory_on_transaction()
RETURNS TRIGGER AS $$
DECLARE
    transaction_type VARCHAR(20);
    warehouse_id_val INTEGER;
BEGIN
    -- Get transaction info
    SELECT type, warehouse_id INTO transaction_type, warehouse_id_val
    FROM transactions WHERE id = NEW.transaction_id;

    -- Insert or update inventory
    IF transaction_type = 'IMPORT' THEN
        INSERT INTO inventory (warehouse_id, product_id, quantity, last_updated)
        VALUES (warehouse_id_val, NEW.product_id, NEW.quantity, CURRENT_TIMESTAMP)
        ON CONFLICT (warehouse_id, product_id)
        DO UPDATE SET 
            quantity = inventory.quantity + NEW.quantity,
            last_updated = CURRENT_TIMESTAMP;
    ELSIF transaction_type = 'EXPORT' THEN
        UPDATE inventory 
        SET quantity = quantity - NEW.quantity,
            last_updated = CURRENT_TIMESTAMP
        WHERE warehouse_id = warehouse_id_val 
          AND product_id = NEW.product_id;
    END IF;

    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER trigger_update_inventory 
    AFTER INSERT ON transaction_details
    FOR EACH ROW 
    EXECUTE FUNCTION update_inventory_on_transaction();

COMMENT ON TRIGGER trigger_update_inventory ON transaction_details IS 
'Automatically update inventory when transaction detail is inserted';