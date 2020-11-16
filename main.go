package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/mikelpsv/auth_service/app"
	"log"
	"net/http"
	"os"
)

type AppCfg struct {
	APP_ADDR string
	APP_PORT string
	DB_HOST  string
	DB_PORT  string
	DB_NAME  string
	DB_USER  string
	DB_PASS  string
}

var Config AppCfg

func main() {
	ReadEnv()
	routeItems := app.Routes{}
	routeItems = RegisterHandlers(routeItems)
	router := NewRouter(routeItems)
	app.InitDb(Config.DB_HOST, Config.DB_NAME, Config.DB_USER, Config.DB_PASS)
	app.Install()
	log.Fatal(http.ListenAndServe(Config.APP_ADDR+":"+Config.APP_PORT, router))
}

func ReadEnv() {
	e := godotenv.Load()
	if e != nil {
		fmt.Print(e)
	}

	Config.APP_ADDR = os.Getenv("APP_ADDR")
	Config.APP_PORT = os.Getenv("APP_PORT")
	Config.DB_HOST = os.Getenv("DB_HOST")
	Config.DB_PORT = os.Getenv("DB_PORT")
	Config.DB_NAME = os.Getenv("DB_NAME")
	Config.DB_USER = os.Getenv("DB_USER")
	Config.DB_PASS = os.Getenv("DB_PASS")
}

func NewRouter(routeItems app.Routes) *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routeItems {
		handlerFunc := route.HandlerFunc
		if route.ValidateToken {
			handlerFunc = app.SetMiddlewareAuth(handlerFunc)
		}

		if route.SetHeaderJSON {
			handlerFunc = app.SetMiddlewareJSON(handlerFunc)
		}

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			HandlerFunc(handlerFunc)
	}

	return router
}
