package dto

type UserData struct {
	ID       int64  `json:"id"`
	FullName string `json:"fullname"`
	Phone    string `json:"phone"`
	Username string `json:"username"`
}
