package repository

import (
	"api/model"
	"gorm.io/gorm"
)

type IUserRepo interface {
	CreateUser(user *model.User) (*model.User, error)
	GetUserByEmail(email string) (*model.User, error)
	GetUserByID(userID uint) (*model.User, error)
	GetUserByGoogleID(googleID string) (*model.User, error)
	GetUserByAppleID(appleID string) (*model.User, error)
	DeleteUser(user *model.User) error
	UpdateUser(user *model.User) error
	Exists(userID uint) bool
}

type PgUserRepo struct {
	db *gorm.DB
}

func NewPgUserRepo(db *gorm.DB) *PgUserRepo {
	return &PgUserRepo{db: db}
}

func (u *PgUserRepo) CreateUser(user *model.User) (*model.User, error) {
	if err := u.db.Create(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (u *PgUserRepo) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	if err := u.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *PgUserRepo) GetUserByID(userID uint) (*model.User, error) {
	var user model.User
	if err := u.db.Preload("Region").Where("user_id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *PgUserRepo) GetUserByGoogleID(googleID string) (*model.User, error) {
	var user model.User
	if err := u.db.Where("google_id = ?", googleID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *PgUserRepo) GetUserByAppleID(appleID string) (*model.User, error) {
	var user model.User
	if err := u.db.Where("apple_id = ?", appleID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *PgUserRepo) UpdateUser(user *model.User) error {
	return u.db.Save(user).Error
}

func (u *PgUserRepo) Exists(userID uint) bool {
	var count int64
	u.db.Model(&model.User{}).Where("user_id = ?", userID).Count(&count)
	return count > 0
}

func (u *PgUserRepo) DeleteUser(user *model.User) error {
	return u.db.Delete(user).Error
}
