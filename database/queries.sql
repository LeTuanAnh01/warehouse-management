-- ============================================
-- USEFUL QUERIES - Các truy vấn hữu ích
-- Warehouse Management System
-- ============================================

-- ============================================
-- 1. USER & PERMISSION QUERIES
-- ============================================

-- Xem tất cả users và roles
SELECT 
    u.id,
    u.username,
    u.email,
    u.full_name,
    r.name as role,
    u.is_active,
    u.last_login
FROM users u
JOIN roles r ON u.role_id = r.id
ORDER BY r.name, u.username;

-- Xem chi tiết quyền của một user cụ thể
SELECT 
    u.username,
    u.full_name,
    w.code as warehouse_code,
    w.name as warehouse_name,
    CASE WHEN uwp.can_view THEN '✓' ELSE '✗' END as view,
    CASE WHEN uwp.can_import THEN '✓' ELSE '✗' END as import,
    CASE WHEN uwp.can_export THEN '✓' ELSE '✗' END as export,
    CASE WHEN uwp.can_manage THEN '✓' ELSE '✗' END as manage
FROM user_warehouse_permissions uwp
JOIN users u ON uwp.user_id = u.id
JOIN warehouses w ON uwp.warehouse_id = w.id
WHERE u.username = 'staff1'  -- Thay đổi username ở đây
ORDER BY w.name;

-- Xem tất cả quyền của tất cả users
SELECT 
    u.username,
    r.name as role,
    w.name as warehouse,
    uwp.can_view,
    uwp.can_import,
    uwp.can_export,
    uwp.can_manage
FROM user_warehouse_permissions uwp
JOIN users u ON uwp.user_id = u.id
JOIN roles r ON u.role_id = r.id
JOIN warehouses w ON uwp.warehouse_id = w.id
ORDER BY u.username, w.name;

-- ============================================
-- 2. WAREHOUSE QUERIES
-- ============================================

-- Xem tất cả kho và quản lý
SELECT 
    w.code,
    w.name,
    w.address,
    u.full_name as manager,
    u.email as manager_email,
    CASE WHEN w.is_active THEN 'Active' ELSE 'Inactive' END as status
FROM warehouses w
LEFT JOIN users u ON w.manager_id = u.id
ORDER BY w.code;

-- Đếm số lượng sản phẩm trong mỗi kho
SELECT 
    w.name as warehouse,
    COUNT(DISTINCT i.product_id) as total_products,
    SUM(i.quantity) as total_items
FROM warehouses w
LEFT JOIN inventory i ON w.id = i.warehouse_id
GROUP BY w.id, w.name
ORDER BY w.name;

-- ============================================
-- 3. PRODUCT QUERIES
-- ============================================

-- Xem tất cả sản phẩm theo category
SELECT 
    category,
    COUNT(*) as product_count,
    STRING_AGG(name, ', ' ORDER BY name) as products
FROM products
WHERE is_active = true
GROUP BY category
ORDER BY category;

-- Xem sản phẩm với tổng tồn kho
SELECT 
    p.code,
    p.name,
    p.unit,
    p.category,
    COALESCE(SUM(i.quantity), 0) as total_stock,
    p.min_stock,
    CASE 
        WHEN COALESCE(SUM(i.quantity), 0) < p.min_stock THEN '⚠️ Low Stock'
        ELSE '✓ OK'
    END as stock_status
FROM products p
LEFT JOIN inventory i ON p.id = i.product_id
WHERE p.is_active = true
GROUP BY p.id, p.code, p.name, p.unit, p.category, p.min_stock
ORDER BY p.code;

-- Tìm sản phẩm sắp hết hàng
SELECT 
    p.code,
    p.name,
    COALESCE(SUM(i.quantity), 0) as current_stock,
    p.min_stock,
    (p.min_stock - COALESCE(SUM(i.quantity), 0)) as shortage
FROM products p
LEFT JOIN inventory i ON p.id = i.product_id
WHERE p.is_active = true
GROUP BY p.id, p.code, p.name, p.min_stock
HAVING COALESCE(SUM(i.quantity), 0) < p.min_stock
ORDER BY shortage DESC;

-- ============================================
-- 4. INVENTORY QUERIES
-- ============================================

-- Xem tồn kho chi tiết theo kho
SELECT 
    w.name as warehouse,
    p.code as product_code,
    p.name as product_name,
    i.quantity,
    p.unit,
    p.min_stock,
    CASE 
        WHEN i.quantity < p.min_stock THEN '⚠️ Low'
        WHEN i.quantity >= p.min_stock * 2 THEN '✓ High'
        ELSE '✓ OK'
    END as status,
    i.last_updated
