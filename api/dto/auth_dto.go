package dto

type UpdateUserRequest struct {
	FirstName  string `json:"first_name" binding:"required"`
	LastName   string `json:"last_name" binding:"required"`
	SecondName string `json:"second_name" binding:"omitempty,min=1"`
	Birthdate  string `json:"birthdate" binding:"required"`
	RegionID   uint   `json:"region_id" binding:"required"`
}

type EmailSignUpRequest struct {
	Email    string
	Password string `json:"password" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type SendRequest struct {
	Email string `json:"email" binding:"required"`
}

type VerifyRequest struct {
	Email string `json:"email" binding:"required"`
	Code  string `json:"code" binding:"required"`
}

type ExternalAuthRequest struct {
	Token string `json:"token" binding:"required"`
}

type UpdatePasswordRequest struct {
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Message   string `json:"message" example:"Logged in"`
	FirstName string `json:"first_name" example:"Арсений"`
	LastName  string `json:"last_name" example:"Титаренко"`
}

type RegistrationIncompleteResponse struct {
	Message string `json:"message" example:"Uncompleted registration"`
	Token   string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

type VerifyCodeResponse struct {
	Message string `json:"message" example:"Email confirmed"`
	Token   string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

type MessageResponse struct {
	Message string `json:"message" example:"Message sent"`
}
