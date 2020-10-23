package models

import (
	"github.com/mikelpsv/auth_service/app"
)

type Client struct {
	Id int64 `json:"client_id"`
	Name string `json:"name"`
	TokenExpires int64 `json:"token_expires"`
	Secret string `json:"client_secret"`
}


func (c *Client) FindById(clientId int64) *Client {
	rows, err := app.Db.Query("SELECT _id, name, secret, expires FROM clients WHERE _id = $1", clientId)
	if err != nil {
		return &Client{}
	}
	defer rows.Close()

	if !rows.Next(){
		return &Client{}
	}
	err = rows.Scan(&c.Id, &c.Name, &c.Secret, &c.TokenExpires)
	if err!= nil{
		return &Client{}
	}
	return c
}