FROM inventory i
JOIN warehouses w ON i.warehouse_id = w.id
JOIN products p ON i.product_id = p.id
ORDER BY w.name, p.code;

-- Tồn kho của một sản phẩm cụ thể trên tất cả kho
SELECT 
    w.code as warehouse_code,
    w.name as warehouse_name,
    COALESCE(i.quantity, 0) as quantity,
    p.unit
FROM warehouses w
CROSS JOIN products p
LEFT JOIN inventory i ON w.id = i.warehouse_id AND p.id = i.product_id
WHERE p.code = 'SP001'  -- Thay đổi mã sản phẩm ở đây
  AND w.is_active = true
ORDER BY w.name;

-- Top 10 sản phẩm có tồn kho nhiều nhất
SELECT 
    p.code,
    p.name,
    SUM(i.quantity) as total_quantity,
    p.unit,
    COUNT(DISTINCT i.warehouse_id) as warehouses_count
FROM products p
JOIN inventory i ON p.id = i.product_id
GROUP BY p.id, p.code, p.name, p.unit
ORDER BY total_quantity DESC
LIMIT 10;

-- ============================================
-- 5. TRANSACTION QUERIES
-- ============================================

-- Xem tất cả giao dịch
SELECT 
    t.transaction_code,
    t.type,
    w.name as warehouse,
    u.full_name as performed_by,
    t.status,
    t.transaction_date,
    t.note
FROM transactions t
JOIN warehouses w ON t.warehouse_id = w.id
JOIN users u ON t.user_id = u.id
ORDER BY t.transaction_date DESC;

-- Xem chi tiết một giao dịch
SELECT 
    t.transaction_code,
    t.type,
    w.name as warehouse,
    u.full_name as performed_by,
    t.transaction_date,
    p.code as product_code,
    p.name as product_name,
    td.quantity,
    p.unit,
    td.unit_price,
    (td.quantity * td.unit_price) as total_price,
    td.note
FROM transactions t
JOIN warehouses w ON t.warehouse_id = w.id
JOIN users u ON t.user_id = u.id
JOIN transaction_details td ON t.id = td.transaction_id
JOIN products p ON td.product_id = p.id
WHERE t.transaction_code = 'IMP-2024-001'  -- Thay đổi mã giao dịch ở đây
ORDER BY p.code;

-- Lịch sử nhập/xuất của một sản phẩm
SELECT 
    t.transaction_date,
    t.transaction_code,
    t.type,
    w.name as warehouse,
    td.quantity,
    p.unit,
    td.unit_price,
    u.full_name as performed_by,
    t.note
FROM transaction_details td
JOIN transactions t ON td.transaction_id = t.id
JOIN products p ON td.product_id = p.id
JOIN warehouses w ON t.warehouse_id = w.id
JOIN users u ON t.user_id = u.id
WHERE p.code = 'SP001'  -- Thay đổi mã sản phẩm ở đây
ORDER BY t.transaction_date DESC;

-- Thống kê giao dịch theo loại và trạng thái
SELECT 
    type,
    status,
    COUNT(*) as transaction_count,
    COUNT(DISTINCT warehouse_id) as warehouses_affected,
    COUNT(DISTINCT user_id) as users_involved
FROM transactions
GROUP BY type, status
ORDER BY type, status;

-- ============================================
-- 6. REPORT QUERIES
-- ============================================

-- Báo cáo nhập/xuất theo tháng
SELECT 
    TO_CHAR(t.transaction_date, 'YYYY-MM') as month,
    t.type,
    COUNT(*) as transaction_count,
    COUNT(DISTINCT td.product_id) as unique_products,
    SUM(td.quantity) as total_quantity
FROM transactions t
JOIN transaction_details td ON t.id = td.transaction_id
WHERE t.status = 'COMPLETED'
GROUP BY TO_CHAR(t.transaction_date, 'YYYY-MM'), t.type
ORDER BY month DESC, t.type;

-- Báo cáo giá trị nhập/xuất theo kho
SELECT 
    w.name as warehouse,
    t.type,
    COUNT(*) as transactions,
    SUM(td.quantity * td.unit_price) as total_value
FROM transactions t
JOIN warehouses w ON t.warehouse_id = w.id
JOIN transaction_details td ON t.id = td.transaction_id
WHERE t.status = 'COMPLETED'
GROUP BY w.id, w.name, t.type
ORDER BY w.name, t.type;

-- Top users hoạt động nhiều nhất
SELECT 
    u.username,
    u.full_name,
    r.name as role,
    COUNT(DISTINCT t.id) as total_transactions,
    SUM(CASE WHEN t.type = 'IMPORT' THEN 1 ELSE 0 END) as imports,
    SUM(CASE WHEN t.type = 'EXPORT' THEN 1 ELSE 0 END) as exports,
    MAX(t.transaction_date) as last_transaction
