package models

import (
	"time"
)

type User struct {
	ID                   uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name                 string    `gorm:"size:100;not null" json:"name"`
	Email                string    `gorm:"size:255;not null;unique" json:"email"`
	Pix                  string    `gorm:"size:100;unique" json:"pix"`
	CPF                  string    `gorm:"size:11;not null;unique" json:"cpf"`
	Password             string    `gorm:"size:255;not null" json:"-"`
	PasswordResetToken   string    `gorm:"size:255;unique" json:"password_reset_token"`
	PasswordResetExpires time.Time `json:"password_reset_expires"`
	CreatedAt            time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
}
