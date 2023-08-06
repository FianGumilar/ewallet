package dto

type UserRegisterReq struct {
	FullName string `json:"fullname"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}
