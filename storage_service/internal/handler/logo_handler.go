package handler

import (
	"context"
	"storage_service/internal/service"
	pb "storage_service/proto/gen"
)

type LogoHandler struct {
	pb.UnimplementedStorageServiceServer
	logoService service.ILogoService
}

func NewLogoHandler(logoService service.ILogoService) *LogoHandler {
	return &LogoHandler{logoService: logoService}
}

func (h *LogoHandler) UploadLogo(ctx context.Context, req *pb.UploadLogoRequest) (*pb.UploadLogoResponse, error) {
	objectName, url, err := h.logoService.UploadLogo(ctx, req.UniversityId, req.FileData, req.FileExtension)
	if err != nil {
		return nil, err
	}
	return &pb.UploadLogoResponse{ObjectName: objectName, Url: url}, nil
}

func (h *LogoHandler) DeleteLogo(ctx context.Context, req *pb.DeleteLogoRequest) (*pb.DeleteLogoResponse, error) {
	err := h.logoService.DeleteLogo(ctx, req.UniversityId, req.FileExtension)
	return &pb.DeleteLogoResponse{Success: err == nil}, err
}
