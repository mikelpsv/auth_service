package models

import (
	"fmt"
	"github.com/mikelpsv/auth_service/app"
)

type User struct {
	Id       int64  `json:"user_id"`
	Name     string `json:"name"`
	Username int64  `json:"username"`
	Password string `json:"password"`
	ClientId int64  `json:"client_id"`
}

func (u *User) FindById(userId int64) *User {
	rows, err := app.Db.Query("SELECT _id, username, password, client_id FROM users WHERE _id = $1", userId)
	if err != nil {
		return &User{}
	}
	defer rows.Close()

	if !rows.Next() {
		return &User{}
	}
	err = rows.Scan(&u.Id, &u.Username, &u.Password, &u.ClientId)
	if err != nil {
		return &User{}
	}
	return u
}

func (u *User) ValidPassword(passwordString string) (bool, error) {
	return app.ValidPassword(u.Password, passwordString)
}

func (u *User) UpdatePassword(secretString string) {
	_, err := app.Db.Exec("UPDATE users SET password = $1 WHERE _id=$2", secretString, u.Id)
	if err != nil {
		fmt.Println(err)
	}
}

func (u *User) FindByUserName(username string) *User {
	rows, err := app.Db.Query("SELECT _id, username, password, client_id FROM users WHERE username = $1", username)
	if err != nil {
		return &User{}
	}
	defer rows.Close()

	if !rows.Next() {
		return &User{}
	}
	err = rows.Scan(&u.Id, &u.Username, &u.Password, &u.ClientId)
	if err != nil {
		return &User{}
	}
	return u
}
