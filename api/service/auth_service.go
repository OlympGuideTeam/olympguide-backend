package service

import (
	"api/dto"
	"api/model"
	"api/repository"
	"api/utils/constants"
	"api/utils/errs"
	"context"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
)

type IAuthService interface {
	SendCode(email string) error
	VerifyCode(email, code string) error
	SignUp(request *dto.EmailSignUpRequest) (*model.User, error)
	Login(email, password string) (*model.User, error)
}

type AuthService struct {
	codeRepo   repository.ICodeRepo
	userRepo   repository.IUserRepo
	regionRepo repository.IRegionRepo
}

func NewAuthService(codeRepo repository.ICodeRepo, userRepo repository.IUserRepo, regionRepo repository.IRegionRepo) *AuthService {
	return &AuthService{codeRepo: codeRepo, userRepo: userRepo, regionRepo: regionRepo}
}

func (s *AuthService) SendCode(email string) error {
	exists, err := s.codeRepo.CodeExists(context.Background(), email)
	if err != nil {
		return err
	}

	if exists {
		ttl, err := s.codeRepo.GetCodeTTL(context.Background(), email)
		if err != nil {
			return err
		}
		details := map[string]interface{}{"ttl": ttl.Seconds()}
		return errs.PreviousCodeNotExpired.WithAdditional(details)
	}

	code := generateCode()

	if err := s.codeRepo.SetCode(context.Background(), email, code, constants.MaxVerifyCodeAttempts, constants.EmailCodeTtl); err != nil {
		return err
	}
	if err := s.codeRepo.PublishEmailCode(context.Background(), email, code); err != nil {
		return err
	}
	return nil
}

func (s *AuthService) VerifyCode(email, requestCode string) error {
	storedCode, attempts, err := s.codeRepo.GetCodeInfo(context.Background(), email)
	if err != nil {
		return err
	} else if attempts == -1 {
		return errs.CodeNotFoundOrExpired
	}

	if attempts == 0 {
		return errs.TooManyAttempts
	}

	if storedCode != requestCode {
		if err = s.codeRepo.DecreaseCodeAttempt(context.Background(), email); err != nil {
			return err
		}
		return errs.InvalidCode
	}

	if err = s.codeRepo.DeleteCode(context.Background(), email); err != nil {
		return err
	}
	return nil
}

func (s *AuthService) SignUp(request *dto.EmailSignUpRequest) (*model.User, error) {
	hashedPassword, err := HashPassword(request.Password)
	if err != nil {
		return nil, err
	}

	user := newPoorUserModel(request, hashedPassword)
	user, err = s.userRepo.CreateUser(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *AuthService) Login(email, password string) (*model.User, error) {
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return nil, errs.UserNotFound
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errs.InvalidPassword
	}

	return user, nil
}

func generateCode() string {
	code := fmt.Sprintf("%04d", rand.Intn(10000))
	return code
}

func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func newPoorUserModel(request *dto.EmailSignUpRequest, pwdHash string) *model.User {
	return &model.User{
		Email:        request.Email,
		PasswordHash: pwdHash,
	}
}
