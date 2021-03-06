package main

import (
	"flag"
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
}

var Config AppCfg

func main() {
	var pFlagInstall, pFlagTestData bool

	flag.BoolVar(&pFlagInstall, "install", false, "Install database. Data will be lost!")
	flag.BoolVar(&pFlagTestData, "testdata", false, "Load test data. Use with a key -install")
	flag.Parse()

	ReadEnv()
	routeItems := app.Routes{}
	routeItems = RegisterHandlers(routeItems)
	router := NewRouter(routeItems)

	if pFlagInstall {
		_, err := os.Stat(".data/authdata.db")
		if err == nil {
			log.Println("Data file is exist!")
			os.Remove(".data/authdata.db")
		}
	}

	app.InitDb("", "", "", "")

	if pFlagInstall {
		app.Install(pFlagTestData)
	}

	log.Fatal(http.ListenAndServe(Config.APP_ADDR+":"+Config.APP_PORT, router))
}

func ReadEnv() {
	e := godotenv.Load()
	if e != nil {
		fmt.Print(e)
	}

	Config.APP_ADDR = os.Getenv("APP_ADDR")
	Config.APP_PORT = os.Getenv("APP_PORT")
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
