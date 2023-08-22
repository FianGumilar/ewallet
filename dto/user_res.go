package dto

type UserRes struct {
	UserId int64  `json:"-"`
	Token  string `json:"token"`
}
