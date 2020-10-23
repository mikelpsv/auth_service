package main

/**
  Данный файл содержит изменяемую часть сервера Api
  Список методов и функции-обработчики
*/

import (
	"github.com/mikelpsv/auth_service/app"
	"github.com/mikelpsv/auth_service/routes"
	//"github.com/mikelpsv/auth_service/routes"
)

func RegisterHandlers(routeItems app.Routes) app.Routes {
	routeItems = routes.RegisterServiceHandler(routeItems)
	//routeItems = routes.RegisterCustomerHandler(routeItems)

	return routeItems
}
