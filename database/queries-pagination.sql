-- ============================================
-- PAGINATION QUERIES - Truy vấn có phân trang
-- Warehouse Management System
-- Sử dụng cho Frontend API
-- ============================================

-- Pattern: LIMIT = số record mỗi trang, OFFSET = (page - 1) * limit
-- Ví dụ: Page 1 -> OFFSET 0, Page 2 -> OFFSET 20, Page 3 -> OFFSET 40...

-- ============================================
-- 1. PRODUCTS WITH PAGINATION
-- ============================================

-- Lấy danh sách sản phẩm có phân trang
-- Params: @limit (default 20), @offset (default 0)
WITH product_list AS (
    SELECT 
        p.id,
        p.code,
        p.name,
        p.description,
        p.unit,
        p.category,
        p.min_stock,
        COALESCE(SUM(i.quantity), 0) as total_stock,
        p.is_active,
        p.created_at
    FROM products p
    LEFT JOIN inventory i ON p.id = i.product_id
    WHERE p.is_active = true
    GROUP BY p.id, p.code, p.name, p.description, p.unit, p.category, p.min_stock, p.is_active, p.created_at
)
SELECT 
    *,
    (SELECT COUNT(*) FROM product_list) as total_count
FROM product_list
ORDER BY code
LIMIT 20 OFFSET 0;  -- Thay đổi LIMIT và OFFSET theo params từ FE

-- ============================================
-- 2. INVENTORY WITH PAGINATION
-- ============================================

-- Lấy danh sách tồn kho có phân trang
-- Params: @warehouse_id (optional), @limit, @offset
WITH inventory_list AS (
    SELECT 
        i.id,
        w.id as warehouse_id,
        w.name as warehouse_name,
        w.code as warehouse_code,
        p.id as product_id,
        p.code as product_code,
        p.name as product_name,
        p.category,
        i.quantity,
        p.unit,
        p.min_stock,
        CASE 
            WHEN i.quantity < p.min_stock THEN 'LOW'
            WHEN i.quantity >= p.min_stock * 2 THEN 'HIGH'
            ELSE 'OK'
        END as status,
        i.last_updated
    FROM inventory i
    JOIN warehouses w ON i.warehouse_id = w.id
    JOIN products p ON i.product_id = p.id
    WHERE w.is_active = true
    -- AND w.id = @warehouse_id  -- Bỏ comment nếu filter theo kho
)
SELECT 
    *,
    (SELECT COUNT(*) FROM inventory_list) as total_count
FROM inventory_list
ORDER BY warehouse_name, product_code
LIMIT 20 OFFSET 0;

-- ============================================
-- 3. TRANSACTIONS WITH PAGINATION
-- ============================================

-- Lấy danh sách giao dịch có phân trang
-- Params: @type (optional), @warehouse_id (optional), @start_date, @end_date, @limit, @offset
WITH transaction_list AS (
    SELECT 
        t.id,
        t.transaction_code,
        t.type,
        t.status,
        t.warehouse_id,
        w.name as warehouse_name,
        w.code as warehouse_code,
        t.user_id,
        u.full_name as user_name,
        u.username,
        t.note,
        t.transaction_date,
        t.created_at,
        (
            SELECT COUNT(*)
            FROM transaction_details td
            WHERE td.transaction_id = t.id
        ) as items_count,
        (
            SELECT SUM(td.quantity * td.unit_price)
            FROM transaction_details td
            WHERE td.transaction_id = t.id
        ) as total_value
    FROM transactions t
    JOIN warehouses w ON t.warehouse_id = w.id
    JOIN users u ON t.user_id = u.id
    WHERE 1=1
    -- AND t.type = @type  -- Filter theo loại (IMPORT/EXPORT)
    -- AND t.warehouse_id = @warehouse_id  -- Filter theo kho
    -- AND t.transaction_date >= @start_date
    -- AND t.transaction_date <= @end_date
)
SELECT 
    *,
    (SELECT COUNT(*) FROM transaction_list) as total_count
FROM transaction_list
ORDER BY transaction_date DESC, id DESC
LIMIT 20 OFFSET 0;

