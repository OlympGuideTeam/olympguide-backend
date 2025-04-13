package service

import (
	"context"
	"fmt"
	"storage_service/internal/repository"
)

type ILogoService interface {
	UploadLogo(ctx context.Context, universityID string, data []byte, ext string) (string, string, error)
	DeleteLogo(ctx context.Context, universityID string) error
}

type LogoService struct {
	logoRepo repository.ILogoRepository
}

func NewLogoService(repo repository.ILogoRepository) *LogoService {
	return &LogoService{logoRepo: repo}
}

func (s *LogoService) UploadLogo(ctx context.Context, universityID string, data []byte, ext string) (string, string, error) {
	objectName := fmt.Sprintf("logos/%s.%s", universityID, ext)
	contentType := "image/" + ext

	url, err := s.logoRepo.Upload(ctx, objectName, data, contentType, ext)
	if err != nil {
		return "", "", err
	}
	return objectName, url, nil
}

func (s *LogoService) DeleteLogo(ctx context.Context, universityID string) error {
	objectName := fmt.Sprintf("logos/%s", universityID)
	ext, err := s.logoRepo.GetFileExtension(ctx, objectName)
	if err != nil {
		return err
	}
	objectNameWithExt := fmt.Sprintf("%s.%s", objectName, ext)
	return s.logoRepo.Delete(ctx, objectNameWithExt)
}
