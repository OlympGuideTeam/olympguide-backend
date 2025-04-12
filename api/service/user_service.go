package service

import (
	"api/dto"
	"api/model"
	"api/repository"
	"api/utils/errs"
	"time"
)

type IUserService interface {
	GetUserData(userID uint) (*dto.UserDataResponse, bool, error)
	DeleteUser(userID uint) error
	UpdateUser(userID uint, request *dto.UpdateUserRequest) error
	UpdatePassword(userID uint, password string) error
}

type UserService struct {
	userRepo   repository.IUserRepo
	regionRepo repository.IRegionRepo
}

func NewUserService(userRepo repository.IUserRepo, regionRepo repository.IRegionRepo) *UserService {
	return &UserService{userRepo: userRepo, regionRepo: regionRepo}
}

func (u *UserService) GetUserData(userID uint) (*dto.UserDataResponse, bool, error) {
	user, err := u.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, false, err
	}
	if user.ProfileComplete {
		return newUserDataResponse(user), true, nil
	} else {
		return newPoorUserDataResponse(user), false, nil
	}

}

func (u *UserService) DeleteUser(userID uint) error {
	user, err := u.userRepo.GetUserByID(userID)
	if err != nil {
		return err
	}
	return u.userRepo.DeleteUser(user)
}

func (u *UserService) UpdateUser(userID uint, request *dto.UpdateUserRequest) error {
	user, err := u.userRepo.GetUserByID(userID)
	if err != nil {
		return err
	}

	parsedBirthday, err := time.Parse("02.01.2006", request.Birthdate)
	if err != nil {
		return errs.InvalidBirthday
	}

	if !u.regionRepo.RegionExists(request.RegionID) {
		return errs.RegionNotFound
	}

	user = updateUserModel(user, request, parsedBirthday)
	return u.userRepo.UpdateUser(user)
}

func (u *UserService) UpdatePassword(userID uint, password string) error {
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return err
	}
	user, err := u.userRepo.GetUserByID(userID)
	user.PasswordHash = hashedPassword

	err = u.userRepo.UpdateUser(user)
	if err != nil {
		return err
	}
	return nil
}

func newUserDataResponse(user *model.User) *dto.UserDataResponse {
	return &dto.UserDataResponse{
		PoorUserDataResponse: dto.PoorUserDataResponse{
			Email:      user.Email,
			SyncGoogle: user.GoogleID != nil,
			SyncApple:  user.AppleID != nil,
		},
		FirstName:  *user.FirstName,
		LastName:   *user.LastName,
		SecondName: *user.SecondName,
		Birthday:   user.Birthdate.Format("02.01.2006"),
		Region: dto.RegionResponse{
			RegionID: *user.RegionID,
			Name:     user.Region.Name,
		},
	}
}

func newPoorUserDataResponse(user *model.User) *dto.UserDataResponse {
	return &dto.UserDataResponse{
		PoorUserDataResponse: dto.PoorUserDataResponse{
			Email:      user.Email,
			SyncGoogle: user.GoogleID != nil,
			SyncApple:  user.AppleID != nil,
		},
	}
}

func updateUserModel(user *model.User, request *dto.UpdateUserRequest, birthdate time.Time) *model.User {
	user.FirstName = &request.FirstName
	user.LastName = &request.LastName
	user.SecondName = &request.SecondName
	user.RegionID = &request.RegionID
	user.Birthdate = &birthdate
	user.ProfileComplete = true
	return user
}
