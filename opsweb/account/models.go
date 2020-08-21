package account

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Name string `gorm:"type:varchar(50);not null;unique;uniqueIndex" validate:"required,min=3,max=50" json:"name"`
	CnName string `gorm:"type:varchar(100)" json:"cn_name"`
	Password string `gorm:"type:varchar(100);not null" validate:"required,min=8,max=50,pwdcomplex" json:"password"`
	Email string `validate:"omitempty,email" json:"email"`
	IsAdmin *bool `gorm:"not null;default:false" validate:"exists" json:"is_admin"`
	Role string `gorm:"type:varchar(20);not null" validate:"required,roleoptions" json:"role"`
	Phone string `gorm:"type:char(11);not null;index:idx_phone" validate:"required,len=11,phonecheck" json:"phone"`
	IsActive *bool `gorm:"not null;default:true" validate:"exists" json:"is_active"`
	LastLogin *time.Time `json:"last_login"`
	Group []Group `gorm:"many2many:user_groups;" json:"group"`
}

type UserPasswordUpdate struct {
	ID uint `json:"id" validate:"required"`
	Password string `validate:"required,min=8,max=50,pwdcomplex" json:"password"`
	ConfirmPassword string `validate:"required,min=8,max=50,pwdcomplex,eqfield=Password" json:"confirm_pwd"`
}

type UserStatusUpdate struct {
	ID uint `json:"id" validate:"required"`
	IsActive *bool `gorm:"not null;default:true" validate:"exists" json:"is_active"`
}

type UserAdminUpdate struct {
	ID uint `json:"id" validate:"required"`
	IsAdmin *bool `gorm:"not null;default:true" validate:"exists" json:"is_admin"`
}

type Group struct {
	gorm.Model
	Name string `gorm:"type:varchar(50);not null;unique;uniqueIndex" validate:"required,min=3,max=50" json:"name"`
	User []User `gorm:"many2many:user_groups;" json:"user"`
	Permissions []Permissions `gorm:"many2many:permission_groups" json:"permissions"`
}

type UserGroups struct {
	PType string `gorm:"type:varchar(100);not null;default:g"`
	UserID int `gorm:"primaryKey;index"`
	GroupID int `gorm:"primaryKey;index"`
}

type Permissions struct {
	ID int64 `gorm:"primaryKey" json:"id"`
	Path string `gorm:"varchar(100);not null;index" json:"path"`
	Method string  `gorm:"varchar(10);not null;index" json:"method"`
	Groups []Group `gorm:"many2many:permission_groups" json:"groups"`
	Describe string `gorm:"varchar(100);not null" json:"describe"`
}