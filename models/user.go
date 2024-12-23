package models

import (
	"time"
)

type User struct {
	ID                   uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name                 string    `gorm:"size:100;not null" json:"name" validate:"required"`
	Email                string    `gorm:"size:255;not null;unique" json:"email" validate:"required,email"`
	Pix                  string    `gorm:"size:100;unique" json:"pix" validate:"required"`
	CPF                  string    `gorm:"size:11;not null;unique" json:"cpf" validate:"required"`
	Password             string    `gorm:"size:255;not null" json:"-" validate:"required"`
	PasswordResetToken   string    `gorm:"size:255" json:"password_reset_token"`
	PasswordResetExpires time.Time `json:"password_reset_expires"`
	CreatedAt            time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
}
