package service

import (
	"api/dto"
	"api/model"
	"api/repository"
)

type IUserService interface {
	GetUserData(userID uint) (*dto.UserDataResponse, error)
	DeleteUser(userID uint) error
}

type UserService struct {
	userRepo repository.IUserRepo
}

func NewUserService(userRepo repository.IUserRepo) *UserService {
	return &UserService{userRepo: userRepo}
}

func (u *UserService) GetUserData(userID uint) (*dto.UserDataResponse, error) {
	user, err := u.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	return newUserDataResponse(user), nil
}

func (u *UserService) DeleteUser(userID uint) error {
	user, err := u.userRepo.GetUserByID(userID)
	if err != nil {
		return err
	}
	return u.userRepo.DeleteUser(user)
}

func newUserDataResponse(user *model.User) *dto.UserDataResponse {
	return &dto.UserDataResponse{
		Email:      user.Email,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		SecondName: user.SecondName,
		Birthday:   user.Birthday.Format("02.01.2006"),
		Region: dto.RegionResponse{
			RegionID: user.RegionID,
			Name:     user.Region.Name,
		},
		SyncGoogle: user.GoogleID != "",
	}
}
