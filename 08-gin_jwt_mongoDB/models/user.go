package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive" // mongo驱动程序(官方)
	"time"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id"`
	UserId       string             `json:"user_id"`
	Token        *string            `json:"token"`
	RefreshToken *string            `json:"refresh_token"`
	FirstName    *string            `json:"first_name" validate:"required,min=2,max=10"`
	LastName     *string            `json:"last_name" validate:"required"`
	Password     *string            `json:"password" validate:"required,min=6"`
	Email        *string            `json:"email" validate:"email,required"`
	Phone        *string            `json:"phone" validate:"required"`
	UserType     *string            `json:"user_type" validate:"required,eq=ADMIN|eq=USER"` // eq=ADMIN|eq=USER 相当于枚举出他的类型
	CreatedAt    time.Time          `json:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at"`
	DeletedAt    time.Time          `json:"deleted_at"`
}
