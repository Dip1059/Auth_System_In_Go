package Globals

import (
	Mod "Auth_System_In_Go/Models"
	"html/template"
)

type Message struct {
	Success template.HTML
	Fail template.HTML
}

type EmailGenerals struct {
	From, To, Subject, HtmlString string
}

type UserDataForEmail struct {
	EncEmail string
	User Mod.User
	PS Mod.PasswordReset
}

type DB_ENV struct {
	Host, Port, Dialect, Username, Password, DBname string
}

var (
	Msg Message
	User Mod.User
	Role Mod.Role
	PS Mod.PasswordReset
	DbEnv DB_ENV
)

const (
	DbName string = "go_crud"
)