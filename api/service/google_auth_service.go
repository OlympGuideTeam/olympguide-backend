package service

import (
	"api/dto"
	"api/model"
	"api/repository"
	"api/utils/errs"
	"context"
	"errors"
	"google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
	"time"
)

type IGoogleAuthService interface {
	GoogleAuth(token string) (*model.User, error)
	CompleteProfile(userID uint, req *dto.SignUpRequest) (*model.User, error)
}

type GoogleAuthService struct {
	userRepo   repository.IUserRepo
	regionRepo repository.IRegionRepo
}

func NewGoogleAuthService(userRepo repository.IUserRepo, regionRepo repository.IRegionRepo) *GoogleAuthService {
	return &GoogleAuthService{userRepo: userRepo, regionRepo: regionRepo}
}

func (s *GoogleAuthService) GoogleAuth(token string) (*model.User, error) {
	tokenInfo, err := validateGoogleToken(token)
	if err != nil {
		return nil, err
	}

	return s.findOrCreateGoogleUser(tokenInfo)
}

func (s *GoogleAuthService) CompleteProfile(userID uint, req *dto.SignUpRequest) (*model.User, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	if user.ProfileComplete {
		return nil, errs.RegistrationAlreadyCompleted
	}

	parsedBirthday, err := time.Parse("02.01.2006", req.Birthday)
	if err != nil {
		return nil, errs.InvalidBirthday
	}

	if !s.regionRepo.RegionExists(req.RegionID) {
		return nil, errs.RegionNotFound
	}

	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.SecondName = req.SecondName
	user.Birthday = parsedBirthday
	user.RegionID = req.RegionID
	user.PasswordHash = hashedPassword
	user.ProfileComplete = true

	return user, s.userRepo.UpdateUser(user)
}

func validateGoogleToken(token string) (*oauth2.Tokeninfo, error) {
	ctx := context.Background()
	oauth2Service, err := oauth2.NewService(ctx, option.WithoutAuthentication())
	if err != nil {
		return nil, err
	}

	tokenInfo, err := oauth2Service.Tokeninfo().IdToken(token).Do()
	if err != nil {
		return nil, err
	}

	if tokenInfo.Email == "" || !tokenInfo.VerifiedEmail {
		return nil, errors.New("invalid Google token")
	}
	return tokenInfo, nil
}

func (s *GoogleAuthService) findOrCreateGoogleUser(tokenInfo *oauth2.Tokeninfo) (*model.User, error) {
	if user, err := s.userRepo.GetUserByGoogleID(tokenInfo.UserId); err == nil {
		return user, nil
	}

	if user, err := s.userRepo.GetUserByEmail(tokenInfo.Email); err == nil {
		user.GoogleID = tokenInfo.UserId
		if err = s.userRepo.UpdateUser(user); err != nil {
			return nil, err
		}
		return user, nil
	}

	newUser := &model.User{
		Email:    tokenInfo.Email,
		GoogleID: tokenInfo.UserId,
	}

	if _, err := s.userRepo.CreateGoogleUser(newUser); err != nil {
		return nil, err
	}
	return newUser, nil
}
