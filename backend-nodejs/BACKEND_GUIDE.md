# 📚 Tài liệu Backend – Quản lý Kho

> Tài liệu này giải thích toàn bộ kiến thức backend trong project theo ngôn ngữ đơn giản,
> dựa trực tiếp trên code đang chạy. Không cần kinh nghiệm trước.

---

## Mục lục

1. [Backend là gì?](#1-backend-là-gì)
2. [Kiến trúc tổng quan](#2-kiến-trúc-tổng-quan)
3. [HTTP & API là gì?](#3-http--api-là-gì)
4. [NestJS – Framework đang dùng](#4-nestjs--framework-đang-dùng)
5. [Module – Cách tổ chức code](#5-module--cách-tổ-chức-code)
6. [Entity – Bảng trong database](#6-entity--bảng-trong-database)
7. [Repository & TypeORM – Truy vấn database](#7-repository--typeorm--truy-vấn-database)
8. [Service – Logic nghiệp vụ](#8-service--logic-nghiệp-vụ)
9. [Controller – Nhận request từ client](#9-controller--nhận-request-từ-client)
10. [DTO – Dữ liệu đầu vào](#10-dto--dữ-liệu-đầu-vào)
11. [Authentication – Xác thực người dùng](#11-authentication--xác-thực-người-dùng)
12. [JWT – Token bảo mật](#12-jwt--token-bảo-mật)
13. [Guard – Bảo vệ API](#13-guard--bảo-vệ-api)
14. [Biến môi trường (.env)](#14-biến-môi-trường-env)
15. [Luồng chạy đầy đủ – Ví dụ login](#15-luồng-chạy-đầy-đủ--ví-dụ-login)
16. [Danh sách API hiện có](#16-danh-sách-api-hiện-có)

---

## 1. Backend là gì?

**Frontend** (Vue.js) là phần người dùng nhìn thấy: giao diện, nút bấm, form nhập liệu.

**Backend** (NestJS) là phần chạy sau, người dùng không nhìn thấy, chịu trách nhiệm:
- Nhận yêu cầu từ frontend (ví dụ: "lấy danh sách sản phẩm")
- Xử lý logic (kiểm tra quyền, tính toán, validate dữ liệu)
- Đọc/ghi dữ liệu vào database (PostgreSQL)
- Trả kết quả về cho frontend (dạng JSON)

```
Người dùng
    │
    ▼
Frontend (Vue.js :5173)
    │  HTTP Request (GET /api/products)
    ▼
Backend (NestJS :3000)  ◄──► Database (PostgreSQL :5432)
    │  JSON Response
    ▼
Frontend hiển thị dữ liệu
```

---

## 2. Kiến trúc tổng quan

Project này theo mô hình **3 tầng**:

```
src/
├── auth/           ← Đăng nhập, đăng ký
├── users/          ← Quản lý người dùng
├── products/       ← Quản lý sản phẩm
├── warehouses/     ← Quản lý kho
├── inventory/      ← Tồn kho
└── common/
    └── guards/     ← Bảo vệ API (yêu cầu đăng nhập)
```

Mỗi thư mục là một **module**, gồm 4 file chính:

| File | Vai trò |
|------|---------|
| `*.entity.ts` | Định nghĩa bảng database |
| `*.service.ts` | Xử lý logic |
| `*.controller.ts` | Nhận/trả HTTP request |
| `*.module.ts` | Kết nối các phần lại với nhau |

---

## 3. HTTP & API là gì?

Khi frontend cần dữ liệu, nó gửi **HTTP Request** đến backend. Request có:

- **Method**: loại hành động
- **URL**: địa chỉ tài nguyên
- **Body**: dữ liệu gửi kèm (với POST/PATCH)
- **Header**: thông tin phụ (token xác thực, loại content)

| Method | Ý nghĩa | Ví dụ |
|--------|---------|-------|
| `GET` | Lấy dữ liệu | Lấy danh sách sản phẩm |
| `POST` | Tạo mới | Thêm sản phẩm mới |
| `PATCH` | Cập nhật | Sửa tên sản phẩm |
| `DELETE` | Xóa | Xóa sản phẩm |

Backend trả về **HTTP Response** gồm:
- **Status code**: kết quả (200 OK, 201 Created, 400 lỗi, 401 chưa đăng nhập, 404 không tìm thấy)
- **Body**: dữ liệu JSON

---

## 4. NestJS – Framework đang dùng

**NestJS** là framework viết backend bằng TypeScript (ngôn ngữ giống JS nhưng có kiểu dữ liệu).

NestJS dùng khái niệm **Decorator** – những annotation bắt đầu bằng `@` để mô tả chức năng:

```typescript
@Controller('api/products')  // ← URL prefix: /api/products
export class ProductsController {

  @Get()          // ← GET /api/products
  findAll() { ... }

  @Get(':id')     // ← GET /api/products/1
  findOne(@Param('id') id: string) { ... }

  @Post()         // ← POST /api/products
  create(@Body() dto: any) { ... }
}
```

---

## 5. Module – Cách tổ chức code

Mỗi tính năng được đóng gói thành 1 **Module**. Ví dụ `products.module.ts`:

```typescript
@Module({
  imports: [TypeOrmModule.forFeature([Product])],  // ← Đăng ký entity Product
  controllers: [ProductsController],               // ← Controller xử lý request
  providers: [ProductsService],                    // ← Service chứa logic
})
export class ProductsModule {}
```

`app.module.ts` là **module gốc** – tập hợp tất cả module lại:

```typescript
@Module({
  imports: [
    ConfigModule.forRoot({ isGlobal: true }),    // ← Đọc file .env
    TypeOrmModule.forRootAsync({ ... }),         // ← Kết nối database
    AuthModule,
    UsersModule,
    ProductsModule,
    WarehousesModule,
    InventoryModule,
  ],
})
export class AppModule {}
```

---

## 6. Entity – Bảng trong database

**Entity** là class TypeScript ánh xạ tới một bảng trong database.
Mỗi property = một cột trong bảng.

Ví dụ `product.entity.ts`:

```typescript
@Entity('products')          // ← Ánh xạ tới bảng "products" trong PostgreSQL
export class Product {

  @PrimaryGeneratedColumn()  // ← Cột id, tự tăng (1, 2, 3...)
  id: number;

  @Column({ unique: true })  // ← Cột code, không được trùng
  code: string;

  @Column()
  name: string;

  @Column({ default: 0 })    // ← Giá trị mặc định là 0
  minStock: number;

  @Column({ name: 'is_active', default: true })  // ← Tên cột trong DB là is_active
  isActive: boolean;

  @CreateDateColumn({ name: 'created_at' })  // ← Tự điền thời gian tạo
  createdAt: Date;
}
```

**Quan hệ giữa bảng** – Ví dụ `inventory.entity.ts`:

```typescript
@Entity('inventory')
export class Inventory {
  @Column({ name: 'warehouse_id' })
  warehouseId: number;

  // ManyToOne = nhiều inventory thuộc 1 warehouse
  @ManyToOne(() => Warehouse, { eager: true })  // eager: tự join khi query
  @JoinColumn({ name: 'warehouse_id' })
  warehouse: Warehouse;
}
```

---

## 7. Repository & TypeORM – Truy vấn database

**TypeORM** là thư viện giúp thao tác database mà không cần viết SQL thuần.

**Repository** là object cung cấp các hàm CRUD cho 1 entity:

```typescript
// Inject repository vào service
constructor(
  @InjectRepository(Product)
  private productsRepository: Repository<Product>,
) {}

// Lấy tất cả sản phẩm → SELECT * FROM products
this.productsRepository.find()

// Lấy theo điều kiện → SELECT * FROM products WHERE id = 1
this.productsRepository.findOne({ where: { id: 1 } })

// Tạo mới → INSERT INTO products (...)
const product = this.productsRepository.create({ name: 'Gạo', code: 'SP001' })
this.productsRepository.save(product)

// Cập nhật → UPDATE products SET name='...' WHERE id=1
this.productsRepository.update(1, { name: 'Gạo thơm' })

// Xóa → DELETE FROM products WHERE id=1
this.productsRepository.remove(product)
```

---

## 8. Service – Logic nghiệp vụ

**Service** chứa toàn bộ logic xử lý. Controller gọi Service, Service gọi Repository.

Ví dụ `products.service.ts`:

```typescript
@Injectable()  // ← Cho phép NestJS inject (truyền) service này vào nơi khác
export class ProductsService {

  async findOne(id: number) {
    const product = await this.productsRepository.findOne({ where: { id } });

    if (!product) {
      throw new NotFoundException('Product not found');  // ← Tự động trả 404
    }

    return product;
  }

  async remove(id: number) {
    const product = await this.findOne(id);  // ← Gọi hàm nội bộ
    await this.productsRepository.remove(product);
    return { message: 'Product deleted' };
  }
}
```

**Dependency Injection (DI)** – Khái niệm quan trọng:

Thay vì tự tạo object (`new ProductsService()`), NestJS **tự động tạo và truyền vào**:

```typescript
// Controller không cần tự tạo service, NestJS làm điều đó
constructor(private productsService: ProductsService) {}
//          ^--- NestJS tự inject ProductsService vào đây
```

---

## 9. Controller – Nhận request từ client

**Controller** là cổng vào – nhận HTTP request và trả về response.
Controller chỉ nhận/trả dữ liệu, **không chứa logic**.

```typescript
@Controller('api/products')  // ← Base URL: /api/products
@UseGuards(JwtAuthGuard)     // ← Yêu cầu đăng nhập mới vào được
export class ProductsController {

  constructor(private productsService: ProductsService) {}

  // GET /api/products
  @Get()
  findAll() {
    return this.productsService.findAll();
  }

  // GET /api/products/5
  @Get(':id')
  findOne(@Param('id') id: string) {   // ← :id từ URL
    return this.productsService.findOne(+id);  // +id chuyển string → number
  }

  // POST /api/products
  @Post()
  create(@Body() dto: any) {   // ← Dữ liệu từ body request
    return this.productsService.create(dto);
  }

  // PATCH /api/products/5
  @Patch(':id')
  update(@Param('id') id: string, @Body() dto: any) {
    return this.productsService.update(+id, dto);
  }

  // DELETE /api/products/5
  @Delete(':id')
  remove(@Param('id') id: string) {
    return this.productsService.remove(+id);
  }
}
```

---

## 10. DTO – Dữ liệu đầu vào

**DTO (Data Transfer Object)** định nghĩa cấu trúc dữ liệu được phép gửi lên.

Ví dụ `auth.dto.ts`:

```typescript
export class RegisterDto {
  @IsString()
  @IsNotEmpty()
  username: string;   // ← Bắt buộc phải có, phải là string

  @IsEmail()
  email: string;      // ← Bắt buộc phải đúng định dạng email

  @IsString()
  @MinLength(6)
  password: string;   // ← Tối thiểu 6 ký tự

  fullName: string;
  phone?: string;     // ← ? nghĩa là không bắt buộc
}
```

Khi client gửi sai (thiếu email, password quá ngắn...), NestJS tự động trả lỗi `400 Bad Request` nhờ `ValidationPipe` trong `main.ts`:

```typescript
app.useGlobalPipes(
  new ValidationPipe({ whitelist: true, transform: true }),
  // whitelist: bỏ qua các field không được khai báo trong DTO
  // transform: tự chuyển kiểu dữ liệu (string → number...)
);
```

---

## 11. Authentication – Xác thực người dùng

Authentication = **"Bạn là ai?"** – Xác nhận danh tính người dùng.

Luồng trong project:

**Đăng ký (`POST /api/auth/register`)**:
```
Client gửi { username, email, password }
    │
    ▼
Kiểm tra username/email đã tồn tại chưa
    │
    ▼
Hash password bằng bcrypt (không lưu mật khẩu thô)
    │
    ▼
Lưu user vào database
    │
    ▼
Trả về thông tin user (không có password)
```

**Đăng nhập (`POST /api/auth/login`)**:
```
Client gửi { username, password }
    │
    ▼
Tìm user theo username
    │
    ▼
So sánh password với hash trong DB (bcrypt.compare)
    │
    ▼
Tạo JWT token chứa { userId, username, roleId }
    │
    ▼
Trả về { access_token: "eyJ...", user: {...} }
```

**bcrypt** – Tại sao không lưu mật khẩu thô?

```
Password: "123456"
    │ bcrypt.hash("123456", 10)
    ▼
Hash:  "$2b$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy"

// Hash này KHÔNG thể giải mã ngược → database bị lộ cũng không lấy được mật khẩu
// Để kiểm tra: bcrypt.compare("123456", hash) → true/false
```

---

## 12. JWT – Token bảo mật

**JWT (JSON Web Token)** là chuỗi mã hóa chứa thông tin người dùng.

Cấu trúc JWT gồm 3 phần ngăn cách bởi dấu `.`:

```
eyJhbGciOiJIUzI1NiJ9  .  eyJzdWIiOjEsInVzZXJuYW1lIjoiYWRtaW4ifQ  .  SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c
      Header                            Payload                                    Signature
   (thuật toán)              (dữ liệu: userId, username)              (chữ ký xác thực)
```

- **Header**: thuật toán mã hóa (HS256)
- **Payload**: dữ liệu lưu trong token – có thể decode ra đọc được nhưng **không thể giả mạo**
- **Signature**: ký bằng `JWT_SECRET` trong `.env` – nếu ai sửa payload thì signature sai

```typescript
// Tạo token khi login (auth.service.ts)
const payload = { sub: user.id, username: user.username, roleId: user.roleId };
const token = this.jwtService.sign(payload);
// Token có hiệu lực 7 ngày (JWT_EXPIRES_IN=7d trong .env)
```

**Cách dùng token**: Client lưu token, mỗi request tiếp theo gửi kèm trong Header:

```
Authorization: Bearer eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOjF9.xxx
```

---

## 13. Guard – Bảo vệ API

**Guard** quyết định request có được phép vào controller không.

`JwtAuthGuard` kiểm tra token hợp lệ hay không:

```typescript
// jwt.strategy.ts – Logic kiểm tra token
async validate(payload: { sub: number; username: string; roleId: number }) {
  const user = await this.authService.validateUser(payload.sub);
  if (!user) throw new UnauthorizedException();
  return user;  // ← user này sẽ được gán vào req.user
}
```

Áp dụng Guard cho toàn bộ controller:

```typescript
@Controller('api/products')
@UseGuards(JwtAuthGuard)   // ← Tất cả route trong này đều cần đăng nhập
export class ProductsController { ... }
```

Luồng khi có Guard:

```
Request đến GET /api/products
    │
    ▼
JwtAuthGuard: Có token không?
    ├── Không có → 401 Unauthorized
    ├── Token hết hạn → 401 Unauthorized
    └── Token hợp lệ → gán user vào req.user → vào Controller
```

`@Request() req` trong controller cho phép lấy thông tin user đang đăng nhập:

```typescript
@Get('me')
getProfile(@Request() req) {
  return this.usersService.findOne(req.user.id);  // ← lấy từ token
}
```

---

## 14. Biến môi trường (.env)

File `.env` chứa **cấu hình bí mật** không được commit lên git:

```env
PORT=3000                   # ← Backend chạy port nào
NODE_ENV=development        # ← Môi trường (development/production)

DB_HOST=localhost           # ← Địa chỉ database
DB_PORT=5432
DB_USERNAME=admin
DB_PASSWORD=123456
DB_NAME=warehouse

JWT_SECRET=your-secret-key  # ← Khóa ký JWT, PHẢI đổi khi lên production
JWT_EXPIRES_IN=7d
```

Đọc trong code qua `ConfigService`:

```typescript
// app.module.ts
host: config.get('DB_HOST'),      // ← Đọc DB_HOST từ .env
port: config.get<number>('DB_PORT'),
```

**Tại sao không hardcode trực tiếp?**
- Môi trường dev và production có DB khác nhau
- Không lộ thông tin nhạy cảm khi push code lên GitHub

---

## 15. Luồng chạy đầy đủ – Ví dụ login

Dưới đây là toàn bộ hành trình khi frontend gọi `POST /api/auth/login`:

```
1. Frontend gửi request:
   POST http://localhost:3000/api/auth/login
   Body: { "username": "admin", "password": "123456" }

2. main.ts: ValidationPipe kiểm tra body có đúng LoginDto không

3. auth.controller.ts:
   @Post('login')
   login(@Body() dto: LoginDto) {
     return this.authService.login(dto);  ← gọi service
   }

4. auth.service.ts:
   - Tìm user có username="admin" trong database
   - bcrypt.compare("123456", user.passwordHash) → true
   - Tạo JWT: { sub: 1, username: "admin", roleId: 1 }
   - jwtService.sign(payload) → "eyJhbGci..."

5. Trả về response:
   {
     "access_token": "eyJhbGci...",
     "user": { "id": 1, "username": "admin", "email": "..." }
   }

6. Frontend lưu access_token vào localStorage
   → Mỗi request tiếp theo thêm: Authorization: Bearer eyJhbGci...
```

---

## 16. Danh sách API hiện có

### Auth (không cần token)
| Method | URL | Mô tả |
|--------|-----|-------|
| POST | `/api/auth/register` | Đăng ký tài khoản mới |
| POST | `/api/auth/login` | Đăng nhập, nhận token |

### Users (cần token)
| Method | URL | Mô tả |
|--------|-----|-------|
| GET | `/api/users/me` | Lấy thông tin bản thân |
| GET | `/api/users` | Danh sách tất cả user |
| GET | `/api/users/:id` | Chi tiết 1 user |
| PATCH | `/api/users/:id` | Cập nhật user |
| DELETE | `/api/users/:id` | Xóa user |

### Products (cần token)
| Method | URL | Mô tả |
|--------|-----|-------|
| GET | `/api/products` | Danh sách sản phẩm |
| GET | `/api/products/:id` | Chi tiết sản phẩm |
| POST | `/api/products` | Thêm sản phẩm |
| PATCH | `/api/products/:id` | Sửa sản phẩm |
| DELETE | `/api/products/:id` | Xóa sản phẩm |

### Warehouses (cần token)
| Method | URL | Mô tả |
|--------|-----|-------|
| GET | `/api/warehouses` | Danh sách kho |
| GET | `/api/warehouses/:id` | Chi tiết kho |
| POST | `/api/warehouses` | Thêm kho |
| PATCH | `/api/warehouses/:id` | Sửa kho |
| DELETE | `/api/warehouses/:id` | Xóa kho |

### Inventory (cần token)
| Method | URL | Mô tả |
|--------|-----|-------|
| GET | `/api/inventory` | Toàn bộ tồn kho |
| GET | `/api/inventory/warehouse/:warehouseId` | Tồn kho theo kho |
| GET | `/api/inventory/:id` | Chi tiết 1 bản ghi |
| POST | `/api/inventory` | Thêm bản ghi tồn kho |
| PATCH | `/api/inventory/:id` | Cập nhật số lượng |
| DELETE | `/api/inventory/:id` | Xóa bản ghi |

---

## Tóm tắt nhanh

```
HTTP Request
    │
    ▼
main.ts          → Cấu hình chung (CORS, ValidationPipe)
    │
    ▼
app.module.ts    → Kết nối DB, đăng ký tất cả module
    │
    ▼
Controller       → Nhận request, gọi Service
    │
    ▼
Service          → Xử lý logic, gọi Repository
    │
    ▼
Repository       → Đọc/ghi PostgreSQL
    │
    ▼
Entity           → Ánh xạ bảng database
```

**Bộ từ khóa cần nhớ:**

| Khái niệm | Nghĩa ngắn gọn |
|-----------|----------------|
| Module | Nhóm tính năng độc lập |
| Entity | Class = Bảng DB |
| Repository | Công cụ query DB |
| Service | Nơi chứa logic |
| Controller | Cổng nhận HTTP |
| DTO | Khuôn dữ liệu đầu vào |
| Guard | Bảo vệ route |
| JWT | Token xác thực |
| bcrypt | Mã hóa mật khẩu |
| .env | Cấu hình bí mật |
