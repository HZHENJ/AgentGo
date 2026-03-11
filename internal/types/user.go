package types

type UserRegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Captcha string `json:"captcha" binding:"required"`
	Password string `json:"password"`
}

type UserInfo struct {
	Id 	 uint   `json:"id"`
	Email string `json:"email"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
}

type UserRegisterResponse struct {
	Token string `json:"token"`
	User UserInfo `json:"user"`
}

type UserLoginRequest struct  {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserLoginResponse struct {
	Token string `json:"token"`
	User UserInfo `json:"user"`
}

type SendCaptchaRequest struct {
	Email string `json:"email" binding:"required,email"`
}