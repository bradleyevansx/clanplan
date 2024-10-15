package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HandlerFunc func(r http.ResponseWriter, w *http.Request)

type App struct {
	Engine *gin.Engine
}

func NewApp() *App{
	e := gin.Default()
	return &App{
		Engine: e,
	}
}

func (a *App) Start(){
	a.Engine.Run()
}

func (a *App) HandlerFunc(method string, group string, path string, handlerFunc HandlerFunc){
	a.routeMethod(method, group, path, handlerFunc)
}

func (a *App) HandlerFuncNoMid(method string, group string, handlerFunc HandlerFunc){

}

func (a *App) routeMethod(method string, group string, path string, handlerFunc HandlerFunc){
	path = group + path
	switch method {
	case "GET":
		a.Engine.GET(path, func(c *gin.Context){
			handlerFunc(c.Writer, c.Request)
		})
	case "POST":
		a.Engine.POST(path, func(c *gin.Context){
			handlerFunc(c.Writer, c.Request)
		})
	}
}
