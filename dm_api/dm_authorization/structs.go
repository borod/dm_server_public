package dm_authorization

import (
	_ "dm_server/dm_db"
)

type ResetPasswordBody struct {
	Email    string `json:"Email"`
	Token    string `json:"Token"`
	Password string `json:"Password"`
}

type UserNamePasswordBody struct {
	Email    string `json:"Email"`
	Password string `json:"Password"`
	Name     string `json:"Name"`
}

type LoginResponse struct {
	UserID             int64  `json:"UserID"`
	AuthorizationToken string `json:"Token"`
	RefreshToken       string `json:"RefreshToken"`
	RedmineUserID      string `json:"RedmineUserID"`
	RedmineToken       string `json:"RedmineToken"`
	Name               string `json:"Name"`
	Nickname           string `json:"Nickname"`
	Email              string `json:"Email"`
	Phone              string `json:"Phone"`
}
