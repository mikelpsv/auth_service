package models

import (
	"fmt"
	"github.com/mikelpsv/auth_service/app"
)

type Client struct {
	Id           int64  `json:"client_id"`
	Name         string `json:"name"`
	TokenExpires int64  `json:"token_expires"`
	Secret       string `json:"client_secret"`
	SecretKey    string `json:"secret_key"`
}

func (c *Client) FindById(clientId int64) *Client {
	rows, err := app.Db.Query("SELECT _id, name, secret, key, expires FROM clients WHERE _id = $1", clientId)
	if err != nil {
		return &Client{}
	}
	defer rows.Close()

	if !rows.Next() {
		return &Client{}
	}
	err = rows.Scan(&c.Id, &c.Name, &c.Secret, &c.SecretKey, &c.TokenExpires)
	if err != nil {
		return &Client{}
	}
	return c
}

func (c *Client) ValidSecret(secretString string) (bool, error) {
	return app.ValidPassword(c.Secret, secretString)
}

func (c *Client) UpdateSecret(secretString string) {
	_, err := app.Db.Exec("UPDATE clients SET secret = $1 WHERE _id=$2", secretString, c.Id)
	if err != nil {
		fmt.Println(err)
	}
}
