# 📚 Database Queries Guide - Hướng dẫn sử dụng

## 📂 File Structure

```
database/
├── schema.sql              # Cấu trúc database (tables, triggers, indexes)
├── seed.sql                # Dữ liệu mẫu
├── queries.sql             # Tất cả queries (bao gồm pagination)
├── .env.example            # Config kết nối
└── README.md               # Tài liệu
```

## 🎯 Cách sử dụng queries.sql

File `queries.sql` được chia thành **8 sections chính**:

### 1️⃣ USER & PERMISSION QUERIES
- Xem danh sách users và roles
- Kiểm tra quyền của user theo kho
- Quản lý phân quyền

### 2️⃣ WAREHOUSE QUERIES
- Danh sách kho và quản lý
- Thống kê sản phẩm theo kho

### 3️⃣ PRODUCT QUERIES
- Danh sách sản phẩm theo category
- Tìm sản phẩm sắp hết hàng
- Xem tổng tồn kho

### 4️⃣ INVENTORY QUERIES
- Tồn kho chi tiết theo kho
- Tồn kho của sản phẩm trên tất cả kho
- Top sản phẩm tồn kho nhiều nhất

### 5️⃣ TRANSACTION QUERIES
- Lịch sử nhập/xuất
- Chi tiết giao dịch
- Lịch sử theo sản phẩm
- Thống kê giao dịch

### 6️⃣ REPORT QUERIES
- Báo cáo nhập/xuất theo tháng
- Báo cáo giá trị theo kho
- Top users hoạt động nhiều nhất
- Báo cáo tồn kho với giá trị

### 7️⃣ MAINTENANCE QUERIES
- Kiểm tra database integrity
- Xem kích thước bảng
- Xem indexes

### 8️⃣ PAGINATION QUERIES (⭐ CHO FRONTEND API)
- **8.1** Products với pagination
- **8.2** Inventory với pagination
- **8.3** Transactions với pagination
- **8.4** Transaction details với pagination
- **8.5** Users với pagination
- **8.6** Search products
- **8.7** Warehouse inventory report
- **8.8** Product history
- **8.9** Low stock alert
- **8.10** Pagination metadata

## 🚀 Cách chạy queries

### Option 1: Từ command line
```bash
# Chạy một query cụ thể
psql -d warehouse_db -c "SELECT * FROM users LIMIT 5;"

# Chạy từ file
psql -d warehouse_db -f queries.sql
```

### Option 2: Từ psql interactive
```bash
# Vào psql
psql -d warehouse_db

# Copy query từ queries.sql và paste vào
# Hoặc sử dụng \i command
\i /path/to/queries.sql
```

## 📖 Pagination Pattern cho Frontend

### Concept
```
LIMIT = số records mỗi trang (default: 20)
OFFSET = (page - 1) * limit

Page 1: OFFSET 0   (records 1-20)
Page 2: OFFSET 20  (records 21-40)
Page 3: OFFSET 40  (records 41-60)
```

### Example: Lấy sản phẩm trang 2

```sql
-- Page 2, 20 records per page
-- OFFSET = (2 - 1) * 20 = 20

WITH product_list AS (
    SELECT 
        p.id,
        p.code,
        p.name,
        COALESCE(SUM(i.quantity), 0) as total_stock
    FROM products p
    LEFT JOIN inventory i ON p.id = i.product_id
    WHERE p.is_active = true
    GROUP BY p.id
)
SELECT 
    *,
    (SELECT COUNT(*) FROM product_list) as total_count
FROM product_list
ORDER BY code
LIMIT 20 OFFSET 20;  -- Page 2
```

### Response Format cho API

```json
{
  "data": [...],           // Dữ liệu trang hiện tại
  "pagination": {
    "total_records": 150,
    "page_size": 20,
    "total_pages": 8,
    "current_page": 2,
    "has_previous": true,
    "has_next": true
  }
}
```

## 🎨 Backend API Endpoints (Suggested)

```
GET /api/products?page=1&limit=20
GET /api/inventory?warehouse_id=1&page=1&limit=20
GET /api/transactions?type=IMPORT&page=1&limit=20
GET /api/users?role_id=2&page=1&limit=20
GET /api/search?keyword=laptop&page=1&limit=20
GET /api/low-stock?warehouse_id=1&page=1&limit=20
```

## 🔧 Parameters Template

Khi implement trong Backend, thay thế:

```sql
-- Thay vì hard-code:
WHERE p.code = 'SP001'

-- Sử dụng parameter (Golang example):
WHERE p.code = $1

-- Vue.js sẽ gọi:
GET /api/products/SP001
```

### Common Parameters

| Parameter | Type | Description | Example |
|-----------|------|-------------|---------|
| `page` | int | Trang hiện tại | 1, 2, 3... |
| `limit` | int | Records mỗi trang | 10, 20, 50 |
| `warehouse_id` | int | Filter theo kho | 1, 2, 3 |
| `product_id` | int | Filter theo sản phẩm | 1, 2, 3 |
| `type` | string | IMPORT/EXPORT | "IMPORT" |
| `category` | string | Category sản phẩm | "Electronics" |
| `keyword` | string | Search keyword | "laptop" |
| `start_date` | date | Từ ngày | "2024-01-01" |
| `end_date` | date | Đến ngày | "2024-12-31" |

## 💡 Tips & Best Practices

### 1. Luôn trả về total_count
```sql
SELECT 
    *,
    (SELECT COUNT(*) FROM product_list) as total_count
FROM product_list
```
Frontend cần `total_count` để tính số trang.

### 2. Validate page parameters
```go
// Backend validation
if page < 1 {
    page = 1
}
if limit < 1 || limit > 100 {
    limit = 20 // default
}
offset := (page - 1) * limit
```

### 3. Index quan trọng
Đã có indexes cho:
- `products.code`
- `inventory.warehouse_id`, `inventory.product_id`
- `transactions.warehouse_id`, `transactions.type`, `transactions.transaction_date`

### 4. Cache khi có thể
Với dữ liệu ít thay đổi (roles, warehouses), nên cache ở Backend.

### 5. Sorting options
Cho phép Frontend chọn sort:
```sql
ORDER BY 
    CASE WHEN $sort_by = 'code' THEN p.code END,
    CASE WHEN $sort_by = 'name' THEN p.name END
```

## 🧪 Testing Queries

```bash
# Test tồn kho trang 1
psql -d warehouse_db << 'EOF'
WITH inventory_list AS (
    SELECT w.name as warehouse, p.code, i.quantity
    FROM inventory i
    JOIN warehouses w ON i.warehouse_id = w.id
    JOIN products p ON i.product_id = p.id
)
SELECT *, (SELECT COUNT(*) FROM inventory_list) as total
FROM inventory_list
ORDER BY warehouse, code
LIMIT 10 OFFSET 0;
EOF
```

## 📊 Performance Tips

1. **Use EXPLAIN ANALYZE** để check performance:
```sql
EXPLAIN ANALYZE
SELECT * FROM inventory WHERE warehouse_id = 1;
```

2. **Limit kết quả** khi test:
```sql
LIMIT 10  -- Test với ít data trước
```

3. **Monitor slow queries**:
```sql
-- Enable slow query log
ALTER DATABASE warehouse_db SET log_min_duration_statement = 1000;
```

## 🎓 Next Steps

1. ✅ Database setup xong
2. ⏭️ Implement Backend API (Golang)
3. ⏭️ Implement Frontend (Vue.js)
4. ⏭️ Connect FE ↔️ BE ↔️ DB

---

**Ready để build Backend API! 🚀**
