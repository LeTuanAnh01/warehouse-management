-- ============================================
-- SEED DATA - Dữ liệu mẫu cho hệ thống
-- ============================================

-- ============================================
-- 1. Insert Roles
-- ============================================
INSERT INTO roles (name, description) VALUES
('ADMIN', 'Quản trị viên hệ thống - Full quyền'),
('WAREHOUSE_MANAGER', 'Quản lý kho - Quản lý toàn bộ kho được phân quyền'),
('WAREHOUSE_STAFF', 'Nhân viên kho - Nhập/xuất hàng'),
('VIEWER', 'Người xem - Chỉ xem thông tin');

-- ============================================
-- 2. Insert Users
-- Password: "123456" hashed with bcrypt
-- ============================================
INSERT INTO users (username, email, password_hash, full_name, phone, role_id, is_active) VALUES
('admin', 'admin@warehouse.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'Administrator', '0901234567', 1, true),
('manager1', 'manager1@warehouse.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'Nguyễn Văn A', '0901234568', 2, true),
('manager2', 'manager2@warehouse.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'Trần Thị B', '0901234569', 2, true),
('staff1', 'staff1@warehouse.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'Lê Văn C', '0901234570', 3, true),
('staff2', 'staff2@warehouse.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'Phạm Thị D', '0901234571', 3, true),
('viewer1', 'viewer1@warehouse.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'Hoàng Văn E', '0901234572', 4, true);

-- ============================================
-- 3. Insert Warehouses
-- ============================================
INSERT INTO warehouses (code, name, address, manager_id, is_active) VALUES
('WH001', 'Kho Hà Nội', '123 Đường Láng, Đống Đa, Hà Nội', 2, true),
('WH002', 'Kho Hồ Chí Minh', '456 Nguyễn Huệ, Quận 1, TP.HCM', 3, true),
('WH003', 'Kho Đà Nẵng', '789 Trần Phú, Hải Châu, Đà Nẵng', 2, true);

-- ============================================
-- 4. Insert User Warehouse Permissions
-- ============================================
-- Admin có full quyền tất cả kho
INSERT INTO user_warehouse_permissions (user_id, warehouse_id, can_view, can_import, can_export, can_manage) VALUES
(1, 1, true, true, true, true),
(1, 2, true, true, true, true),
(1, 3, true, true, true, true);

-- Manager1 quản lý kho HN và Đà Nẵng
INSERT INTO user_warehouse_permissions (user_id, warehouse_id, can_view, can_import, can_export, can_manage) VALUES
(2, 1, true, true, true, true),
(2, 3, true, true, true, true);

-- Manager2 quản lý kho HCM
INSERT INTO user_warehouse_permissions (user_id, warehouse_id, can_view, can_import, can_export, can_manage) VALUES
(3, 2, true, true, true, true);

-- Staff1 làm việc tại kho HN (nhập/xuất)
INSERT INTO user_warehouse_permissions (user_id, warehouse_id, can_view, can_import, can_export, can_manage) VALUES
(4, 1, true, true, true, false);

-- Staff2 làm việc tại kho HCM (nhập/xuất)
INSERT INTO user_warehouse_permissions (user_id, warehouse_id, can_view, can_import, can_export, can_manage) VALUES
(5, 2, true, true, true, false);

-- Viewer1 chỉ xem kho HN
INSERT INTO user_warehouse_permissions (user_id, warehouse_id, can_view, can_import, can_export, can_manage) VALUES
(6, 1, true, false, false, false);

-- ============================================
-- 5. Insert Products
-- ============================================
INSERT INTO products (code, name, description, unit, category, min_stock, is_active) VALUES
('SP001', 'Laptop Dell XPS 15', 'Laptop cao cấp Dell XPS 15 inch', 'Cái', 'Electronics', 10, true),
('SP002', 'iPhone 15 Pro', 'Điện thoại iPhone 15 Pro 256GB', 'Cái', 'Electronics', 20, true),
('SP003', 'Samsung Galaxy S24', 'Điện thoại Samsung Galaxy S24', 'Cái', 'Electronics', 15, true),
('SP004', 'Bàn phím cơ Keychron K2', 'Bàn phím cơ không dây', 'Cái', 'Accessories', 30, true),
('SP005', 'Chuột Logitech MX Master 3', 'Chuột không dây cao cấp', 'Cái', 'Accessories', 25, true),
('SP006', 'Tai nghe Sony WH-1000XM5', 'Tai nghe chống ồn', 'Cái', 'Accessories', 20, true),
('SP007', 'Màn hình Dell U2723DE', 'Màn hình 27 inch 4K', 'Cái', 'Electronics', 15, true),
('SP008', 'Ổ cứng SSD Samsung 1TB', 'Ổ cứng SSD NVMe 1TB', 'Cái', 'Components', 50, true),
('SP009', 'RAM Corsair 32GB DDR5', 'RAM 32GB DDR5 6000MHz', 'Bộ', 'Components', 40, true),
('SP010', 'Card đồ họa RTX 4070', 'VGA NVIDIA RTX 4070 12GB', 'Cái', 'Components', 10, true);

-- ============================================
-- 6. Insert Initial Inventory
-- ============================================
-- Kho Hà Nội
INSERT INTO inventory (warehouse_id, product_id, quantity) VALUES
(1, 1, 15),
(1, 2, 25),
(1, 3, 20),
(1, 4, 35),
(1, 5, 30);

-- Kho HCM
INSERT INTO inventory (warehouse_id, product_id, quantity) VALUES
(2, 1, 10),
(2, 2, 30),
(2, 6, 25),
(2, 7, 18),
(2, 8, 60);

-- Kho Đà Nẵng
INSERT INTO inventory (warehouse_id, product_id, quantity) VALUES
(3, 3, 18),
(3, 4, 40),
(3, 9, 45),
(3, 10, 12);

-- ============================================
-- 7. Insert Sample Transactions (IMPORT)
-- ============================================
INSERT INTO transactions (transaction_code, warehouse_id, type, status, user_id, note, transaction_date) VALUES
('IMP-2024-001', 1, 'IMPORT', 'COMPLETED', 2, 'Nhập hàng đợt 1 tháng 1', '2024-01-15 09:00:00'),
('IMP-2024-002', 2, 'IMPORT', 'COMPLETED', 3, 'Nhập hàng đợt 1 tháng 1', '2024-01-16 10:30:00'),
('IMP-2024-003', 1, 'IMPORT', 'COMPLETED', 4, 'Nhập bổ sung tháng 2', '2024-02-10 14:00:00');

-- Transaction Details for IMPORT
INSERT INTO transaction_details (transaction_id, product_id, quantity, unit_price, note) VALUES
-- IMP-2024-001
(1, 1, 15, 25000000, 'Laptop Dell XPS 15'),
(1, 2, 25, 28000000, 'iPhone 15 Pro'),
(1, 3, 20, 22000000, 'Samsung S24'),
-- IMP-2024-002
(2, 1, 10, 25000000, 'Laptop Dell XPS 15'),
(2, 2, 30, 28000000, 'iPhone 15 Pro'),
(2, 6, 25, 8000000, 'Tai nghe Sony'),
-- IMP-2024-003
(3, 4, 35, 2500000, 'Bàn phím Keychron'),
(3, 5, 30, 2200000, 'Chuột Logitech');

-- ============================================
-- 8. Insert Sample Transactions (EXPORT)
-- ============================================
INSERT INTO transactions (transaction_code, warehouse_id, type, status, user_id, note, transaction_date) VALUES
('EXP-2024-001', 1, 'EXPORT', 'COMPLETED', 4, 'Xuất hàng cho đơn hàng #001', '2024-02-20 11:00:00'),
('EXP-2024-002', 2, 'EXPORT', 'COMPLETED', 5, 'Xuất hàng cho đơn hàng #002', '2024-02-22 15:30:00');

-- Transaction Details for EXPORT
-- Note: Inventory will be auto-updated by trigger
-- But since we already inserted initial inventory, we'll skip trigger execution for seed data

-- ============================================
-- SUMMARY
-- ============================================
SELECT 'Database seeded successfully!' as message;
SELECT 'Total Roles: ' || COUNT(*) as info FROM roles;
SELECT 'Total Users: ' || COUNT(*) as info FROM users;
SELECT 'Total Warehouses: ' || COUNT(*) as info FROM warehouses;
SELECT 'Total Products: ' || COUNT(*) as info FROM products;
SELECT 'Total Inventory Records: ' || COUNT(*) as info FROM inventory;
SELECT 'Total Transactions: ' || COUNT(*) as info FROM transactions;