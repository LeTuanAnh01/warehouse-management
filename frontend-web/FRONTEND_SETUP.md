# 🎨 Warehouse Management Frontend

Frontend Vue.js 3 application for Warehouse Management System.

## ✅ Đã Setup

### 1. Dependencies
- ✅ Vue 3 với Composition API
- ✅ Vue Router 4
- ✅ Pinia (State Management)
- ✅ Axios (API calls)
- ✅ TailwindCSS (Styling)
- ✅ Heroicons (Icons)

### 2. Project Structure
```
src/
├── services/         # API service layer
│   ├── api.js       # Axios instance với interceptors
│   └── index.js     # Auth & Product services
├── stores/          # Pinia stores
│   ├── auth.js      # Authentication state
│   └── product.js   # Product state
├── router/          # Vue Router
│   └── index.js     # Routes & navigation guards
├── layouts/         # Layout components (cần tạo)
├── views/           # Page components (cần tạo)
└── components/      # Reusable components (cần tạo)
```

### 3. Features Implemented
- ✅ API client với authentication interceptor
- ✅ Auth service (login, register, logout, refresh token)
- ✅ Product service (CRUD operations)
- ✅ Auth store với localStorage persistence
- ✅ Product store với loading states
- ✅ Router với authentication guards
- ✅ TailwindCSS configured

---

## 🚀 Next Steps - Tạo UI Components

### Bước 1: Tạo App.vue
```bash
# File: src/App.vue
```

### Bước 2: Tạo Layouts
- `src/layouts/DashboardLayout.vue` - Main dashboard layout với sidebar

### Bước 3: Tạo Views (Pages)
- `src/views/LoginView.vue` - Login page
- `src/views/RegisterView.vue` - Register page  
- `src/views/DashboardView.vue` - Dashboard home
- `src/views/ProductsView.vue` - Products list
- `src/views/ProductFormView.vue` - Create/Edit product
- `src/views/ProfileView.vue` - User profile

### Bước 4: Tạo Components
- `src/components/Navbar.vue` - Top navigation bar
- `src/components/Sidebar.vue` - Sidebar menu
- `src/components/ProductCard.vue` - Product card component
- `src/components/Pagination.vue` - Pagination component

---

## 🏃 Run Development Server

```bash
cd "/Users/anhlt/Documents/LTA/Quản lý kho/frontend-web"
npm run dev
```

Server sẽ chạy tại: **http://localhost:5173**

---

## 🔗 Connect với Backend

Backend API đang chạy tại: **http://localhost:8080**

Axios đã được config để tự động:
- Thêm `Authorization: Bearer <token>` vào mọi request
- Redirect về `/login` khi token expires (401)
- Handle refresh token

---

## 📝 API Endpoints Available

### Authentication
- `POST /auth/login` - Login
- `POST /auth/register` - Register
- `POST /auth/refresh` - Refresh token
- `POST /auth/logout` - Logout

### Products (Protected)
- `GET /products` - List products (pagination)
- `GET /products/:id` - Get product
- `POST /products` - Create product
- `PUT /products/:id` - Update product
- `DELETE /products/:id` - Delete product
- `GET /products/:id/inventory` - Get inventory

---

## 🎨 TailwindCSS Classes

Custom classes available:
- `.btn` - Base button
- `.btn-primary` - Primary button
- `.btn-secondary` - Secondary button
- `.btn-danger` - Danger button
- `.input` - Input field
- `.card` - Card container
- `.label` - Form label

---

## 💡 Tôi sẽ tạo tiếp UI components?

Bạn muốn tôi:
1. **Tạo tất cả UI components** - Login, Dashboard, Products, etc.
2. **Tạo từng phần** - Bắt đầu với Login page trước
3. **Hướng dẫn bạn tự tạo** - Tôi chỉ cung cấp template

Chọn option nào? 🤔