-- ============================================
-- 4. TRANSACTION DETAILS WITH PAGINATION
-- ============================================

-- Lấy chi tiết giao dịch theo transaction_id
-- Params: @transaction_id, @limit, @offset
WITH detail_list AS (
    SELECT 
        td.id,
        td.transaction_id,
        td.product_id,
        p.code as product_code,
        p.name as product_name,
        p.unit,
        p.category,
        td.quantity,
        td.unit_price,
        (td.quantity * td.unit_price) as total_price,
        td.note,
        td.created_at
    FROM transaction_details td
    JOIN products p ON td.product_id = p.id
    WHERE td.transaction_id = 1  -- Thay bằng @transaction_id
)
SELECT 
    *,
    (SELECT COUNT(*) FROM detail_list) as total_count
FROM detail_list
ORDER BY product_code
LIMIT 20 OFFSET 0;

-- ============================================
-- 5. USERS WITH PAGINATION
-- ============================================

-- Lấy danh sách users có phân trang
-- Params: @role_id (optional), @is_active (optional), @limit, @offset
WITH user_list AS (
    SELECT 
        u.id,
        u.username,
        u.email,
        u.full_name,
        u.phone,
        u.role_id,
        r.name as role_name,
        u.is_active,
        u.last_login,
        u.created_at,
        (
            SELECT COUNT(DISTINCT uwp.warehouse_id)
            FROM user_warehouse_permissions uwp
            WHERE uwp.user_id = u.id
        ) as assigned_warehouses
    FROM users u
    JOIN roles r ON u.role_id = r.id
    WHERE 1=1
    -- AND u.role_id = @role_id
    -- AND u.is_active = @is_active
)
SELECT 
    *,
    (SELECT COUNT(*) FROM user_list) as total_count
FROM user_list
ORDER BY created_at DESC, id DESC
LIMIT 20 OFFSET 0;

-- ============================================
-- 6. SEARCH PRODUCTS (Full-text search với pagination)
-- ============================================

-- Tìm kiếm sản phẩm theo keyword
-- Params: @keyword, @category (optional), @limit, @offset
WITH search_results AS (
    SELECT 
        p.id,
        p.code,
        p.name,
        p.description,
        p.unit,
        p.category,
        p.min_stock,
        COALESCE(SUM(i.quantity), 0) as total_stock,
        p.is_active
    FROM products p
    LEFT JOIN inventory i ON p.id = i.product_id
    WHERE p.is_active = true
        AND (
            LOWER(p.name) LIKE LOWER('%laptop%')  -- Thay bằng @keyword
            OR LOWER(p.code) LIKE LOWER('%laptop%')
            OR LOWER(p.description) LIKE LOWER('%laptop%')
        )
        -- AND p.category = @category
    GROUP BY p.id
)
SELECT 
    *,
    (SELECT COUNT(*) FROM search_results) as total_count
FROM search_results
ORDER BY name
LIMIT 20 OFFSET 0;

-- ============================================
-- 7. WAREHOUSE INVENTORY REPORT (với pagination)
-- ============================================

-- Báo cáo tồn kho theo kho với phân trang
-- Params: @warehouse_id, @category (optional), @limit, @offset
WITH warehouse_inventory AS (
    SELECT 
        w.id as warehouse_id,
        w.name as warehouse_name,
        p.id as product_id,
        p.code as product_code,
        p.name as product_name,
        p.category,
        COALESCE(i.quantity, 0) as quantity,
        p.unit,
        p.min_stock,
        CASE 
            WHEN COALESCE(i.quantity, 0) < p.min_stock THEN 'LOW'
            WHEN COALESCE(i.quantity, 0) >= p.min_stock * 2 THEN 'HIGH'
            WHEN COALESCE(i.quantity, 0) = 0 THEN 'OUT_OF_STOCK'
            ELSE 'OK'
        END as stock_status,
        i.last_updated
    FROM warehouses w
    CROSS JOIN products p
    LEFT JOIN inventory i ON w.id = i.warehouse_id AND p.id = i.product_id
    WHERE w.is_active = true 
        AND p.is_active = true
        AND w.id = 1  -- Thay bằng @warehouse_id
        -- AND p.category = @category
)
SELECT 
    *,
    (SELECT COUNT(*) FROM warehouse_inventory) as total_count
