# Database - Warehouse Management System

## 📊 Database Structure

### Tables

1. **roles** - Vai trò trong hệ thống
2. **users** - Người dùng
3. **warehouses** - Kho hàng
4. **user_warehouse_permissions** - Phân quyền user theo kho
5. **products** - Sản phẩm
6. **inventory** - Tồn kho
7. **transactions** - Giao dịch nhập/xuất
8. **transaction_details** - Chi tiết giao dịch

## 🚀 Setup Instructions

### 1. Cài đặt PostgreSQL

```bash
# macOS
brew install postgresql@16
brew services start postgresql@16

# hoặc sử dụng Docker
docker run --name warehouse-db \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=warehouse_db \
  -p 5432:5432 \
  -d postgres:16
```

### 2. Tạo Database

```bash
# Kết nối PostgreSQL
psql -U postgres

# Tạo database
CREATE DATABASE warehouse_db;

# Thoát
\q
```

### 3. Import Schema

```bash
# Import schema
psql -U postgres -d warehouse_db -f schema.sql

# Import seed data
psql -U postgres -d warehouse_db -f seed.sql
```

## 👥 Default Accounts (Password: `123456`)

| Username | Email | Role | Description |
|----------|-------|------|-------------|
| admin | admin@warehouse.com | ADMIN | Quản trị viên |
| manager1 | manager1@warehouse.com | WAREHOUSE_MANAGER | Quản lý kho HN, ĐN |
| manager2 | manager2@warehouse.com | WAREHOUSE_MANAGER | Quản lý kho HCM |
| staff1 | staff1@warehouse.com | WAREHOUSE_STAFF | Nhân viên kho HN |
| staff2 | staff2@warehouse.com | WAREHOUSE_STAFF | Nhân viên kho HCM |
| viewer1 | viewer1@warehouse.com | VIEWER | Người xem (kho HN) |

## 🔑 Permission System

### Warehouse Permissions

| Permission | Description |
|------------|-------------|
| can_view | Xem thông tin kho |
| can_import | Nhập hàng vào kho |
| can_export | Xuất hàng ra khỏi kho |
| can_manage | Quản lý kho (full quyền) |

## 📝 Database Features

### Auto Update Inventory
- Khi tạo transaction detail (nhập/xuất), inventory sẽ tự động cập nhật
- IMPORT: Tăng số lượng tồn kho
- EXPORT: Giảm số lượng tồn kho

### Auto Update Timestamps
- Các bảng có trường `updated_at` sẽ tự động cập nhật khi record được update

### Indexes
- Đã tạo indexes cho các trường thường xuyên query để tối ưu performance

## 🔍 Useful Queries

### Xem tồn kho theo kho

```sql
SELECT 
    w.name as warehouse_name,
    p.code as product_code,
    p.name as product_name,
    i.quantity,
    p.unit
FROM inventory i
JOIN warehouses w ON i.warehouse_id = w.id
JOIN products p ON i.product_id = p.id
WHERE w.id = 1
ORDER BY p.name;
```

### Xem lịch sử nhập/xuất

```sql
SELECT 
    t.transaction_code,
    t.type,
    w.name as warehouse_name,
    u.full_name as user_name,
    t.transaction_date,
    t.note
FROM transactions t
JOIN warehouses w ON t.warehouse_id = w.id
JOIN users u ON t.user_id = u.id
ORDER BY t.transaction_date DESC;
```

### Kiểm tra quyền của user

```sql
SELECT 
    u.username,
    w.name as warehouse_name,
    uwp.can_view,
    uwp.can_import,
    uwp.can_export,
    uwp.can_manage
FROM user_warehouse_permissions uwp
JOIN users u ON uwp.user_id = u.id
JOIN warehouses w ON uwp.warehouse_id = w.id
WHERE u.username = 'staff1';
```

## 🛠️ Maintenance

### Backup Database

```bash
pg_dump -U postgres warehouse_db > backup.sql
```

### Restore Database

```bash
psql -U postgres warehouse_db < backup.sql
```

### Reset Database

```bash
psql -U postgres -d warehouse_db -f schema.sql
psql -U postgres -d warehouse_db -f seed.sql
```
