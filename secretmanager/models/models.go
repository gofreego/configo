package models

type Secret struct {
	Token    *string `json:"token"`
	Username *string `json:"username"`
	Password *string `json:"password"`
}
