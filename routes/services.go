package routes

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/mikelpsv/auth_service/app"
	"github.com/mikelpsv/auth_service/models"
	"log"
	"net/http"
	"strconv"
)

func RegisterServiceHandler(routeItems app.Routes) app.Routes {
	routeItems = append(routeItems, app.Route{
		Name:          "Token",
		Method:        "POST",
		Pattern:       "/token",
		SetHeaderJSON: true,
		ValidateToken: false,
		HandlerFunc:   requestToken,
	})
	return routeItems
}

func requestToken(w http.ResponseWriter, r *http.Request) {
	var clientSecret = ""

	grantType, isVarFound := app.GetSimpleValue(r, "grant_type")
	if !isVarFound {
		log.Println("Parameter grant_type is required")
		app.ResponseERROR(w, http.StatusBadRequest, errors.New("invalid_request"))
		return
	}

	clientId, isVarFound, err := app.GetSimpleValueAsInt(r, "client_id")
	if err != nil {
		log.Println("Invalid parameter client_id")
		app.ResponseERROR(w, http.StatusBadRequest, errors.New("invalid_request"))
		return
	}

	if !isVarFound {
		log.Println("Parameter client_id is required")
		app.ResponseERROR(w, http.StatusBadRequest, errors.New("invalid_request"))
		return
	}

	clientSecret, isVarFound = app.GetSimpleValue(r, "client_secret")

	if !isVarFound {
		log.Println("Parameter client_secret is required")
		app.ResponseERROR(w, http.StatusBadRequest, errors.New("invalid_request"))
		return
	}

	client := new(models.Client)
	client.FindById(clientId)
	res, err := client.ValidSecret(clientSecret)
	if !res || err != nil {
		log.Printf("Client valid secret %s, error %s", res, err)
		app.ResponseERROR(w, http.StatusBadRequest, errors.New("unauthorized_client"))
		return
	}

	if grantType == "password" {
		// Resource Owner Password Credentials
		// https://tools.ietf.org/html/draft-ietf-oauth-v2-13#section-4.3

		var username = ""
		var password = ""

		if username, isVarFound = app.GetSimpleValue(r, "username"); !isVarFound {
			app.ResponseERROR(w, http.StatusBadRequest, errors.New("invalid_request"))
			log.Println("Parameter username is required")
			return
		}

		if password, isVarFound = app.GetSimpleValue(r, "password"); !isVarFound {
			app.ResponseERROR(w, http.StatusBadRequest, errors.New("invalid_request"))
			log.Println("Parameter password is required")
			return
		}

		user := new(models.User)
		user.FindByUserName(username)

		if user.Id == 0 {
			app.ResponseERROR(w, http.StatusBadRequest, errors.New("invalid_grant"))
			return
		}
		if valid, _ := user.ValidPassword(password); !valid {
			app.ResponseERROR(w, http.StatusBadRequest, errors.New("invalid_grant"))
			return
		}

		tokenPair, _ := app.CreateTokenPair(user.Id, client.SecretKey, client.TokenExpires)
		app.ResponseJSON(w, http.StatusOK, tokenPair)

	} else if grantType == "refresh_token" {
		// Refreshing an Access Token
		// https://tools.ietf.org/html/draft-ietf-oauth-v2-13#section-6
		var refreshToken = ""

		if refreshToken, isVarFound = app.GetSimpleValue(r, "refresh_token"); !isVarFound {
			http.Error(w, "Parameter refresh_token is required", http.StatusBadRequest)
			return
		}

		token, _ := app.ReadToken(client.SecretKey, refreshToken)

		if token.Valid {

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				// return
			}

			userId, err := strconv.ParseInt(fmt.Sprintf("%.0f", claims["sub"]), 10, 64)
			if err != nil {
				return // ошибка чтения/конвертирования sub
			}

			tokenPair, _ := app.CreateTokenPair(userId, client.SecretKey, client.TokenExpires)
			app.ResponseJSON(w, http.StatusOK, tokenPair)
		}

	} else {
		app.ResponseERROR(w, http.StatusBadRequest, errors.New("unsupported_grant_type"))
		log.Println("unsupported_grant_type")
	}
}
