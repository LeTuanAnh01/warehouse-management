# 🎯 Server Modes - Test vs Production

Dự án có **2 servers riêng biệt** để phân tách môi trường test và production.

---

## 🧪 TEST Server (cmd/test/main.go)

### Mục đích:
- **Development & Testing ONLY**
- Không cần authentication
- Tất cả endpoints đều public
- Dùng để test API nhanh chóng

### Chạy Test Server:
```bash
cd "/Users/anhlt/Documents/LTA/Quản lý kho/backend-go"
go run cmd/test/main.go
```

### Đặc điểm:
- ⚠️ **NO AUTHENTICATION** - Tất cả endpoints đều mở
- 🔓 **Public Access** - Không cần token
- 🧪 **Test Data** - Dùng để test features
- ⚡ **Quick Testing** - Setup nhanh, test nhanh

### Endpoints (Test Mode):
```bash
# Health check
curl http://localhost:8080/health

# Get products (NO AUTH required)
curl http://localhost:8080/api/v1/products

# Create product (NO AUTH required)
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{"code":"SP999","name":"Test"}'
```

---

## 🚀 PRODUCTION Server (cmd/api/main.go)

### Mục đích:
- **Production Ready**
- Authentication bắt buộc
- Protected endpoints
- Dùng cho website thật

### Chạy Production Server:
```bash
cd "/Users/anhlt/Documents/LTA/Quản lý kho/backend-go"
go run cmd/api/main.go
```

### Đặc điểm:
- 🔒 **AUTHENTICATION REQUIRED** - Cần login
- 🛡️ **Protected Routes** - Phải có JWT token
- 👥 **User Management** - Register, Login, Logout
- 🔐 **Secure** - Password hashing, JWT tokens

### Authentication Flow:

#### 1. Register (Đăng ký user mới)
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "newuser",
    "email": "user@example.com",
    "password": "123456",
    "full_name": "New User",
    "phone": "0123456789"
  }'
```

**Response:**
```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "id": 7,
    "username": "newuser",
    "email": "user@example.com",
    "full_name": "New User",
    "role": "VIEWER"
  }
}
```

#### 2. Login (Đăng nhập)
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "newuser",
    "password": "123456"
  }'
```

**Response:**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "access_token": "eyJhbGci...",
    "refresh_token": "eyJhbGci...",
    "token_type": "Bearer",
    "expires_in": 86400,
    "user": {
      "id": 7,
      "username": "newuser",
      "email": "user@example.com",
      "full_name": "New User",
      "role": "VIEWER"
    }
  }
}
```

#### 3. Access Protected Endpoints (với token)
```bash
# Lưu token
TOKEN="eyJhbGci..."

# Get products (WITH AUTH)
curl http://localhost:8080/api/v1/products \
  -H "Authorization: Bearer $TOKEN"

# Create product (WITH AUTH)
curl -X POST http://localhost:8080/api/v1/products \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "code": "SP012",
    "name": "Secure Product",
    "unit": "Cái",
    "min_stock": 10
  }'
```

#### 4. Refresh Token (khi access token hết hạn)
```bash
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "eyJhbGci..."
  }'
```

#### 5. Change Password
```bash
curl -X POST http://localhost:8080/api/v1/auth/change-password \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "old_password": "123456",
    "new_password": "newpass123"
  }'
```

#### 6. Logout
```bash
curl -X POST http://localhost:8080/api/v1/auth/logout \
  -H "Authorization: Bearer $TOKEN"
```

---

## 📊 So sánh 2 Servers

| Feature | Test Server | Production Server |
|---------|-------------|-------------------|
| **Authentication** | ❌ Không có | ✅ Bắt buộc |
| **Public Endpoints** | ✅ Tất cả | ❌ Chỉ auth endpoints |
| **Protected Endpoints** | ❌ Không có | ✅ Products, User mgmt |
| **JWT Tokens** | ❌ Không dùng | ✅ Required |
| **User Management** | ❌ Không có | ✅ Register, Login |
| **Use Case** | Testing | Production Web |
| **Security** | ⚠️ Không secure | 🔒 Secure |

---

## 🎯 Khi nào dùng cái nào?

### Dùng TEST Server khi:
- ✅ Develop features mới
- ✅ Test API endpoints
- ✅ Debug nhanh
- ✅ Không cần authentication
- ✅ Thử nghiệm integration

### Dùng PRODUCTION Server khi:
- ✅ Connect với frontend (Vue.js)
- ✅ Deploy lên server thật
- ✅ Test authentication flow
- ✅ Test user permissions
- ✅ Production environment

---

## 🔧 Configuration

Cả 2 servers đều dùng cùng `.env` file:

```env
# Server
SERVER_PORT=8080
SERVER_HOST=localhost
ENV=development

# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=warehouse_db
DB_USER=anhlt
DB_PASSWORD=

# JWT (cho production server)
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
JWT_EXPIRATION=24h
```

---

## 🚦 Quick Start Commands

### Test Mode (No Auth):
```bash
# Terminal 1: Start test server
go run cmd/test/main.go

# Terminal 2: Test API
curl http://localhost:8080/api/v1/products
```

### Production Mode (With Auth):
```bash
# Terminal 1: Start production server
go run cmd/api/main.go

# Terminal 2: Register & Login
curl -X POST http://localhost:8080/api/v1/auth/register -H "Content-Type: application/json" -d '{"username":"test","email":"test@test.com","password":"123456","full_name":"Test"}'

curl -X POST http://localhost:8080/api/v1/auth/login -H "Content-Type: application/json" -d '{"username":"test","password":"123456"}'

# Use the token from login response
curl http://localhost:8080/api/v1/products -H "Authorization: Bearer YOUR_TOKEN"
```

---

## 📝 Notes

1. **Chỉ chạy 1 server tại 1 thời điểm** (cùng port 8080)
2. **Test server**: Dùng cho development, không deploy
3. **Production server**: Dùng cho web thật, có authentication
4. **Database**: Cả 2 đều dùng chung database `warehouse_db`
5. **Users trong DB**: Password đã hash bằng bcrypt, không thể login trực tiếp với seed data cũ

---

## 🎊 Summary

- **2 servers riêng biệt** = Quản lý dễ dàng
- **Test server** = Nhanh, đơn giản, không auth
- **Production server** = Secure, đầy đủ, có auth
- **Best practice** = Test trên test server, deploy production server

---

**Created**: March 9, 2026  
**Servers**: Test (cmd/test) | Production (cmd/api)
