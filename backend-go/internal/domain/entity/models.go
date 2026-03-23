package entity

import "time"

type Product struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Code        string    `gorm:"uniqueIndex;not null" json:"code"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `json:"description"`
	Unit        string    `json:"unit"`
	Category    string    `json:"category"`
	MinStock    int       `gorm:"default:0" json:"min_stock"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Product) TableName() string {
	return "products"
}

type Warehouse struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Code      string    `gorm:"uniqueIndex;not null" json:"code"`
	Name      string    `gorm:"not null" json:"name"`
	Address   string    `json:"address"`
	ManagerID *uint     `json:"manager_id"`
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Warehouse) TableName() string {
	return "warehouses"
}

type Inventory struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	WarehouseID uint      `gorm:"not null" json:"warehouse_id"`
	ProductID   uint      `gorm:"not null" json:"product_id"`
	Quantity    int       `gorm:"default:0" json:"quantity"`
	LastUpdated time.Time `gorm:"autoUpdateTime" json:"last_updated"`

	// Relationships
	Warehouse Warehouse `gorm:"foreignKey:WarehouseID" json:"warehouse,omitempty"`
	Product   Product   `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

func (Inventory) TableName() string {
	return "inventory"
}

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Username     string    `gorm:"uniqueIndex;not null" json:"username"`
	Email        string    `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string    `gorm:"not null" json:"-"`
	FullName     string    `json:"full_name"`
	Phone        string    `json:"phone"`
	RoleID       uint      `json:"role_id"`
	IsActive     bool      `gorm:"default:true" json:"is_active"`
	LastLogin    *time.Time `json:"last_login"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Relationships
	Role Role `gorm:"foreignKey:RoleID" json:"role,omitempty"`
}

func (User) TableName() string {
	return "users"
}

type Role struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"uniqueIndex;not null" json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Role) TableName() string {
	return "roles"
}