FROM users u
JOIN roles r ON u.role_id = r.id
LEFT JOIN transactions t ON u.id = t.user_id
GROUP BY u.id, u.username, u.full_name, r.name
HAVING COUNT(DISTINCT t.id) > 0
ORDER BY total_transactions DESC;

-- Báo cáo tồn kho chi tiết với giá trị
SELECT 
    w.name as warehouse,
    p.category,
    p.code,
    p.name as product,
    i.quantity,
    p.unit,
    COALESCE(
        (SELECT td.unit_price 
         FROM transaction_details td
         JOIN transactions t ON td.transaction_id = t.id
         WHERE td.product_id = p.id 
           AND t.warehouse_id = w.id
           AND t.type = 'IMPORT'
         ORDER BY t.transaction_date DESC
         LIMIT 1
        ), 0
    ) as last_import_price,
    i.quantity * COALESCE(
        (SELECT td.unit_price 
         FROM transaction_details td
         JOIN transactions t ON td.transaction_id = t.id
         WHERE td.product_id = p.id 
           AND t.warehouse_id = w.id
           AND t.type = 'IMPORT'
         ORDER BY t.transaction_date DESC
         LIMIT 1
        ), 0
    ) as estimated_value
FROM inventory i
JOIN warehouses w ON i.warehouse_id = w.id
JOIN products p ON i.product_id = p.id
WHERE i.quantity > 0
ORDER BY w.name, p.category, p.code;

-- ============================================
-- 7. MAINTENANCE QUERIES
-- ============================================

-- Kiểm tra database integrity
SELECT 
    'Total Users' as metric, 
    COUNT(*)::TEXT as value 
FROM users
UNION ALL
SELECT 
    'Active Users', 
    COUNT(*)::TEXT 
FROM users 
WHERE is_active = true
UNION ALL
SELECT 
    'Total Products', 
    COUNT(*)::TEXT 
FROM products
UNION ALL
SELECT 
    'Active Products', 
    COUNT(*)::TEXT 
FROM products 
WHERE is_active = true
UNION ALL
SELECT 
    'Total Warehouses', 
    COUNT(*)::TEXT 
FROM warehouses
UNION ALL
SELECT 
    'Total Transactions', 
    COUNT(*)::TEXT 
FROM transactions
UNION ALL
SELECT 
    'Completed Transactions', 
    COUNT(*)::TEXT 
FROM transactions 
WHERE status = 'COMPLETED';

-- Xem kích thước các bảng
SELECT 
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size
FROM pg_tables
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;

-- Xem các indexes
SELECT 
    tablename,
    indexname,
    indexdef
FROM pg_indexes
WHERE schemaname = 'public'
ORDER BY tablename, indexname;

-- ============================================
-- 8. PAGINATION QUERIES - Cho Frontend API
-- ============================================

-- ⚡ Pattern: LIMIT = số record mỗi trang, OFFSET = (page - 1) * limit
-- Ví dụ: Page 1 -> OFFSET 0, Page 2 -> OFFSET 20, Page 3 -> OFFSET 40...

-- 8.1. Lấy danh sách sản phẩm có phân trang
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

-- 8.2. Lấy danh sách tồn kho có phân trang
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

-- 8.3. Lấy danh sách giao dịch có phân trang
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

-- 8.4. Lấy chi tiết giao dịch theo transaction_id
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

-- 8.5. Lấy danh sách users có phân trang
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

-- 8.6. Tìm kiếm sản phẩm theo keyword (Full-text search)
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

-- 8.7. Báo cáo tồn kho theo kho với phân trang
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

-- 8.8. Lịch sử giao dịch của một sản phẩm (có phân trang)
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

-- 8.9. Cảnh báo sản phẩm sắp hết hàng (có phân trang)
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

-- 8.10. Pagination Metadata - Thông tin phân trang
-- Helper query để tính toán thông tin phân trang
-- Params: @page (current page), @page_size (default 20)
SELECT 
    COUNT(*) as total_records,
    20 as page_size,
    CEILING(COUNT(*)::DECIMAL / 20) as total_pages,
    1 as current_page,  -- Từ params @page
    CASE WHEN 1 > 1 THEN true ELSE false END as has_previous,
    CASE WHEN 1 < CEILING(COUNT(*)::DECIMAL / 20) THEN true ELSE false END as has_next,
    (1 - 1) * 20 as offset_value  -- Dùng trong query chính
FROM products
WHERE is_active = true;
