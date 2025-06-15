package models

import (
	"context"
	"linn221/shop/utils"
	"log"
	"sync"

	"gorm.io/gorm"
)

type User struct {
	Id       int    `gorm:"primaryKey"`
	Username string `gorm:"unique;not null"`
	Email    string `gorm:"unique"`
	Password string `gorm:"index;not null"`
	HasIsActive
}

type UserService struct {
	db *gorm.DB
	mu sync.Mutex
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		db: db,
	}
}

func (s *UserService) Register(ctx context.Context, user *User) (*User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	err := Validate(s.db.WithContext(ctx),
		NewUniqueRule("users", "username", user.Username, 0, "duplicate username", nil),
		NewUniqueRule("users", "email", user.Email, 0, "duplicate email", nil),
	)
	if err != nil {
		return nil, err
	}

	h, err := utils.HashPassword(user.Password)
	if err != nil {
		log.Fatal(err.Error())
	}
	user.Password = string(h)

	if err := s.db.WithContext(ctx).Create(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) Login(ctx context.Context, username string, password string) (*User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var user User
	if err := s.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}

	err := utils.ComparePassword(user.Password, password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserService) ChangePassword(ctx context.Context, userId int, oldPassword string, newPassword string) (*User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var user User
	if err := s.db.WithContext(ctx).First(&user, userId).Error; err != nil {
		return nil, err
	}

	err := utils.ComparePassword(user.Password, oldPassword)
	if err != nil {
		return nil, err
	}
	h, err := utils.HashPassword(newPassword)
	if err != nil {
		return nil, err
	}
	if err := s.db.WithContext(ctx).Model(&user).UpdateColumn("password", string(h)).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
