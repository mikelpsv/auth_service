package routes

import (
	"fmt"
	"github.com/mikelpsv/auth_service/app"
	"github.com/mikelpsv/auth_service/models"
	"net/http"
)



func RegisterServiceHandler(routeItems app.Routes) app.Routes {
	routeItems = append(routeItems, app.Route{
		Name:          "Authorize",
		Method:        "POST",
		Pattern:       "/authorize",
		SetHeaderJSON: true,
		ValidateToken: false,
		HandlerFunc:   nil,
	})
	routeItems = append(routeItems, app.Route{
		Name:          "Token",
		Method:        "GET",
		Pattern:       "/token",
		SetHeaderJSON: true,
		ValidateToken: false,
		HandlerFunc:   requestToken,
	})
	return routeItems
}

func requestToken(w http.ResponseWriter, r *http.Request) {

	grantType, isVarFound := app.GetSimpleValue(r, "grant_type")
	if !isVarFound {
		http.Error(w, "Parameter grant_type is required", http.StatusBadRequest)
		return
	}



	if grantType == "password"{
		// Resource Owner Password Credentials
		// https://tools.ietf.org/html/draft-ietf-oauth-v2-13#section-4.3


//		var clientId int64
		var clientSecret = ""
		var username  = ""
		var password = ""

		clientId, isVarFound, err := app.GetSimpleValueAsInt(r, "client_id")
		if err!=nil{
			http.Error(w, "Invalid parameter client_id", http.StatusBadRequest)
			return
		}

		if !isVarFound{
			http.Error(w, "Parameter client_id is required", http.StatusBadRequest)
			return
		}

		if clientSecret, isVarFound = app.GetSimpleValue(r, "client_secret"); !isVarFound{
			http.Error(w, "Parameter client_secret is required", http.StatusBadRequest)
			return
		}

		if username, isVarFound = app.GetSimpleValue(r, "username"); !isVarFound{
			http.Error(w, "Parameter username is required", http.StatusBadRequest)
			return
		}

		if password, isVarFound = app.GetSimpleValue(r, "password"); !isVarFound{
			http.Error(w, "Parameter password is required", http.StatusBadRequest)
			return
		}

		getToken(grantType, clientId, clientSecret, username, password)
	}


}

func getToken(grantType string, clientId int64, clientSecret string, username string, password string){
	client := models.Client{}
	client.FindById(clientId)
	tokenPair, _ := app.CreateTokenPair(10, client.Secret, client.TokenExpires)
	fmt.Println(tokenPair.AccessToken)
}