package dto

type RegisterRequest struct {
	Username string `json:"username" binding:"required,alpha,min=3,max=20"`
	Password string `json:"password" binding:"required,ascii,min=6"`
}
