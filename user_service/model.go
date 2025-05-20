package userservice

import "time"

type UserModel struct {
	ID           string    `gorm:"primarykey" json:"id"`
	Name         string    `json:"name"`
	Password     string    `json:"_"`
	Email        string    `gorm:"uniqueIndex" json:"email"`
	Phone        string    `json:"phone"` // must for admin
	UserType     string    `json:"user_type"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UserModelList struct {
	TotalPage int         `json:"total_page"`
	TotalUser int         `json:"total_user"`
	UserModel []UserModel `json:"user_list"`
}
