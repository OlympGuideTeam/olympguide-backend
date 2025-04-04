package middleware

import "api/service"

type Mw struct {
	adminService service.IAdminService
	tokenService service.ITokenService
}

func NewMw(adminService service.IAdminService, tokenService service.ITokenService) *Mw {
	return &Mw{adminService: adminService, tokenService: tokenService}
}