FROM warehouse_inventory
ORDER BY 
    CASE stock_status
        WHEN 'OUT_OF_STOCK' THEN 1
        WHEN 'LOW' THEN 2
        WHEN 'OK' THEN 3
        WHEN 'HIGH' THEN 4
    END,
    product_code
LIMIT 20 OFFSET 0;

-- ============================================
-- 8. TRANSACTION HISTORY BY PRODUCT (với pagination)
-- ============================================

-- Lịch sử giao dịch của một sản phẩm
-- Params: @product_id, @type (optional), @limit, @offset
WITH product_history AS (
    SELECT 
        t.id,
        t.transaction_code,
        t.type,
        t.transaction_date,
        w.name as warehouse_name,
        td.quantity,
        p.unit,
        td.unit_price,
        (td.quantity * td.unit_price) as total_value,
        u.full_name as performed_by,
        t.note
    FROM transaction_details td
    JOIN transactions t ON td.transaction_id = t.id
    JOIN products p ON td.product_id = p.id
    JOIN warehouses w ON t.warehouse_id = w.id
    JOIN users u ON t.user_id = u.id
    WHERE td.product_id = 1  -- Thay bằng @product_id
        AND t.status = 'COMPLETED'
        -- AND t.type = @type
)
SELECT 
    *,
    (SELECT COUNT(*) FROM product_history) as total_count
FROM product_history
ORDER BY transaction_date DESC
LIMIT 20 OFFSET 0;

-- ============================================
-- 9. LOW STOCK ALERT (với pagination)
-- ============================================

-- Cảnh báo sản phẩm sắp hết hàng
-- Params: @warehouse_id (optional), @limit, @offset
WITH low_stock_items AS (
    SELECT 
        w.id as warehouse_id,
        w.name as warehouse_name,
        p.id as product_id,
        p.code as product_code,
        p.name as product_name,
        p.category,
        COALESCE(i.quantity, 0) as current_stock,
        p.min_stock,
        p.unit,
        (p.min_stock - COALESCE(i.quantity, 0)) as shortage,
        CASE 
            WHEN COALESCE(i.quantity, 0) = 0 THEN 'CRITICAL'
            WHEN COALESCE(i.quantity, 0) < p.min_stock * 0.5 THEN 'URGENT'
            ELSE 'WARNING'
        END as alert_level
    FROM products p
    CROSS JOIN warehouses w
    LEFT JOIN inventory i ON p.id = i.product_id AND w.id = i.warehouse_id
    WHERE p.is_active = true 
        AND w.is_active = true
        AND COALESCE(i.quantity, 0) < p.min_stock
        -- AND w.id = @warehouse_id
)
SELECT 
    *,
    (SELECT COUNT(*) FROM low_stock_items) as total_count
FROM low_stock_items
ORDER BY 
    CASE alert_level
        WHEN 'CRITICAL' THEN 1
        WHEN 'URGENT' THEN 2
        WHEN 'WARNING' THEN 3
    END,
    shortage DESC
LIMIT 20 OFFSET 0;

-- ============================================
-- 10. HELPER: Tính số trang
-- ============================================

-- Function để tính tổng số trang
-- Example usage:
-- total_count = 150, page_size = 20 -> total_pages = 8
-- total_pages = CEILING(total_count::DECIMAL / page_size)

-- Query metadata cho pagination
SELECT 
    COUNT(*) as total_records,
    20 as page_size,
    CEILING(COUNT(*)::DECIMAL / 20) as total_pages,
    1 as current_page,  -- Từ params
    CASE 
        WHEN 1 > 1 THEN true 
        ELSE false 
    END as has_previous,
    CASE 
        WHEN 1 < CEILING(COUNT(*)::DECIMAL / 20) THEN true 
        ELSE false 
    END as has_next
FROM products;
