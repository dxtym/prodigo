package dto

type LoginRequest struct {
	Username string `json:"username" binding:"required,alpha,min=3,max=20"`
	Password string `json:"password" binding:"required,ascii,min=8"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
