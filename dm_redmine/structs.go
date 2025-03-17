package dm_redmine

import (
	_ "dm_server/dm_db"
	"time"
)

type RedmineUserResponse struct {
	User RedmineUser `json:"user"`
}

type RedmineUser struct {
	ID              int       `json:"id"`
	Login           string    `json:"login"`
	Admin           bool      `json:"admin"`
	Firstname       string    `json:"firstname"`
	Lastname        string    `json:"lastname"`
	Mail            string    `json:"mail"`
	CreatedOn       time.Time `json:"created_on"`
	UpdatedOn       time.Time `json:"updated_on"`
	LastLoginOn     time.Time `json:"last_login_on,omitempty"`
	PasswdChangedOn time.Time `json:"passwd_changed_on"`
	TwofaScheme     string    `json:"twofa_scheme,omitempty"`
	ApiKey          string    `json:"api_key"`
	Status          int       `json:"status"`
}

type UserBody struct {
	Login     string `json:"login"`
	Password  string `json:"password"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Mail      string `json:"mail"`
}

type RedmineCreateUserBody struct {
	User UserBody `json:"user"`
}
